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

// SetUsernameHandlerFunc turns a function with the right signature into a set username handler
type SetUsernameHandlerFunc func(SetUsernameParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn SetUsernameHandlerFunc) Handle(params SetUsernameParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// SetUsernameHandler interface for that can handle valid set username params
type SetUsernameHandler interface {
	Handle(SetUsernameParams, interface{}) middleware.Responder
}

// NewSetUsername creates a new http.Handler for the set username operation
func NewSetUsername(ctx *middleware.Context, handler SetUsernameHandler) *SetUsername {
	return &SetUsername{Context: ctx, Handler: handler}
}

/*SetUsername swagger:route POST /set-username setUsername

Set user username.

*/
type SetUsername struct {
	Context *middleware.Context
	Handler SetUsernameHandler
}

func (o *SetUsername) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewSetUsernameParams()

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

// SetUsernameBody set username body
//
// swagger:model SetUsernameBody
type SetUsernameBody struct {

	// username
	// Required: true
	Username models.Username `json:"username"`
}

// Validate validates this set username body
func (o *SetUsernameBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateUsername(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SetUsernameBody) validateUsername(formats strfmt.Registry) error {

	if err := o.Username.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "username")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SetUsernameBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SetUsernameBody) UnmarshalBinary(b []byte) error {
	var res SetUsernameBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
