// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/Lisss13/french-back-template/snouki-mobile/mtmb-extauthapi/models"
)

// SetBlockedNoContentCode is the HTTP code returned for type SetBlockedNoContent
const SetBlockedNoContentCode int = 204

/*SetBlockedNoContent The server successfully processed the request and is not returning any content.

swagger:response setBlockedNoContent
*/
type SetBlockedNoContent struct {
}

// NewSetBlockedNoContent creates SetBlockedNoContent with default headers values
func NewSetBlockedNoContent() *SetBlockedNoContent {

	return &SetBlockedNoContent{}
}

// WriteResponse to the client
func (o *SetBlockedNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

/*SetBlockedDefault Generic error response.

swagger:response setBlockedDefault
*/
type SetBlockedDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSetBlockedDefault creates SetBlockedDefault with default headers values
func NewSetBlockedDefault(code int) *SetBlockedDefault {
	if code <= 0 {
		code = 500
	}

	return &SetBlockedDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the set blocked default response
func (o *SetBlockedDefault) WithStatusCode(code int) *SetBlockedDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the set blocked default response
func (o *SetBlockedDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the set blocked default response
func (o *SetBlockedDefault) WithPayload(payload *models.Error) *SetBlockedDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the set blocked default response
func (o *SetBlockedDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SetBlockedDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
