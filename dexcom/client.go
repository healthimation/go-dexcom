package dexcom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/healthimation/go-client/client"
	"github.com/healthimation/go-glitch/glitch"
)

//Error codes
const (
	ErrorAPI          = "ERROR_API"
	ErrorJSON         = "ERROR_JSON"
	ErrorMissingParam = "ERROR_MISSING_PARAM"

	// grant types
	grantTypeAuthorizationCode = "authorization_code"
	grantTypeRefreshToken      = "refresh_token"

	// params
	paramClientID          = "client_id"
	paramClientSecret      = "client_secret"
	paramAuthorizationCode = "code"
	paramGrantType         = "grant_type"
	paramRedirectURI       = "redirect_uri"
	paramRefreshToken      = "refresh_token"

	paramStartDate = "startDate"
	paramEndDate   = "endDate"

	timeformat = "2006-01-02T15:04:05"
)

// Client can make requests to the pushy api
type Client interface {
	GetUser(ctx context.Context, authorizationCode, redirectURI string) (*UserToken, glitch.DataError)
	RefreshUser(ctx context.Context, refreshToken, redirectURI string) (*UserToken, glitch.DataError)
	GetDevices(ctx context.Context, accessToken string, startDate, endDate time.Time) (*DeviceResponse, glitch.DataError)
	GetEGVs(ctx context.Context, accessToken string, startDate, endDate time.Time) (*EGVResponse, glitch.DataError)
	GetEvents(ctx context.Context, accessToken string, startDate, endDate time.Time) (*EventResponse, glitch.DataError)
	GetStatistics(ctx context.Context, accessToken string, startDate, endDate time.Time, stats map[string][]StatRequest) (*Statistics, glitch.DataError)
}

type dexcomClient struct {
	c            client.BaseClient
	clientID     string
	clientSecret string
}

