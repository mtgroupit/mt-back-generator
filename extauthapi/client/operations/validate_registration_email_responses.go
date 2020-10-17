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

// ValidateRegistrationEmailReader is a Reader for the ValidateRegistrationEmail structure.
type ValidateRegistrationEmailReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ValidateRegistrationEmailReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewValidateRegistrationEmailNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewValidateRegistrationEmailDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewValidateRegistrationEmailNoContent creates a ValidateRegistrationEmailNoContent with default headers values
func NewValidateRegistrationEmailNoContent() *ValidateRegistrationEmailNoContent {
	return &ValidateRegistrationEmailNoContent{}
}

/*ValidateRegistrationEmailNoContent handles this case with default header values.

The server successfully processed the request and is not returning any content.
*/
type ValidateRegistrationEmailNoContent struct {
}

func (o *ValidateRegistrationEmailNoContent) Error() string {
	return fmt.Sprintf("[POST /validate-registration-email][%d] validateRegistrationEmailNoContent ", 204)
}

func (o *ValidateRegistrationEmailNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewValidateRegistrationEmailDefault creates a ValidateRegistrationEmailDefault with default headers values
func NewValidateRegistrationEmailDefault(code int) *ValidateRegistrationEmailDefault {
	return &ValidateRegistrationEmailDefault{
		_statusCode: code,
	}
}

/*ValidateRegistrationEmailDefault handles this case with default header values.

- 409.700: email is not available
- 422.703: wrong captcha answer

*/
type ValidateRegistrationEmailDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the validate registration email default response
func (o *ValidateRegistrationEmailDefault) Code() int {
	return o._statusCode
}

func (o *ValidateRegistrationEmailDefault) Error() string {
	return fmt.Sprintf("[POST /validate-registration-email][%d] validateRegistrationEmail default  %+v", o._statusCode, o.Payload)
}

func (o *ValidateRegistrationEmailDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *ValidateRegistrationEmailDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*ValidateRegistrationEmailBody validate registration email body
swagger:model ValidateRegistrationEmailBody
*/
type ValidateRegistrationEmailBody struct {

	// email
	// Required: true
	// Format: email
	Email models.Email `json:"email"`
}

// Validate validates this validate registration email body
func (o *ValidateRegistrationEmailBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateEmail(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ValidateRegistrationEmailBody) validateEmail(formats strfmt.Registry) error {

	if err := o.Email.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "email")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *ValidateRegistrationEmailBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ValidateRegistrationEmailBody) UnmarshalBinary(b []byte) error {
	var res ValidateRegistrationEmailBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
