// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Profile User profile with auth details.
//
// swagger:model Profile
type Profile struct {

	// Is user authenticated by credentials provided in request.
	// Required: true
	Authn *bool `json:"authn"`

	// authz
	// Required: true
	Authz *ProfileAuthz `json:"authz"`

	// email
	// Format: email
	Email Email `json:"email,omitempty"`

	// id
	// Format: uuid4
	ID UserID `json:"id,omitempty"`

	// persdata endpoint
	// Required: true
	// Format: uri
	PersdataEndpoint Endpoint `json:"persdataEndpoint"`

	// username
	Username Username `json:"username,omitempty"`
}

// Validate validates this profile
func (m *Profile) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAuthn(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAuthz(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEmail(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePersdataEndpoint(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUsername(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Profile) validateAuthn(formats strfmt.Registry) error {

	if err := validate.Required("authn", "body", m.Authn); err != nil {
		return err
	}

	return nil
}

func (m *Profile) validateAuthz(formats strfmt.Registry) error {

	if err := validate.Required("authz", "body", m.Authz); err != nil {
		return err
	}

	if m.Authz != nil {
		if err := m.Authz.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("authz")
			}
			return err
		}
	}

	return nil
}

func (m *Profile) validateEmail(formats strfmt.Registry) error {

	if swag.IsZero(m.Email) { // not required
		return nil
	}

	if err := m.Email.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("email")
		}
		return err
	}

	return nil
}

func (m *Profile) validateID(formats strfmt.Registry) error {

	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := m.ID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("id")
		}
		return err
	}

	return nil
}

func (m *Profile) validatePersdataEndpoint(formats strfmt.Registry) error {

	if err := m.PersdataEndpoint.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("persdataEndpoint")
		}
		return err
	}

	return nil
}

func (m *Profile) validateUsername(formats strfmt.Registry) error {

	if swag.IsZero(m.Username) { // not required
		return nil
	}

	if err := m.Username.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("username")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Profile) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Profile) UnmarshalBinary(b []byte) error {
	var res Profile
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ProfileAuthz User roles/permissions needed for authorization.
//
// swagger:model ProfileAuthz
type ProfileAuthz struct {

	// Is user an admin.
	// Required: true
	Admin *bool `json:"admin"`

	// Is user an manager.
	// Required: true
	Manager *bool `json:"manager"`

	// Is user has validated email.
	// Required: true
	User *bool `json:"user"`
}

// Validate validates this profile authz
func (m *ProfileAuthz) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAdmin(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateManager(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUser(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProfileAuthz) validateAdmin(formats strfmt.Registry) error {

	if err := validate.Required("authz"+"."+"admin", "body", m.Admin); err != nil {
		return err
	}

	return nil
}

func (m *ProfileAuthz) validateManager(formats strfmt.Registry) error {

	if err := validate.Required("authz"+"."+"manager", "body", m.Manager); err != nil {
		return err
	}

	return nil
}

func (m *ProfileAuthz) validateUser(formats strfmt.Registry) error {

	if err := validate.Required("authz"+"."+"user", "body", m.User); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProfileAuthz) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProfileAuthz) UnmarshalBinary(b []byte) error {
	var res ProfileAuthz
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}