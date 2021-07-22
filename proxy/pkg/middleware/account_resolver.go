package middleware

import (
	"errors"
	"net/http"

	"github.com/cs3org/reva/pkg/auth/scope"
	"github.com/owncloud/ocis/proxy/pkg/user/backend"

	tokenPkg "github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	revauser "github.com/cs3org/reva/pkg/user"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
)

// AccountResolver provides a middleware which mints a jwt and adds it to the proxied request based
// on the oidc-claims
func AccountResolver(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)
	logger := options.Logger

	return func(next http.Handler) http.Handler {
		tokenManager, err := jwt.New(map[string]interface{}{
			"secret":  options.TokenManagerConfig.JWTSecret,
			"expires": int64(60),
		})
		if err != nil {
			logger.Fatal().Err(err).Msgf("Could not initialize token-manager")
		}

		return &accountResolver{
			next:                  next,
			logger:                logger,
			tokenManager:          tokenManager,
			userProvider:          options.UserProvider,
			userOIDCClaim:         options.UserOIDCClaim,
			userCS3Claim:          options.UserCS3Claim,
			autoProvisionAccounts: options.AutoprovisionAccounts,
		}
	}
}

type accountResolver struct {
	next                  http.Handler
	logger                log.Logger
	tokenManager          tokenPkg.Manager
	userProvider          backend.UserBackend
	autoProvisionAccounts bool
	userOIDCClaim         string
	userCS3Claim          string
}

// TODO do not use the context to store values: https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39
func (m accountResolver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	claims := oidc.FromContext(req.Context())
	u, ok := revauser.ContextGetUser(req.Context())

	if claims == nil && !ok {
		m.next.ServeHTTP(w, req)
		return
	}

	if u == nil && claims != nil {

		var err error
		var value string
		var ok bool
		if value, ok = claims[m.userOIDCClaim].(string); !ok || value == "" {
			m.logger.Error().Str("claim", m.userOIDCClaim).Interface("claims", claims).Msg("claim not set or empty")
			w.WriteHeader(http.StatusInternalServerError) // admin needs to make the idp send the right claim
			return
		}

		u, err = m.userProvider.GetUserByClaims(req.Context(), m.userCS3Claim, value, true)

		if errors.Is(err, backend.ErrAccountNotFound) {
			m.logger.Debug().Str("claim", m.userOIDCClaim).Str("value", value).Msg("User by claim not found")
			if !m.autoProvisionAccounts {
				m.logger.Debug().Interface("claims", claims).Msg("Autoprovisioning disabled")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			m.logger.Debug().Interface("claims", claims).Msg("Autoprovisioning user")
			u, err = m.userProvider.CreateUserFromClaims(req.Context(), claims)
		}

		if errors.Is(err, backend.ErrAccountDisabled) {
			m.logger.Debug().Interface("claims", claims).Msg("Disabled")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err != nil {
			m.logger.Error().Err(err).Msg("Could not get user by claim")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		m.logger.Debug().Interface("claims", claims).Interface("user", u).Msg("associated claims with user")
	}

	s, err := scope.AddOwnerScope(nil)
	if err != nil {
		m.logger.Error().Err(err).Msg("could not get owner scope")
		return
	}
	token, err := m.tokenManager.MintToken(req.Context(), u, s)
	if err != nil {
		m.logger.Error().Err(err).Msg("could not mint token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Header.Set(tokenPkg.TokenHeader, token)

	m.next.ServeHTTP(w, req)
}
