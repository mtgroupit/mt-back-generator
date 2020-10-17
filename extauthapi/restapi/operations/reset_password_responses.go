// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/mtgroupit/mt-back-generator/extauthapi/models"
)

// ResetPasswordNoContentCode is the HTTP code returned for type ResetPasswordNoContent
const ResetPasswordNoContentCode int = 204

/*ResetPasswordNoContent The server successfully processed the request and is not returning any content.

swagger:response resetPasswordNoContent
*/
type ResetPasswordNoContent struct {
}

// NewResetPasswordNoContent creates ResetPasswordNoContent with default headers values
func NewResetPasswordNoContent() *ResetPasswordNoContent {

	return &ResetPasswordNoContent{}
}

// WriteResponse to the client
func (o *ResetPasswordNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

/*ResetPasswordDefault - 404.101: invalid credentials
- 404.102: organisation not found


swagger:response resetPasswordDefault
*/
type ResetPasswordDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewResetPasswordDefault creates ResetPasswordDefault with default headers values
func NewResetPasswordDefault(code int) *ResetPasswordDefault {
	if code <= 0 {
		code = 500
	}

	return &ResetPasswordDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the reset password default response
func (o *ResetPasswordDefault) WithStatusCode(code int) *ResetPasswordDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the reset password default response
func (o *ResetPasswordDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the reset password default response
func (o *ResetPasswordDefault) WithPayload(payload *models.Error) *ResetPasswordDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the reset password default response
func (o *ResetPasswordDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ResetPasswordDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
