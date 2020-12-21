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

	"github.com/mtgroupit/mt-back-generator/extauthapi/models"
)

// ResetPasswordReader is a Reader for the ResetPassword structure.
type ResetPasswordReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ResetPasswordReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewResetPasswordNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewResetPasswordDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewResetPasswordNoContent creates a ResetPasswordNoContent with default headers values
func NewResetPasswordNoContent() *ResetPasswordNoContent {
	return &ResetPasswordNoContent{}
}

/*ResetPasswordNoContent handles this case with default header values.

The server successfully processed the request and is not returning any content.
*/
type ResetPasswordNoContent struct {
}

func (o *ResetPasswordNoContent) Error() string {
	return fmt.Sprintf("[POST /reset-password][%d] resetPasswordNoContent ", 204)
}

func (o *ResetPasswordNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewResetPasswordDefault creates a ResetPasswordDefault with default headers values
func NewResetPasswordDefault(code int) *ResetPasswordDefault {
	return &ResetPasswordDefault{
		_statusCode: code,
	}
}

/*ResetPasswordDefault handles this case with default header values.

- 404.101: invalid credentials
- 404.707: no such email

*/
type ResetPasswordDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the reset password default response
func (o *ResetPasswordDefault) Code() int {
	return o._statusCode
}

func (o *ResetPasswordDefault) Error() string {
	return fmt.Sprintf("[POST /reset-password][%d] resetPassword default  %+v", o._statusCode, o.Payload)
}

func (o *ResetPasswordDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *ResetPasswordDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*ResetPasswordBody reset password body
swagger:model ResetPasswordBody
*/
type ResetPasswordBody struct {

	// email
	// Required: true
	// Format: email
	Email models.Email `json:"email"`
}

// Validate validates this reset password body
func (o *ResetPasswordBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateEmail(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ResetPasswordBody) validateEmail(formats strfmt.Registry) error {

	if err := o.Email.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "email")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *ResetPasswordBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ResetPasswordBody) UnmarshalBinary(b []byte) error {
	var res ResetPasswordBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
