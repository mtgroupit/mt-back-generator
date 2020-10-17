// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// ID ID
//
// swagger:model ID
type ID strfmt.UUID4

// Validate validates this ID
func (m ID) Validate(formats strfmt.Registry) error {
	var res []error

	if err := validate.FormatOf("", "body", "uuid4", strfmt.UUID4(m).String(), formats); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
