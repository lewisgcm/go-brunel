package security

import (
	"go-brunel/internal/pkg/server/endpoint/api"
	"net/http"

	"github.com/casbin/casbin"
	"github.com/pkg/errors"
)

// Middleware takes casbin key match and policy file for enforcing security authorization for users accessing http routes.
func Middleware(keyMatchFile string, policyFile string, serializer TokenSerializer) func(next http.Handler) http.Handler {
	e := casbin.NewEnforcer(keyMatchFile, policyFile)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				identity, err := serializer.Decode(r)
				if err != nil {
					api.HandleError(w, http.StatusUnauthorized, errors.Wrap(err, "error parsing jwt token"))
					return
				}

				subject := ""
				if identity != nil {
					subject = identity.Username
					e.AddRoleForUser(subject, string(identity.Role))
					defer e.DeleteRolesForUser(subject)
				}

				if e.Enforce(subject, r.URL.Path, r.Method) {
					next.ServeHTTP(w, r)
					return
				}
				api.HandleError(w, http.StatusUnauthorized, errors.New("unauthorized access"))
			},
		)
	}
}
