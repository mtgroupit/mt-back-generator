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

// DeleteUserHandlerFunc turns a function with the right signature into a delete user handler
type DeleteUserHandlerFunc func(DeleteUserParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteUserHandlerFunc) Handle(params DeleteUserParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeleteUserHandler interface for that can handle valid delete user params
type DeleteUserHandler interface {
	Handle(DeleteUserParams, interface{}) middleware.Responder
}

// NewDeleteUser creates a new http.Handler for the delete user operation
func NewDeleteUser(ctx *middleware.Context, handler DeleteUserHandler) *DeleteUser {
	return &DeleteUser{Context: ctx, Handler: handler}
}

/*DeleteUser swagger:route POST /delete-user deleteUser

Removes user account and expire all user sessions including current one.

*/
type DeleteUser struct {
	Context *middleware.Context
	Handler DeleteUserHandler
}

func (o *DeleteUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteUserParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// DeleteUserBody delete user body
//
// swagger:model DeleteUserBody
type DeleteUserBody struct {

	// password
	// Required: true
	Password models.Password `json:"password"`
}

// Validate validates this delete user body
func (o *DeleteUserBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validatePassword(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *DeleteUserBody) validatePassword(formats strfmt.Registry) error {

	if err := o.Password.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "password")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *DeleteUserBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *DeleteUserBody) UnmarshalBinary(b []byte) error {
	var res DeleteUserBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
