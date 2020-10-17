// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/Lisss13/french-back-template/snouki-mobile/mtmb-extauthapi/models"
)

// SetNewPasswordNoContentCode is the HTTP code returned for type SetNewPasswordNoContent
const SetNewPasswordNoContentCode int = 204

/*SetNewPasswordNoContent The server successfully processed the request and is not returning any content.

swagger:response setNewPasswordNoContent
*/
type SetNewPasswordNoContent struct {
}

// NewSetNewPasswordNoContent creates SetNewPasswordNoContent with default headers values
func NewSetNewPasswordNoContent() *SetNewPasswordNoContent {

	return &SetNewPasswordNoContent{}
}

// WriteResponse to the client
func (o *SetNewPasswordNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

/*SetNewPasswordDefault - 404.102: organisation not found
- 422.103: password is too weak
- 403.104: invalid password reset token


swagger:response setNewPasswordDefault
*/
type SetNewPasswordDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSetNewPasswordDefault creates SetNewPasswordDefault with default headers values
func NewSetNewPasswordDefault(code int) *SetNewPasswordDefault {
	if code <= 0 {
		code = 500
	}

	return &SetNewPasswordDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the set new password default response
func (o *SetNewPasswordDefault) WithStatusCode(code int) *SetNewPasswordDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the set new password default response
func (o *SetNewPasswordDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the set new password default response
func (o *SetNewPasswordDefault) WithPayload(payload *models.Error) *SetNewPasswordDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the set new password default response
func (o *SetNewPasswordDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SetNewPasswordDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
