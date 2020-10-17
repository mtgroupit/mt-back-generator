// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/Lisss13/french-back-template/snouki-mobile/mtmb-extauthapi/models"
)

// DeleteUserNoContentCode is the HTTP code returned for type DeleteUserNoContent
const DeleteUserNoContentCode int = 204

/*DeleteUserNoContent The server successfully processed the request and is not returning any content.

swagger:response deleteUserNoContent
*/
type DeleteUserNoContent struct {
	/*Session token.

	 */
	SetCookie string `json:"Set-Cookie"`
}

// NewDeleteUserNoContent creates DeleteUserNoContent with default headers values
func NewDeleteUserNoContent() *DeleteUserNoContent {

	return &DeleteUserNoContent{}
}

// WithSetCookie adds the setCookie to the delete user no content response
func (o *DeleteUserNoContent) WithSetCookie(setCookie string) *DeleteUserNoContent {
	o.SetCookie = setCookie
	return o
}

// SetSetCookie sets the setCookie to the delete user no content response
func (o *DeleteUserNoContent) SetSetCookie(setCookie string) {
	o.SetCookie = setCookie
}

// WriteResponse to the client
func (o *DeleteUserNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Set-Cookie

	setCookie := o.SetCookie
	if setCookie != "" {
		rw.Header().Set("Set-Cookie", setCookie)
	}

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

/*DeleteUserDefault - 403.710: invalid password


swagger:response deleteUserDefault
*/
type DeleteUserDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteUserDefault creates DeleteUserDefault with default headers values
func NewDeleteUserDefault(code int) *DeleteUserDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteUserDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete user default response
func (o *DeleteUserDefault) WithStatusCode(code int) *DeleteUserDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete user default response
func (o *DeleteUserDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the delete user default response
func (o *DeleteUserDefault) WithPayload(payload *models.Error) *DeleteUserDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete user default response
func (o *DeleteUserDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteUserDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
