// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewSetEmailParams creates a new SetEmailParams object
// with the default values initialized.
func NewSetEmailParams() *SetEmailParams {
	var ()
	return &SetEmailParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewSetEmailParamsWithTimeout creates a new SetEmailParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewSetEmailParamsWithTimeout(timeout time.Duration) *SetEmailParams {
	var ()
	return &SetEmailParams{

		timeout: timeout,
	}
}

// NewSetEmailParamsWithContext creates a new SetEmailParams object
// with the default values initialized, and the ability to set a context for a request
func NewSetEmailParamsWithContext(ctx context.Context) *SetEmailParams {
	var ()
	return &SetEmailParams{

		Context: ctx,
	}
}

// NewSetEmailParamsWithHTTPClient creates a new SetEmailParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewSetEmailParamsWithHTTPClient(client *http.Client) *SetEmailParams {
	var ()
	return &SetEmailParams{
		HTTPClient: client,
	}
}

/*SetEmailParams contains all the parameters to send to the API endpoint
for the set email operation typically these are written to a http.Request
*/
type SetEmailParams struct {

	/*Args*/
	Args SetEmailBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the set email params
func (o *SetEmailParams) WithTimeout(timeout time.Duration) *SetEmailParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the set email params
func (o *SetEmailParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the set email params
func (o *SetEmailParams) WithContext(ctx context.Context) *SetEmailParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the set email params
func (o *SetEmailParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the set email params
func (o *SetEmailParams) WithHTTPClient(client *http.Client) *SetEmailParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the set email params
func (o *SetEmailParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithArgs adds the args to the set email params
func (o *SetEmailParams) WithArgs(args SetEmailBody) *SetEmailParams {
	o.SetArgs(args)
	return o
}

// SetArgs adds the args to the set email params
func (o *SetEmailParams) SetArgs(args SetEmailBody) {
	o.Args = args
}

// WriteToRequest writes these params to a swagger request
func (o *SetEmailParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if err := r.SetBodyParam(o.Args); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}