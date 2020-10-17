package extauthapi_test

import (
	"testing"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/powerman/check"

	extauthapi "github.com/Lisss13/french-back-template/snouki-mobile/mtmb-extauthapi"
	"github.com/Lisss13/french-back-template/snouki-mobile/mtmb-extauthapi/client/operations"
	"github.com/Lisss13/french-back-template/snouki-mobile/mtmb-extauthapi/models"
)

func TestGetUserProfile(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	ts, tlsConfig := testNewServer()
	defer ts.Close()
	c, err := extauthapi.NewClient(ts.URL, tlsConfig, false)
	t.Nil(err)

	testCases := []struct {
		cookies string
		want    *extauthapi.Profile
	}{
		{"", &extauthapi.Profile{}},
		{"bad", &extauthapi.Profile{}},
		{sessUser, profileUser},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			profile, err := c.GetUserProfile(ctx, tc.cookies)
			t.Nil(err)
			t.DeepEqual(profile, tc.want)
		})
	}
}

func TestClient_Auth(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	ts, tlsConfig := testNewServer()
	defer ts.Close()
	c, err := extauthapi.NewClient(ts.URL, tlsConfig, false)
	t.Nil(err)

	params := op.NewSetUsernameParamsWithContext(ctx).WithArgs(op.SetUsernameBody{Username: "NewUsername"})

	var (
		authEmpty     = httptransport.APIKeyAuth("Cookie", "header", "")
		authUser1     = httptransport.APIKeyAuth("Cookie", "header", "user1")
		authEmptyCSRF = httptransport.APIKeyAuth("X-CSRFTokenBound", "header", "user1csrf")
		authUser1CSRF = runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
			err := r.SetHeaderParam("Cookie", "user1")
			if err == nil {
				err = r.SetHeaderParam("X-CSRFTokenBound", "user1csrf")
			}
			return err
		})
	)

	testCases := []struct {
		auth    runtime.ClientAuthInfoWriter
		wantmsg string
	}{
		{nil, "unauthenticated for invalid credentials"},
		{authEmpty, "unauthenticated for invalid credentials"},
		{authUser1, ""},
		{authEmptyCSRF, "unauthenticated for invalid credentials"},
		{authUser1CSRF, ""},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := c.Op.SetUsername(params, tc.auth)
			if tc.wantmsg == "" {
				t.Nil(errPayload(err))
				t.DeepEqual(res, op.NewSetUsernameNoContent())
			} else {
				t.DeepEqual(errPayload(err), &models.Error{
					Code:    swag.Int32(401),
					Message: swag.String(tc.wantmsg),
				})
				t.Nil(res)
			}
		})
	}
}
