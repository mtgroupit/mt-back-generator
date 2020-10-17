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

// RegisterLoginOAuthReader is a Reader for the RegisterLoginOAuth structure.
type RegisterLoginOAuthReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RegisterLoginOAuthReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewRegisterLoginOAuthOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewRegisterLoginOAuthDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewRegisterLoginOAuthOK creates a RegisterLoginOAuthOK with default headers values
func NewRegisterLoginOAuthOK() *RegisterLoginOAuthOK {
	return &RegisterLoginOAuthOK{}
}

/*RegisterLoginOAuthOK handles this case with default header values.

OK
*/
type RegisterLoginOAuthOK struct {
	/*Session token.
	 */
	SetCookie string

	Payload *RegisterLoginOAuthOKBody
}

func (o *RegisterLoginOAuthOK) Error() string {
	return fmt.Sprintf("[POST /register-login-oauth][%d] registerLoginOAuthOK  %+v", 200, o.Payload)
}

func (o *RegisterLoginOAuthOK) GetPayload() *RegisterLoginOAuthOKBody {
	return o.Payload
}

func (o *RegisterLoginOAuthOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response header Set-Cookie
	o.SetCookie = response.GetHeader("Set-Cookie")

	o.Payload = new(RegisterLoginOAuthOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewRegisterLoginOAuthDefault creates a RegisterLoginOAuthDefault with default headers values
func NewRegisterLoginOAuthDefault(code int) *RegisterLoginOAuthDefault {
	return &RegisterLoginOAuthDefault{
		_statusCode: code,
	}
}

/*RegisterLoginOAuthDefault handles this case with default header values.

- 502.713: failed to get user profile from oauth server
- 403.714: state does not match
- 403.715: user is blocked

*/
type RegisterLoginOAuthDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the register login o auth default response
func (o *RegisterLoginOAuthDefault) Code() int {
	return o._statusCode
}

func (o *RegisterLoginOAuthDefault) Error() string {
	return fmt.Sprintf("[POST /register-login-oauth][%d] registerLoginOAuth default  %+v", o._statusCode, o.Payload)
}

func (o *RegisterLoginOAuthDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *RegisterLoginOAuthDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*RegisterLoginOAuthBody register login o auth body
swagger:model RegisterLoginOAuthBody
*/
type RegisterLoginOAuthBody struct {

	// code
	// Required: true
	Code models.OAuthCode `json:"code"`

	// recv state
	// Required: true
	RecvState models.OAuthState `json:"recvState"`

	// sent state
	// Required: true
	SentState models.OAuthState `json:"sentState"`

	// server
	// Required: true
	Server models.OAuthServer `json:"server"`
}

// Validate validates this register login o auth body
func (o *RegisterLoginOAuthBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateCode(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateRecvState(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateSentState(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateServer(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *RegisterLoginOAuthBody) validateCode(formats strfmt.Registry) error {

	if err := o.Code.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "code")
		}
		return err
	}

	return nil
}

func (o *RegisterLoginOAuthBody) validateRecvState(formats strfmt.Registry) error {

	if err := o.RecvState.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "recvState")
		}
		return err
	}

	return nil
}

func (o *RegisterLoginOAuthBody) validateSentState(formats strfmt.Registry) error {

	if err := o.SentState.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "sentState")
		}
		return err
	}

	return nil
}

func (o *RegisterLoginOAuthBody) validateServer(formats strfmt.Registry) error {

	if err := o.Server.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "server")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *RegisterLoginOAuthBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *RegisterLoginOAuthBody) UnmarshalBinary(b []byte) error {
	var res RegisterLoginOAuthBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*RegisterLoginOAuthOKBody register login o auth o k body
swagger:model RegisterLoginOAuthOKBody
*/
type RegisterLoginOAuthOKBody struct {

	// user ID
	// Required: true
	// Format: uuid4
	UserID models.UserID `json:"userID"`
}

// Validate validates this register login o auth o k body
func (o *RegisterLoginOAuthOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateUserID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *RegisterLoginOAuthOKBody) validateUserID(formats strfmt.Registry) error {

	if err := o.UserID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("registerLoginOAuthOK" + "." + "userID")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *RegisterLoginOAuthOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *RegisterLoginOAuthOKBody) UnmarshalBinary(b []byte) error {
	var res RegisterLoginOAuthOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}