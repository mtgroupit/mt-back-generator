// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/mtgroupit/mt-back-generator/extauthapi/models"
)

// ValidateNewEmailNoContentCode is the HTTP code returned for type ValidateNewEmailNoContent
const ValidateNewEmailNoContentCode int = 204

/*ValidateNewEmailNoContent The server successfully processed the request and is not returning any content.

swagger:response validateNewEmailNoContent
*/
type ValidateNewEmailNoContent struct {
}

// NewValidateNewEmailNoContent creates ValidateNewEmailNoContent with default headers values
func NewValidateNewEmailNoContent() *ValidateNewEmailNoContent {

	return &ValidateNewEmailNoContent{}
}

// WriteResponse to the client
func (o *ValidateNewEmailNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

/*ValidateNewEmailDefault - 404.2003: email is not available
- 403.710: invalid password


swagger:response validateNewEmailDefault
*/
type ValidateNewEmailDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewValidateNewEmailDefault creates ValidateNewEmailDefault with default headers values
func NewValidateNewEmailDefault(code int) *ValidateNewEmailDefault {
	if code <= 0 {
		code = 500
	}

	return &ValidateNewEmailDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the validate new email default response
func (o *ValidateNewEmailDefault) WithStatusCode(code int) *ValidateNewEmailDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the validate new email default response
func (o *ValidateNewEmailDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the validate new email default response
func (o *ValidateNewEmailDefault) WithPayload(payload *models.Error) *ValidateNewEmailDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the validate new email default response
func (o *ValidateNewEmailDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ValidateNewEmailDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
