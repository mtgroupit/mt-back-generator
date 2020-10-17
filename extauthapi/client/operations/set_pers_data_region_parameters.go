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

// NewSetPersDataRegionParams creates a new SetPersDataRegionParams object
// with the default values initialized.
func NewSetPersDataRegionParams() *SetPersDataRegionParams {
	var ()
	return &SetPersDataRegionParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewSetPersDataRegionParamsWithTimeout creates a new SetPersDataRegionParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewSetPersDataRegionParamsWithTimeout(timeout time.Duration) *SetPersDataRegionParams {
	var ()
	return &SetPersDataRegionParams{

		timeout: timeout,
	}
}

// NewSetPersDataRegionParamsWithContext creates a new SetPersDataRegionParams object
// with the default values initialized, and the ability to set a context for a request
func NewSetPersDataRegionParamsWithContext(ctx context.Context) *SetPersDataRegionParams {
	var ()
	return &SetPersDataRegionParams{

		Context: ctx,
	}
}

// NewSetPersDataRegionParamsWithHTTPClient creates a new SetPersDataRegionParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewSetPersDataRegionParamsWithHTTPClient(client *http.Client) *SetPersDataRegionParams {
	var ()
	return &SetPersDataRegionParams{
		HTTPClient: client,
	}
}

/*SetPersDataRegionParams contains all the parameters to send to the API endpoint
for the set pers data region operation typically these are written to a http.Request
*/
type SetPersDataRegionParams struct {

	/*Args*/
	Args SetPersDataRegionBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the set pers data region params
func (o *SetPersDataRegionParams) WithTimeout(timeout time.Duration) *SetPersDataRegionParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the set pers data region params
func (o *SetPersDataRegionParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the set pers data region params
func (o *SetPersDataRegionParams) WithContext(ctx context.Context) *SetPersDataRegionParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the set pers data region params
func (o *SetPersDataRegionParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the set pers data region params
func (o *SetPersDataRegionParams) WithHTTPClient(client *http.Client) *SetPersDataRegionParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the set pers data region params
func (o *SetPersDataRegionParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithArgs adds the args to the set pers data region params
func (o *SetPersDataRegionParams) WithArgs(args SetPersDataRegionBody) *SetPersDataRegionParams {
	o.SetArgs(args)
	return o
}

// SetArgs adds the args to the set pers data region params
func (o *SetPersDataRegionParams) SetArgs(args SetPersDataRegionBody) {
	o.Args = args
}

// WriteToRequest writes these params to a swagger request
func (o *SetPersDataRegionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
