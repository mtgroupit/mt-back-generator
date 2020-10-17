// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/mtgroupit/mt-back-generator/extauthapi/models"
)

// SearchUsersByUsernameReader is a Reader for the SearchUsersByUsername structure.
type SearchUsersByUsernameReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SearchUsersByUsernameReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewSearchUsersByUsernameOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewSearchUsersByUsernameDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewSearchUsersByUsernameOK creates a SearchUsersByUsernameOK with default headers values
func NewSearchUsersByUsernameOK() *SearchUsersByUsernameOK {
	return &SearchUsersByUsernameOK{}
}

/*SearchUsersByUsernameOK handles this case with default header values.

OK
*/
type SearchUsersByUsernameOK struct {
	Payload *SearchUsersByUsernameOKBody
}

func (o *SearchUsersByUsernameOK) Error() string {
	return fmt.Sprintf("[POST /search-users-by-username][%d] searchUsersByUsernameOK  %+v", 200, o.Payload)
}

func (o *SearchUsersByUsernameOK) GetPayload() *SearchUsersByUsernameOKBody {
	return o.Payload
}

func (o *SearchUsersByUsernameOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(SearchUsersByUsernameOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewSearchUsersByUsernameDefault creates a SearchUsersByUsernameDefault with default headers values
func NewSearchUsersByUsernameDefault(code int) *SearchUsersByUsernameDefault {
	return &SearchUsersByUsernameDefault{
		_statusCode: code,
	}
}

/*SearchUsersByUsernameDefault handles this case with default header values.

Generic error response.
*/
type SearchUsersByUsernameDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the search users by username default response
func (o *SearchUsersByUsernameDefault) Code() int {
	return o._statusCode
}

func (o *SearchUsersByUsernameDefault) Error() string {
	return fmt.Sprintf("[POST /search-users-by-username][%d] searchUsersByUsername default  %+v", o._statusCode, o.Payload)
}

func (o *SearchUsersByUsernameDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *SearchUsersByUsernameDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*SearchUsersByUsernameBody search users by username body
swagger:model SearchUsersByUsernameBody
*/
type SearchUsersByUsernameBody struct {

	// username
	// Required: true
	Username models.Username `json:"username"`
}

// Validate validates this search users by username body
func (o *SearchUsersByUsernameBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateUsername(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SearchUsersByUsernameBody) validateUsername(formats strfmt.Registry) error {

	if err := o.Username.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "username")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SearchUsersByUsernameBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SearchUsersByUsernameBody) UnmarshalBinary(b []byte) error {
	var res SearchUsersByUsernameBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*SearchUsersByUsernameOKBody search users by username o k body
swagger:model SearchUsersByUsernameOKBody
*/
type SearchUsersByUsernameOKBody struct {

	// profiles
	// Required: true
	// Max Items: 10
	// Min Items: 0
	Profiles []*models.Profile `json:"profiles"`
}

// Validate validates this search users by username o k body
func (o *SearchUsersByUsernameOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateProfiles(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SearchUsersByUsernameOKBody) validateProfiles(formats strfmt.Registry) error {

	if err := validate.Required("searchUsersByUsernameOK"+"."+"profiles", "body", o.Profiles); err != nil {
		return err
	}

	iProfilesSize := int64(len(o.Profiles))

	if err := validate.MinItems("searchUsersByUsernameOK"+"."+"profiles", "body", iProfilesSize, 0); err != nil {
		return err
	}

	if err := validate.MaxItems("searchUsersByUsernameOK"+"."+"profiles", "body", iProfilesSize, 10); err != nil {
		return err
	}

	for i := 0; i < len(o.Profiles); i++ {
		if swag.IsZero(o.Profiles[i]) { // not required
			continue
		}

		if o.Profiles[i] != nil {
			if err := o.Profiles[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("searchUsersByUsernameOK" + "." + "profiles" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *SearchUsersByUsernameOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SearchUsersByUsernameOKBody) UnmarshalBinary(b []byte) error {
	var res SearchUsersByUsernameOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}