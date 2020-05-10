package user

import (
	"fmt"
	"go-brunel/internal/pkg/server/endpoint/api"
	"go-brunel/internal/pkg/server/security"
	"go-brunel/internal/pkg/server/store"
	"net/http"

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

func (handler *authHandler) oAuthComplete(w http.ResponseWriter, user goth.User) {
	roles, err := handler.userStore.AddOrUpdate(store.User{
		Username:  user.Email,
		Email:     user.Email,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
	})
	if err != nil {
		log.Error("error occurred updating user ", err)
		api.HandleError(w, http.StatusInternalServerError, errors.New("error authenticating"))
		return
	}

	if roles.Username == handler.defaultAdminUser {
		roles.Role = security.UserRoleAdmin
	}

	token, err := handler.serializer.Encode(security.Identity{Username: roles.Username, Role: roles.Role})
	if err != nil {
		log.Error("error occurred creating jwt token ", err)
		api.HandleError(w, http.StatusInternalServerError, errors.New("error authenticating"))
		return
	}

	// Write our current user item to local storage and redirect to the dashboard
	// TODO this could be nicer: https://stackoverflow.com/questions/9153445/how-to-communicate-between-iframe-and-the-parent-site
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
	return router
}
