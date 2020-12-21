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

// RequestRegistrationHandlerFunc turns a function with the right signature into a request registration handler
type RequestRegistrationHandlerFunc func(RequestRegistrationParams) middleware.Responder

// Handle executing the request and returning a response
func (fn RequestRegistrationHandlerFunc) Handle(params RequestRegistrationParams) middleware.Responder {
	return fn(params)
}

// RequestRegistrationHandler interface for that can handle valid request registration params
type RequestRegistrationHandler interface {
	Handle(RequestRegistrationParams) middleware.Responder
}

// NewRequestRegistration creates a new http.Handler for the request registration operation
func NewRequestRegistration(ctx *middleware.Context, handler RequestRegistrationHandler) *RequestRegistration {
	return &RequestRegistration{Context: ctx, Handler: handler}
}

/*RequestRegistration swagger:route POST /request-registration requestRegistration

Sends email with validation token.

*/
type RequestRegistration struct {
	Context *middleware.Context
	Handler RequestRegistrationHandler
}

func (o *RequestRegistration) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewRequestRegistrationParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// RequestRegistrationBody request registration body
//
// swagger:model RequestRegistrationBody
type RequestRegistrationBody struct {

	// email
	// Required: true
	// Format: email
	Email models.Email `json:"email"`

	// language
	Language models.Language `json:"language,omitempty"`

	// password
	// Required: true
	Password models.Password `json:"password"`
}

// Validate validates this request registration body
func (o *RequestRegistrationBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateEmail(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateLanguage(formats); err != nil {
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

func (o *RequestRegistrationBody) validateEmail(formats strfmt.Registry) error {

	if err := o.Email.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "email")
		}
		return err
	}

	return nil
}

func (o *RequestRegistrationBody) validateLanguage(formats strfmt.Registry) error {

	if swag.IsZero(o.Language) { // not required
		return nil
	}

	if err := o.Language.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "language")
		}
		return err
	}

	return nil
}

func (o *RequestRegistrationBody) validatePassword(formats strfmt.Registry) error {

	if err := o.Password.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "password")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *RequestRegistrationBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *RequestRegistrationBody) UnmarshalBinary(b []byte) error {
	var res RequestRegistrationBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
