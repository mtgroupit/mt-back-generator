// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/mtgroupit/mt-back-generator/extauthapi/models"
)

// SetEmailNoContentCode is the HTTP code returned for type SetEmailNoContent
const SetEmailNoContentCode int = 204

/*SetEmailNoContent The server successfully processed the request and is not returning any content.

swagger:response setEmailNoContent
*/
type SetEmailNoContent struct {
}

// NewSetEmailNoContent creates SetEmailNoContent with default headers values
func NewSetEmailNoContent() *SetEmailNoContent {

	return &SetEmailNoContent{}
}

// WriteResponse to the client
func (o *SetEmailNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

/*SetEmailDefault - 409.700: email is not available
- 409.709: invalid email validation token


swagger:response setEmailDefault
*/
type SetEmailDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSetEmailDefault creates SetEmailDefault with default headers values
func NewSetEmailDefault(code int) *SetEmailDefault {
	if code <= 0 {
		code = 500
	}

	return &SetEmailDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the set email default response
func (o *SetEmailDefault) WithStatusCode(code int) *SetEmailDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the set email default response
func (o *SetEmailDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the set email default response
func (o *SetEmailDefault) WithPayload(payload *models.Error) *SetEmailDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the set email default response
func (o *SetEmailDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SetEmailDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}