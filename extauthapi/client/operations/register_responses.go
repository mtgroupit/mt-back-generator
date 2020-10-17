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

// RegisterReader is a Reader for the Register structure.
type RegisterReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RegisterReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewRegisterOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewRegisterDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewRegisterOK creates a RegisterOK with default headers values
func NewRegisterOK() *RegisterOK {
	return &RegisterOK{}
}

/*RegisterOK handles this case with default header values.

OK
*/
type RegisterOK struct {
	/*Session token.
	 */
	SetCookie string

	Payload *RegisterOKBody
}

func (o *RegisterOK) Error() string {
	return fmt.Sprintf("[POST /register][%d] registerOK  %+v", 200, o.Payload)
}

func (o *RegisterOK) GetPayload() *RegisterOKBody {
	return o.Payload
}

func (o *RegisterOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response header Set-Cookie
	o.SetCookie = response.GetHeader("Set-Cookie")

	o.Payload = new(RegisterOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewRegisterDefault creates a RegisterDefault with default headers values
func NewRegisterDefault(code int) *RegisterDefault {
	return &RegisterDefault{
		_statusCode: code,
	}
}

/*RegisterDefault handles this case with default header values.

- 409.700: email is not available
- 422.702: password is too weak
- 409.709: invalid email validation token

*/
type RegisterDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the register default response
func (o *RegisterDefault) Code() int {
	return o._statusCode
}

func (o *RegisterDefault) Error() string {
	return fmt.Sprintf("[POST /register][%d] register default  %+v", o._statusCode, o.Payload)
}

func (o *RegisterDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *RegisterDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*RegisterBody register body
swagger:model RegisterBody
*/
type RegisterBody struct {

	// email token
	// Required: true
	EmailToken models.JWT `json:"emailToken"`

	// password
	// Required: true
	Password models.Password `json:"password"`
}

// Validate validates this register body
func (o *RegisterBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateEmailToken(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validatePassword(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *RegisterBody) validateEmailToken(formats strfmt.Registry) error {

	if err := o.EmailToken.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "emailToken")
		}
		return err
	}

	return nil
}

func (o *RegisterBody) validatePassword(formats strfmt.Registry) error {

	if err := o.Password.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "password")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *RegisterBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *RegisterBody) UnmarshalBinary(b []byte) error {
	var res RegisterBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*RegisterOKBody register o k body
swagger:model RegisterOKBody
*/
type RegisterOKBody struct {

	// user ID
	// Required: true
	// Format: uuid4
	UserID models.UserID `json:"userID"`
}

// Validate validates this register o k body
func (o *RegisterOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateUserID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *RegisterOKBody) validateUserID(formats strfmt.Registry) error {

	if err := o.UserID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("registerOK" + "." + "userID")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *RegisterOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *RegisterOKBody) UnmarshalBinary(b []byte) error {
	var res RegisterOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}