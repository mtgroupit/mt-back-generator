// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/Lisss13/french-back-template/snouki-mobile/mtmb-extauthapi/restapi/operations"
)

//go:generate swagger generate server --target ../../mtmb-extauthapi --name Authentication --spec ../swagger.yml --principal interface{} --exclude-main

func configureFlags(api *operations.AuthenticationAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AuthenticationAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "Cookie" header is set
	if api.CookieKeyAuth == nil {
		api.CookieKeyAuth = func(token string) (interface{}, error) {
			return nil, errors.NotImplemented("api key auth (cookieKey) Cookie from header param [Cookie] has not yet been implemented")
		}
	}
	// Applies when the "X-CSRFTokenBound" header is set
	if api.CsrfTokenAuth == nil {
		api.CsrfTokenAuth = func(token string) (interface{}, error) {
			return nil, errors.NotImplemented("api key auth (csrfToken) X-CSRFTokenBound from header param [X-CSRFTokenBound] has not yet been implemented")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	if api.ChangePasswordHandler == nil {
		api.ChangePasswordHandler = operations.ChangePasswordHandlerFunc(func(params operations.ChangePasswordParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.ChangePassword has not yet been implemented")
		})
	}
	if api.DeleteUserHandler == nil {
		api.DeleteUserHandler = operations.DeleteUserHandlerFunc(func(params operations.DeleteUserParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.DeleteUser has not yet been implemented")
		})
	}
	if api.GetUserProfileHandler == nil {
		api.GetUserProfileHandler = operations.GetUserProfileHandlerFunc(func(params operations.GetUserProfileParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetUserProfile has not yet been implemented")
		})
	}
	if api.GetUserProfileByIDHandler == nil {
		api.GetUserProfileByIDHandler = operations.GetUserProfileByIDHandlerFunc(func(params operations.GetUserProfileByIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetUserProfileByID has not yet been implemented")
		})
	}
	if api.IsEmailAvailableHandler == nil {
		api.IsEmailAvailableHandler = operations.IsEmailAvailableHandlerFunc(func(params operations.IsEmailAvailableParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.IsEmailAvailable has not yet been implemented")
		})
	}
	if api.IsUsernameAvailableHandler == nil {
		api.IsUsernameAvailableHandler = operations.IsUsernameAvailableHandlerFunc(func(params operations.IsUsernameAvailableParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.IsUsernameAvailable has not yet been implemented")
		})
	}
	if api.LoginHandler == nil {
		api.LoginHandler = operations.LoginHandlerFunc(func(params operations.LoginParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.Login has not yet been implemented")
		})
	}
	if api.LogoutHandler == nil {
		api.LogoutHandler = operations.LogoutHandlerFunc(func(params operations.LogoutParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.Logout has not yet been implemented")
		})
	}
	if api.RegisterHandler == nil {
		api.RegisterHandler = operations.RegisterHandlerFunc(func(params operations.RegisterParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.Register has not yet been implemented")
		})
	}
	if api.RegisterLoginOAuthHandler == nil {
		api.RegisterLoginOAuthHandler = operations.RegisterLoginOAuthHandlerFunc(func(params operations.RegisterLoginOAuthParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.RegisterLoginOAuth has not yet been implemented")
		})
	}
	if api.ResetPasswordHandler == nil {
		api.ResetPasswordHandler = operations.ResetPasswordHandlerFunc(func(params operations.ResetPasswordParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.ResetPassword has not yet been implemented")
		})
	}
	if api.SearchUsersByUsernameHandler == nil {
		api.SearchUsersByUsernameHandler = operations.SearchUsersByUsernameHandlerFunc(func(params operations.SearchUsersByUsernameParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.SearchUsersByUsername has not yet been implemented")
		})
	}
	if api.SetBlockedHandler == nil {
		api.SetBlockedHandler = operations.SetBlockedHandlerFunc(func(params operations.SetBlockedParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.SetBlocked has not yet been implemented")
		})
	}
	if api.SetEmailHandler == nil {
		api.SetEmailHandler = operations.SetEmailHandlerFunc(func(params operations.SetEmailParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.SetEmail has not yet been implemented")
		})
	}
	if api.SetNewPasswordHandler == nil {
		api.SetNewPasswordHandler = operations.SetNewPasswordHandlerFunc(func(params operations.SetNewPasswordParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.SetNewPassword has not yet been implemented")
		})
	}
	if api.SetPersDataRegionHandler == nil {
		api.SetPersDataRegionHandler = operations.SetPersDataRegionHandlerFunc(func(params operations.SetPersDataRegionParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.SetPersDataRegion has not yet been implemented")
		})
	}
	if api.SetUsernameHandler == nil {
		api.SetUsernameHandler = operations.SetUsernameHandlerFunc(func(params operations.SetUsernameParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.SetUsername has not yet been implemented")
		})
	}
	if api.ValidateNewEmailHandler == nil {
		api.ValidateNewEmailHandler = operations.ValidateNewEmailHandlerFunc(func(params operations.ValidateNewEmailParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.ValidateNewEmail has not yet been implemented")
		})
	}
	if api.ValidateRegistrationEmailHandler == nil {
		api.ValidateRegistrationEmailHandler = operations.ValidateRegistrationEmailHandlerFunc(func(params operations.ValidateRegistrationEmailParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.ValidateRegistrationEmail has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
