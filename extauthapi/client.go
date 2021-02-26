// Package extauthapi implementaton for debugging and testing.
package extauthapi

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"
)

// Client for extauthapi.
type Client struct {
}

// NewClient creates and return new client for extauthapi.
func NewClient(endpoint string, tlsConfig *tls.Config, autoRetryCSRF bool) (*Client, error) {
	return &Client{}, nil
}

// Profile describes user profile returned by /get-user-profile.
type Profile struct {
	ID               ID
	Authn            bool
	IsolatedEntityID ID
}

func newProfile(userID, isolatedEntityID ID) *Profile {
	return &Profile{
		ID:               userID,
		Authn:            true,
		IsolatedEntityID: isolatedEntityID,
	}
}

// GetUserProfile gets a cookie with the userID and isolatedEntityID separated by a dot(.) and returns a profile with the values from the cookie.
// Authn is always true.
func (c *Client) GetUserProfile(ctx context.Context, rawCookies string) (*Profile, error) {
	cookie := parseCookieRaw(rawCookies)
	idStrs := strings.SplitN(cookie, ".", 2)

	var ids []ID
	for i := range idStrs {
		id, err := ParseID(idStrs[i])
		if err != nil {
			id = NewID()
		}
		ids = append(ids, id)
	}
	if len(ids) == 1 {
		ids = append(ids, NewID())
	}

	return newProfile(ids[0], ids[1]), nil
}

func parseCookieRaw(rawCookies string) string {
	header := http.Header{}
	header.Add("Cookie", rawCookies)
	request := http.Request{Header: header}

	cookieKey, err := request.Cookie(SessionCookieName)
	if err != nil {
		return ""
	}

	return cookieKey.Value
}
