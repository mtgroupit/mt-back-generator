// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/mtgroupit/mt-back-generator/extauthapi/models"
)

// IsUsernameAvailableReader is a Reader for the IsUsernameAvailable structure.
type IsUsernameAvailableReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *IsUsernameAvailableReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewIsUsernameAvailableOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewIsUsernameAvailableDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewIsUsernameAvailableOK creates a IsUsernameAvailableOK with default headers values
func NewIsUsernameAvailableOK() *IsUsernameAvailableOK {
	return &IsUsernameAvailableOK{}
}

/*IsUsernameAvailableOK handles this case with default header values.

OK
*/
type IsUsernameAvailableOK struct {
	Payload *IsUsernameAvailableOKBody
}

func (o *IsUsernameAvailableOK) Error() string {
	return fmt.Sprintf("[POST /is-username-available][%d] isUsernameAvailableOK  %+v", 200, o.Payload)
}

func (o *IsUsernameAvailableOK) GetPayload() *IsUsernameAvailableOKBody {
	return o.Payload
}

func (o *IsUsernameAvailableOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(IsUsernameAvailableOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewIsUsernameAvailableDefault creates a IsUsernameAvailableDefault with default headers values
func NewIsUsernameAvailableDefault(code int) *IsUsernameAvailableDefault {
	return &IsUsernameAvailableDefault{
		_statusCode: code,
	}
}

/*IsUsernameAvailableDefault handles this case with default header values.

Generic error response.
*/
type IsUsernameAvailableDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the is username available default response
func (o *IsUsernameAvailableDefault) Code() int {
	return o._statusCode
}

func (o *IsUsernameAvailableDefault) Error() string {
	return fmt.Sprintf("[POST /is-username-available][%d] isUsernameAvailable default  %+v", o._statusCode, o.Payload)
}

func (o *IsUsernameAvailableDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *IsUsernameAvailableDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*IsUsernameAvailableBody is username available body
swagger:model IsUsernameAvailableBody
*/
type IsUsernameAvailableBody struct {

	// username
	// Required: true
	Username models.Username `json:"username"`
}

// Validate validates this is username available body
func (o *IsUsernameAvailableBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateUsername(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *IsUsernameAvailableBody) validateUsername(formats strfmt.Registry) error {

	if err := o.Username.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "username")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *IsUsernameAvailableBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *IsUsernameAvailableBody) UnmarshalBinary(b []byte) error {
	var res IsUsernameAvailableBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*IsUsernameAvailableOKBody is username available o k body
swagger:model IsUsernameAvailableOKBody
*/
type IsUsernameAvailableOKBody struct {

	// True if username is available.
	// Required: true
	Available *bool `json:"available"`
}

// Validate validates this is username available o k body
func (o *IsUsernameAvailableOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateAvailable(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *IsUsernameAvailableOKBody) validateAvailable(formats strfmt.Registry) error {

	if err := validate.Required("isUsernameAvailableOK"+"."+"available", "body", o.Available); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *IsUsernameAvailableOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *IsUsernameAvailableOKBody) UnmarshalBinary(b []byte) error {
	var res IsUsernameAvailableOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
