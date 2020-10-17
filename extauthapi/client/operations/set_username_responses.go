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

	"github.com/Lisss13/french-back-template/snouki-mobile/mtmb-extauthapi/models"
)

// SetUsernameReader is a Reader for the SetUsername structure.
type SetUsernameReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SetUsernameReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewSetUsernameNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewSetUsernameDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewSetUsernameNoContent creates a SetUsernameNoContent with default headers values
func NewSetUsernameNoContent() *SetUsernameNoContent {
	return &SetUsernameNoContent{}
}

/*SetUsernameNoContent handles this case with default header values.

The server successfully processed the request and is not returning any content.
*/
type SetUsernameNoContent struct {
}

func (o *SetUsernameNoContent) Error() string {
	return fmt.Sprintf("[POST /set-username][%d] setUsernameNoContent ", 204)
}

func (o *SetUsernameNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewSetUsernameDefault creates a SetUsernameDefault with default headers values
func NewSetUsernameDefault(code int) *SetUsernameDefault {
	return &SetUsernameDefault{
		_statusCode: code,
	}
}

/*SetUsernameDefault handles this case with default header values.

- 409.701: username is not available

*/
type SetUsernameDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the set username default response
func (o *SetUsernameDefault) Code() int {
	return o._statusCode
}

func (o *SetUsernameDefault) Error() string {
	return fmt.Sprintf("[POST /set-username][%d] setUsername default  %+v", o._statusCode, o.Payload)
}

func (o *SetUsernameDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *SetUsernameDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*SetUsernameBody set username body
swagger:model SetUsernameBody
*/
type SetUsernameBody struct {

	// username
	// Required: true
	Username models.Username `json:"username"`
}

// Validate validates this set username body
func (o *SetUsernameBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateUsername(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SetUsernameBody) validateUsername(formats strfmt.Registry) error {

	if err := o.Username.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "username")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SetUsernameBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SetUsernameBody) UnmarshalBinary(b []byte) error {
	var res SetUsernameBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
