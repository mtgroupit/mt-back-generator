package extauthapi

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"

	oapierrors "github.com/go-openapi/errors"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"
	"github.com/mtgroupit/mt-back-generator/extauthapi/client/operations"
	"github.com/pkg/errors"

	"github.com/mtgroupit/mt-back-generator/extauthapi/client"
	"github.com/mtgroupit/mt-back-generator/extauthapi/models"
)

var (
	errNoHost           = errors.New("endpoint must contain host")
	errWrongUUIDVersion = errors.New("wrong UUID version")
)

// Client for extauthapi.
type Client struct {
	*client.Authentication
}

// NewClient creates and return new client for extauthapi.
func NewClient(endpoint string, tlsConfig *tls.Config, autoRetryCSRF bool) (*Client, error) {
	p, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse endpoint: %w", err)
	}
	if p.Host == "" {
		return nil, errNoHost
	}
	basePath := client.DefaultBasePath
	if p.Path != "" {
		basePath = p.Path
	}
	schemes := client.DefaultSchemes
	if p.Scheme != "" {
		schemes = []string{p.Scheme}
	}

	transport := httptransport.New(p.Host, basePath, schemes)
	transport.Transport = newCSRFTransport(autoRetryCSRF, &http.Transport{ // Same as http.DefaultTransport plus tlsConfig.
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tlsConfig,
	})

	c := &Client{
		Authentication: client.New(transport, nil),
	}
	return c, nil
}

type csrfTransport struct {
	*http.Transport
	jar http.CookieJar

	autoRetry bool
	token     string
	muToken   sync.Mutex
}

func newCSRFTransport(autoRetry bool, tr *http.Transport) *csrfTransport {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	return &csrfTransport{
		Transport: tr,
		jar:       jar,
		autoRetry: autoRetry,
	}
}

func (tr *csrfTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	tr.muToken.Lock()
	r.Header.Set(CSRFTokenHeaderName, tr.token)
	tr.muToken.Unlock()
	for _, cookie := range tr.jar.Cookies(r.URL) {
		if _, err := r.Cookie(cookie.Name); err == http.ErrNoCookie {
			r.AddCookie(cookie)
		}
	}
	resp, err := tr.Transport.RoundTrip(r)
	if err == nil {
		token := resp.Header.Get(CSRFTokenHeaderName)
		tr.muToken.Lock()
		retry := resp.StatusCode == http.StatusForbidden && tr.token != token
		tr.token = token
		tr.muToken.Unlock()
		tr.jar.SetCookies(r.URL, resp.Cookies())
		if tr.autoRetry && retry && (r.GetBody != nil || r.ContentLength == 0) {
			if r.GetBody != nil {
				r.Body, err = r.GetBody()
				if err != nil {
					panic(err)
				}
			}
			return tr.RoundTrip(r)
		}
	}
	return resp, err
}

// Authz describes user roles/permissions.
type Authz struct {
	User    bool // true if user has validated email
	Admin   bool // true if user is an admin
	Manager bool // true if user is an manager
}

// Profile describes user profile returned by /get-user-profile.
type Profile struct {
	ID               UserID // May be empty if !Authn.
	Username         string // May be empty.
	Email            string // May be empty even if Authz.User.
	PersdataEndpoint string // Required.
	Authn            bool   // true if user authenticated by credentials provided in request.
	Authz            Authz
	IsolatedEntityID UserID // May be empty if !Authn.
}

func newProfile(id UserID, m *models.Profile) *Profile {

	return &Profile{
		ID:               id,
		Username:         string(m.Username),
		Email:            string(m.Email),
		PersdataEndpoint: string(m.PersdataEndpoint),
		Authn:            swag.BoolValue(m.Authn),
		Authz: Authz{
			User:    swag.BoolValue(m.Authz.User),
			Admin:   swag.BoolValue(m.Authz.Admin),
			Manager: swag.BoolValue(m.Authz.Manager),
		},
		IsolatedEntityID: MustParseUserID(string(m.IsolatedEntityID)),
	}
}

// GetUserProfile send /get-user-profile request to extauth service.
// Return OAPI errors.
func (c *Client) GetUserProfile(ctx context.Context, cookie string) (*Profile, error) {
	params := operations.NewGetUserProfileParamsWithContext(ctx)
	auth := httptransport.APIKeyAuth("Cookie", "header", cookie)
	res, err := c.Operations.GetUserProfile(params, auth)
	switch err := err.(type) {
	case nil:
	case *operations.GetUserProfileDefault:
		return nil, oapierrors.New(int32(err.Code()), swag.StringValue(err.Payload.Message))
	default:
		err = errors.Wrap(err, "failed to call /get-user-profile")
		return nil, oapierrors.New(http.StatusBadGateway, err.Error())
	}
	id := NoUserID
	if swag.BoolValue(res.Payload.Authn) {
		id, err = ParseUserID(string(res.Payload.ID))
		if err != nil {
			err := errors.Wrap(err, "failed to parse user ID")
			return nil, oapierrors.New(http.StatusBadGateway, err.Error())
		}
	}
	return newProfile(id, res.Payload), nil
}
