package user

import (
	"encoding/json"
	"fmt"
	"go-brunel/internal/pkg/server/endpoint/api"
	"go-brunel/internal/pkg/server/security"
	"go-brunel/internal/pkg/server/store"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/pkg/errors"
)

type authHandler struct {
	defaultAdminUser string
	serializer       security.TokenSerializer
	userStore        store.UserStore
}

func (handler *authHandler) sendErrorMessage(w http.ResponseWriter, err string) {
	if _, err := fmt.Fprintf(
		w,
		`<html><head></head><body><script>window.opener.postMessage({error: '%s'}, '*');</script></body></html>`,
		err,
	); err != nil {
		log.Error("error occurred attempting to write error to client ", err)
	}
}

func (handler *authHandler) oAuthComplete(w http.ResponseWriter, user goth.User) {
	email := strings.TrimSpace(user.Email)
	if len(email) == 0 {
		handler.sendErrorMessage(w, "You must have at least one verified primary email address to login using oAuth.")
		return
	}

	roles, err := handler.userStore.AddOrUpdate(store.User{
		Username:  email,
		Email:     email,
		Name:      strings.TrimSpace(user.Name),
		AvatarURL: strings.TrimSpace(user.AvatarURL),
	})
	if err != nil {
		log.Error("error occurred updating user ", err)
		handler.sendErrorMessage(w, "Error attempting to log in.")
		return
	}

	if roles.Username == handler.defaultAdminUser {
		roles.Role = security.UserRoleAdmin
	}

	token, err := handler.serializer.Encode(security.Identity{Username: roles.Username, Role: roles.Role})
	if err != nil {
		log.Error("error occurred creating jwt token ", err)
		handler.sendErrorMessage(w, "Error attempting to log in.")
		return
	}

	// Write our current user item to local storage and redirect to the dashboard
	if _, err := fmt.Fprintf(
		w,
		`<html><head></head><body><script>window.opener.postMessage({token: '%s'}, '*');</script></body></html>`,
		token,
	); err != nil {
		log.Error("error occurred attempting to write jwt to client ", err)
	}
}

func (handler *authHandler) callback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		api.HandleError(w, http.StatusInternalServerError, errors.Wrap(err, "error completing oAuth sign in"))
	} else {
		handler.oAuthComplete(w, user)
	}
}

func (handler *authHandler) login(w http.ResponseWriter, r *http.Request) {
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil { // try to get the user without re-authenticating
		handler.oAuthComplete(w, gothUser)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (handler *authHandler) profile(r *http.Request) api.Response {
	claims, err := handler.serializer.Decode(r)
	if err != nil {
		return api.BadRequest(err, "error getting jwt claims")
	}
	user, err := handler.userStore.GetByUsername(claims.Username)
	if err != nil {
		if err == store.ErrorNotFound {
			return api.NotFound()
		}
		return api.InternalServerError(errors.Wrap(err, "error getting user"))
	}
	return api.Ok(user)
}

func (handler *authHandler) get(r *http.Request) api.Response {
	username := chi.URLParam(r, "username")
	user, err := handler.userStore.GetByUsername(username)
	if err != nil {
		if err == store.ErrorNotFound {
			return api.NotFound()
		}
		return api.InternalServerError(errors.Wrap(err, "error getting user"))
	}
	return api.Ok(user)
}

func (handler *authHandler) list(r *http.Request) api.Response {
	filter := r.URL.Query().Get("filter")
	users, err := handler.userStore.Filter(filter)
	if err != nil {
		return api.InternalServerError(errors.Wrap(err, "error getting users"))
	}
	return api.Ok(users)
}

type userUpdateRequest struct {
	Role security.UserRole
}

func (r *userUpdateRequest) IsValid() error {
	if r.Role != security.UserRoleAdmin && r.Role != security.UserRoleReader {
		return fmt.Errorf("unknown fole: %s", r.Role)
	}

	return nil
}

func (handler *authHandler) setRole(r *http.Request) api.Response {
	username := chi.URLParam(r, "username")
	user, err := handler.userStore.GetByUsername(username)
	if err != nil {
		if err == store.ErrorNotFound {
			return api.NotFound()
		}
		return api.InternalServerError(errors.Wrap(err, "error getting user"))
	}

	updateRequest := userUpdateRequest{}
	if e := json.NewDecoder(r.Body).Decode(&updateRequest); e != nil {
		return api.BadRequest(e, "bad request data")
	}

	if e := updateRequest.IsValid(); e != nil {
		return api.BadRequest(errors.Wrap(e, "error validating user update"), "invalid request data")
	}

	user.Role = updateRequest.Role
	if _, e := handler.userStore.AddOrUpdate(*user); e != nil {
		return api.InternalServerError(errors.Wrap(e, "error saving user"))
	}

	return api.Ok(user)
}

func Routes(
	defaultAdminUser string,
	userStore store.UserStore,
	oauthProviders []goth.Provider,
	serializer security.TokenSerializer,
) *chi.Mux {
	for _, p := range oauthProviders {
		log.Info("registering oauth provider ", p.Name())
		goth.UseProviders(p)
	}
	handler := authHandler{
		userStore:        userStore,
		serializer:       serializer,
		defaultAdminUser: defaultAdminUser,
	}

	router := chi.NewRouter()
	router.Get("/login", handler.login)
	router.Get("/callback", handler.callback)
	router.Get("/profile", api.Handle(handler.profile))
	router.Get("/profile/{username}", api.Handle(handler.get))
	router.Post("/profile/{username}", api.Handle(handler.setRole))
	router.Get("/", api.Handle(handler.list))
	return router
}
