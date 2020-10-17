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

// SetPersDataRegionReader is a Reader for the SetPersDataRegion structure.
type SetPersDataRegionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SetPersDataRegionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewSetPersDataRegionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewSetPersDataRegionDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewSetPersDataRegionOK creates a SetPersDataRegionOK with default headers values
func NewSetPersDataRegionOK() *SetPersDataRegionOK {
	return &SetPersDataRegionOK{}
}

/*SetPersDataRegionOK handles this case with default header values.

OK
*/
type SetPersDataRegionOK struct {
	Payload *SetPersDataRegionOKBody
}

func (o *SetPersDataRegionOK) Error() string {
	return fmt.Sprintf("[POST /set-persdata-region][%d] setPersDataRegionOK  %+v", 200, o.Payload)
}

func (o *SetPersDataRegionOK) GetPayload() *SetPersDataRegionOKBody {
	return o.Payload
}

func (o *SetPersDataRegionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(SetPersDataRegionOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewSetPersDataRegionDefault creates a SetPersDataRegionDefault with default headers values
func NewSetPersDataRegionDefault(code int) *SetPersDataRegionDefault {
	return &SetPersDataRegionDefault{
		_statusCode: code,
	}
}

/*SetPersDataRegionDefault handles this case with default header values.

Generic error response.
*/
type SetPersDataRegionDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the set pers data region default response
func (o *SetPersDataRegionDefault) Code() int {
	return o._statusCode
}

func (o *SetPersDataRegionDefault) Error() string {
	return fmt.Sprintf("[POST /set-persdata-region][%d] setPersDataRegion default  %+v", o._statusCode, o.Payload)
}

func (o *SetPersDataRegionDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *SetPersDataRegionDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*SetPersDataRegionBody set pers data region body
swagger:model SetPersDataRegionBody
*/
type SetPersDataRegionBody struct {

	// country code
	// Required: true
	CountryCode models.CountryCode `json:"countryCode"`
}

// Validate validates this set pers data region body
func (o *SetPersDataRegionBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateCountryCode(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SetPersDataRegionBody) validateCountryCode(formats strfmt.Registry) error {

	if err := o.CountryCode.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("args" + "." + "countryCode")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SetPersDataRegionBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SetPersDataRegionBody) UnmarshalBinary(b []byte) error {
	var res SetPersDataRegionBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*SetPersDataRegionOKBody set pers data region o k body
swagger:model SetPersDataRegionOKBody
*/
type SetPersDataRegionOKBody struct {

	// persdata endpoint
	// Required: true
	// Format: uri
	PersdataEndpoint models.Endpoint `json:"persdataEndpoint"`
}

// Validate validates this set pers data region o k body
func (o *SetPersDataRegionOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validatePersdataEndpoint(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SetPersDataRegionOKBody) validatePersdataEndpoint(formats strfmt.Registry) error {

	if err := o.PersdataEndpoint.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("setPersDataRegionOK" + "." + "persdataEndpoint")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SetPersDataRegionOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SetPersDataRegionOKBody) UnmarshalBinary(b []byte) error {
	var res SetPersDataRegionOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
