// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/Lisss13/french-back-template/snouki-mobile/mtmb-extauthapi/models"
)

// LoginHandlerFunc turns a function with the right signature into a login handler
type LoginHandlerFunc func(LoginParams) middleware.Responder

// Handle executing the request and returning a response
func (fn LoginHandlerFunc) Handle(params LoginParams) middleware.Responder {
	return fn(params)
}

// LoginHandler interface for that can handle valid login params
type LoginHandler interface {
	Handle(LoginParams) middleware.Responder
}

// NewLogin creates a new http.Handler for the login operation
func NewLogin(ctx *middleware.Context, handler LoginHandler) *Login {
	return &Login{Context: ctx, Handler: handler}
}

/*Login swagger:route POST /login login

Sign in user.

*/
type Login struct {
	Context *middleware.Context
	Handler LoginHandler
}

func (o *Login) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewLoginParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// LoginBody login body
//
// swagger:model LoginBody
type LoginBody struct {

	// email
	// Required: true
	// Format: email
	Email models.Email `json:"email"`

	// password
	// Required: true
	Password models.Password `json:"password"`
}

// Validate validates this login body
func (o *LoginBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateEmail(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validatePassword(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *LoginBody) validateEmail(formats strfmt.Registry) error {

	if err := o.Email.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "email")
		}
		return err
	}

	return nil
}

func (o *LoginBody) validatePassword(formats strfmt.Registry) error {

	if err := o.Password.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "password")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *LoginBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *LoginBody) UnmarshalBinary(b []byte) error {
	var res LoginBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
