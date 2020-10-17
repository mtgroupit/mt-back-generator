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

	"github.com/mtgroupit/mt-back-generator/extauthapi/models"
)

// ResetPasswordHandlerFunc turns a function with the right signature into a reset password handler
type ResetPasswordHandlerFunc func(ResetPasswordParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ResetPasswordHandlerFunc) Handle(params ResetPasswordParams) middleware.Responder {
	return fn(params)
}

// ResetPasswordHandler interface for that can handle valid reset password params
type ResetPasswordHandler interface {
	Handle(ResetPasswordParams) middleware.Responder
}

// NewResetPassword creates a new http.Handler for the reset password operation
func NewResetPassword(ctx *middleware.Context, handler ResetPasswordHandler) *ResetPassword {
	return &ResetPassword{Context: ctx, Handler: handler}
}

/*ResetPassword swagger:route POST /reset-password resetPassword

Request password reset by  email.

*/
type ResetPassword struct {
	Context *middleware.Context
	Handler ResetPasswordHandler
}

func (o *ResetPassword) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewResetPasswordParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// ResetPasswordBody reset password body
//
// swagger:model ResetPasswordBody
type ResetPasswordBody struct {

	// email
	// Required: true
	// Format: email
	Email models.Email `json:"email"`
}

// Validate validates this reset password body
func (o *ResetPasswordBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateEmail(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ResetPasswordBody) validateEmail(formats strfmt.Registry) error {

	if err := o.Email.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "email")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *ResetPasswordBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ResetPasswordBody) UnmarshalBinary(b []byte) error {
	var res ResetPasswordBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
