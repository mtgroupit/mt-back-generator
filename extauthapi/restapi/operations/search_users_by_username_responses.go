// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/mtgroupit/mt-back-generator/extauthapi/models"
)

// SearchUsersByUsernameOKCode is the HTTP code returned for type SearchUsersByUsernameOK
const SearchUsersByUsernameOKCode int = 200

/*SearchUsersByUsernameOK OK

swagger:response searchUsersByUsernameOK
*/
type SearchUsersByUsernameOK struct {

	/*
	  In: Body
	*/
	Payload *SearchUsersByUsernameOKBody `json:"body,omitempty"`
}

// NewSearchUsersByUsernameOK creates SearchUsersByUsernameOK with default headers values
func NewSearchUsersByUsernameOK() *SearchUsersByUsernameOK {

	return &SearchUsersByUsernameOK{}
}

// WithPayload adds the payload to the search users by username o k response
func (o *SearchUsersByUsernameOK) WithPayload(payload *SearchUsersByUsernameOKBody) *SearchUsersByUsernameOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the search users by username o k response
func (o *SearchUsersByUsernameOK) SetPayload(payload *SearchUsersByUsernameOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SearchUsersByUsernameOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*SearchUsersByUsernameDefault Generic error response.

swagger:response searchUsersByUsernameDefault
*/
type SearchUsersByUsernameDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSearchUsersByUsernameDefault creates SearchUsersByUsernameDefault with default headers values
func NewSearchUsersByUsernameDefault(code int) *SearchUsersByUsernameDefault {
	if code <= 0 {
		code = 500
	}

	return &SearchUsersByUsernameDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the search users by username default response
func (o *SearchUsersByUsernameDefault) WithStatusCode(code int) *SearchUsersByUsernameDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the search users by username default response
func (o *SearchUsersByUsernameDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the search users by username default response
func (o *SearchUsersByUsernameDefault) WithPayload(payload *models.Error) *SearchUsersByUsernameDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the search users by username default response
func (o *SearchUsersByUsernameDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SearchUsersByUsernameDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
