// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/mtgroupit/mt-back-generator/extauthapi/models"
)

// RegisterOKCode is the HTTP code returned for type RegisterOK
const RegisterOKCode int = 200

/*RegisterOK OK

swagger:response registerOK
*/
type RegisterOK struct {
	/*Session token.

	 */
	SetCookie string `json:"Set-Cookie"`

	/*
	  In: Body
	*/
	Payload *RegisterOKBody `json:"body,omitempty"`
}

// NewRegisterOK creates RegisterOK with default headers values
func NewRegisterOK() *RegisterOK {

	return &RegisterOK{}
}

// WithSetCookie adds the setCookie to the register o k response
func (o *RegisterOK) WithSetCookie(setCookie string) *RegisterOK {
	o.SetCookie = setCookie
	return o
}

// SetSetCookie sets the setCookie to the register o k response
func (o *RegisterOK) SetSetCookie(setCookie string) {
	o.SetCookie = setCookie
}

// WithPayload adds the payload to the register o k response
func (o *RegisterOK) WithPayload(payload *RegisterOKBody) *RegisterOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the register o k response
func (o *RegisterOK) SetPayload(payload *RegisterOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RegisterOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Set-Cookie

	setCookie := o.SetCookie
	if setCookie != "" {
		rw.Header().Set("Set-Cookie", setCookie)
	}

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*RegisterDefault - 404.2003: email is not available
- 404.2001: invalid email validation token


swagger:response registerDefault
*/
type RegisterDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewRegisterDefault creates RegisterDefault with default headers values
func NewRegisterDefault(code int) *RegisterDefault {
	if code <= 0 {
		code = 500
	}

	return &RegisterDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the register default response
func (o *RegisterDefault) WithStatusCode(code int) *RegisterDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the register default response
func (o *RegisterDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the register default response
func (o *RegisterDefault) WithPayload(payload *models.Error) *RegisterDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the register default response
func (o *RegisterDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RegisterDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