// NewClient returns a new pushy client
func NewClient(clientID string, clientSecret string, timeout time.Duration) Client {
	return &dexcomClient{
		c:            client.NewBaseClient(findDexcom, "dexcom", true, timeout),
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// NewSandboxClient gets a client that talks to the dexcom sandbox
func NewSandboxClient(clientID string, clientSecret string, timeout time.Duration) Client {
	return &dexcomClient{
		c:            client.NewBaseClient(findDexcomSandbox, "dexcom", true, timeout),
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (d *dexcomClient) getUser(ctx context.Context, authorizationCode, refreshToken, redirectURI string) (*UserToken, glitch.DataError) {
	slug := "v1/oauth2/token"
	h := http.Header{}
	h.Set("Content-type", "application/x-www-form-urlencoded")
	h.Set("cache-control", "no-cache")

	values := url.Values{}
	values.Set(paramClientID, d.clientID)
	values.Set(paramClientSecret, d.clientSecret)
	values.Set(paramRedirectURI, redirectURI)
	if len(authorizationCode) > 0 {
		values.Set(paramAuthorizationCode, authorizationCode)
		values.Set(paramGrantType, grantTypeAuthorizationCode)
	} else if len(refreshToken) > 0 {
		values.Set(paramRefreshToken, refreshToken)
		values.Set(paramGrantType, grantTypeRefreshToken)
	} else {
		return nil, glitch.NewDataError(nil, ErrorMissingParam, "authorization_code or refresh_token is missing")
	}

	statusCode, ret, err := d.c.MakeRequest(ctx, http.MethodPost, slug, nil, h, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	result := new(UserToken)
	if statusCode >= 200 && statusCode < 300 {
		err := json.Unmarshal(ret, result)
		if err != nil {
			return nil, glitch.NewDataError(err, ErrorJSON, fmt.Sprintf("Could not unmarshal response with code %d | %s", statusCode, err.Error()))
		}
		t := time.Now().Add(time.Duration(result.ExpiresIn-5) * time.Second) //5 second buffer
		result.ExpireTime = &t
		return result, nil
	}
	return nil, glitch.NewDataError(fmt.Errorf("Error from API: %d - %s", statusCode, ret), ErrorAPI, fmt.Sprintf("Status code was not in the 2xx range: %d", statusCode))
}

func (d *dexcomClient) GetUser(ctx context.Context, authorizationCode, redirectURI string) (*UserToken, glitch.DataError) {
	return d.getUser(ctx, authorizationCode, "", redirectURI)
}

func (d *dexcomClient) RefreshUser(ctx context.Context, refreshToken, redirectURI string) (*UserToken, glitch.DataError) {
	return d.getUser(ctx, "", refreshToken, redirectURI)
}

func (d *dexcomClient) GetDevices(ctx context.Context, accessToken string, startDate, endDate time.Time) (*DeviceResponse, glitch.DataError) {
	slug := "/v1/users/self/devices"
	h := http.Header{}
	h.Add("authorization", fmt.Sprintf("Bearer %s", accessToken))

	q := url.Values{}
	q.Set(paramStartDate, startDate.UTC().Format(timeformat))
	q.Set(paramEndDate, endDate.UTC().Format(timeformat))

	statusCode, ret, err := d.c.MakeRequest(ctx, http.MethodGet, slug, q, h, nil)
	if err != nil {
		return nil, err
	}

	result := new(DeviceResponse)
	if statusCode >= 200 && statusCode < 300 {
		err := json.Unmarshal(ret, result)
		if err != nil {
			return nil, glitch.NewDataError(err, ErrorJSON, fmt.Sprintf("Could not unmarshal response with code %d | %s", statusCode, err.Error()))
		}
		return result, nil
	}
	return nil, glitch.NewDataError(fmt.Errorf("Error from API: %d - %s", statusCode, ret), ErrorAPI, fmt.Sprintf("Status code was not in the 2xx range: %d", statusCode))
}

func (d *dexcomClient) GetEGVs(ctx context.Context, accessToken string, startDate, endDate time.Time) (*EGVResponse, glitch.DataError) {
	slug := "/v1/users/self/egvs"
	h := http.Header{}
	h.Add("authorization", fmt.Sprintf("Bearer %s", accessToken))

	q := url.Values{}
	q.Set(paramStartDate, startDate.UTC().Format(timeformat))
	q.Set(paramEndDate, endDate.UTC().Format(timeformat))

	statusCode, ret, err := d.c.MakeRequest(ctx, http.MethodGet, slug, q, h, nil)
	if err != nil {
		return nil, err
	}

	result := new(EGVResponse)
	if statusCode >= 200 && statusCode < 300 {
		err := json.Unmarshal(ret, result)
		if err != nil {
			return nil, glitch.NewDataError(err, ErrorJSON, fmt.Sprintf("Could not unmarshal response with code %d | %s", statusCode, err.Error()))
		}
		return result, nil
	}
	return nil, glitch.NewDataError(fmt.Errorf("Error from API: %d - %s", statusCode, ret), ErrorAPI, fmt.Sprintf("Status code was not in the 2xx range: %d", statusCode))
}

func (d *dexcomClient) GetEvents(ctx context.Context, accessToken string, startDate, endDate time.Time) (*EventResponse, glitch.DataError) {
	slug := "/v1/users/self/events"
	h := http.Header{}
	h.Add("authorization", fmt.Sprintf("Bearer %s", accessToken))

	q := url.Values{}
	q.Set(paramStartDate, startDate.UTC().Format(timeformat))
	q.Set(paramEndDate, endDate.UTC().Format(timeformat))

	statusCode, ret, err := d.c.MakeRequest(ctx, http.MethodGet, slug, q, h, nil)
	if err != nil {
		return nil, err
	}

	result := new(EventResponse)
	if statusCode >= 200 && statusCode < 300 {
		err := json.Unmarshal(ret, result)
		if err != nil {
			return nil, glitch.NewDataError(err, ErrorJSON, fmt.Sprintf("Could not unmarshal response with code %d | %s", statusCode, err.Error()))
		}
		return result, nil
	}
	return nil, glitch.NewDataError(fmt.Errorf("Error from API: %d - %s", statusCode, ret), ErrorAPI, fmt.Sprintf("Status code was not in the 2xx range: %d", statusCode))
}

func (d *dexcomClient) GetStatistics(ctx context.Context, accessToken string, startDate, endDate time.Time, stats map[string][]StatRequest) (*Statistics, glitch.DataError) {
	slug := "/v1/users/self/statistics"
	h := http.Header{}
	h.Add("authorization", fmt.Sprintf("Bearer %s", accessToken))

	q := url.Values{}
	q.Set(paramStartDate, startDate.UTC().Format(timeformat))
	q.Set(paramEndDate, endDate.UTC().Format(timeformat))

	body, err := client.ObjectToJSONReader(stats)
	if err != nil {
		return nil, err
	}

	statusCode, ret, err := d.c.MakeRequest(ctx, http.MethodPost, slug, q, h, body)
	if err != nil {
		return nil, err
	}

	result := new(Statistics)
	if statusCode >= 200 && statusCode < 300 {
		err := json.Unmarshal(ret, result)
		if err != nil {
			return nil, glitch.NewDataError(err, ErrorJSON, fmt.Sprintf("Could not unmarshal response with code %d | %s", statusCode, err.Error()))
		}
		return result, nil
	}
	return nil, glitch.NewDataError(fmt.Errorf("Error from API: %d - %s", statusCode, ret), ErrorAPI, fmt.Sprintf("Status code was not in the 2xx range: %d", statusCode))
}
