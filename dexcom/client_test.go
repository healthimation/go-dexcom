package dexcom

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/healthimation/go-client/client"
)

func testClient(handler http.HandlerFunc, timeout time.Duration) (Client, *httptest.Server) {
	ts := httptest.NewServer(handler)
	finder := func(serviceName string, useTLS bool) (url.URL, error) {
		ret, err := url.Parse(ts.URL)
		if err != nil || ret == nil {
			return url.URL{}, err
		}
		return *ret, err
	}
	c := &dexcomClient{
		c:            client.NewBaseClient(finder, "dexcom", true, timeout),
		clientID:     "123",
		clientSecret: "abc",
	}
	return c, ts
}

func makeStrPtr(v string) *string {
	return &v
}
func makeFloat64Ptr(v float64) *float64 {
	return &v
}

func TestUnit_GetUser(t *testing.T) {

	type testcase struct {
		name             string
		handler          http.HandlerFunc
		timeout          time.Duration
		ctx              context.Context
		authCode         string
		redirectURI      string
		expectedErrCode  string
		expectedResponse *UserToken
	}

	testcases := []testcase{
		{
			name: "base path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"access_token":"access", "expires_in":600, "token_type":"Bearer", "refresh_token":"refresh"}`)
			}),
			timeout:          5 * time.Second,
			ctx:              context.Background(),
			authCode:         "123",
			redirectURI:      "abc",
			expectedResponse: &UserToken{AccessToken: "access", RefreshToken: "refresh", ExpiresIn: 600, TokenType: "Bearer"},
		},
		{
			name: "exceptional path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `invalid_request`)
			}),
			timeout:         5 * time.Second,
			ctx:             context.Background(),
			authCode:        "123",
			redirectURI:     "abc",
			expectedErrCode: ErrorAPI,
		},
		{
			name: "exceptional path - timeout",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Millisecond)
				fmt.Fprint(w, `{"access_token": "access","expires_in": 600,"token_type": "Bearer","refresh_token": "refresh"}`)
			}),
			timeout:         1 * time.Millisecond,
			ctx:             context.Background(),
			authCode:        "123",
			redirectURI:     "abc",
			expectedErrCode: client.ErrorRequestError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c, ts := testClient(tc.handler, tc.timeout)
			defer ts.Close()
			ret, err := c.GetUser(tc.ctx, tc.authCode, tc.redirectURI)
			if tc.expectedErrCode != "" || err != nil {
				if tc.expectedErrCode == "" {
					t.Fatalf("Unexpected error occurred (%#v)", err)
				}
				if err == nil {
					t.Fatalf("Expected error did not occur")
				}
				if err.Code() != tc.expectedErrCode {
					t.Fatalf("Actual error (%#v) did not match expected (%#v)", err.Code(), tc.expectedErrCode)
				}
				if tc.expectedResponse != nil && ret != nil {
					tc.expectedResponse.ExpireTime = ret.ExpireTime
				}
				if !reflect.DeepEqual(tc.expectedResponse, ret) {
					t.Fatalf("Actual response (%#v) did not match expected (%#v)", ret, tc.expectedResponse)
				}
			}
		})
	}
}

func TestUnit_RefreshUser(t *testing.T) {

	type testcase struct {
		name             string
		handler          http.HandlerFunc
		timeout          time.Duration
		ctx              context.Context
		refreshToken     string
		redirectURI      string
		expectedErrCode  string
		expectedResponse *UserToken
	}

	testcases := []testcase{
		{
			name: "base path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"access_token":"access", "expires_in":600, "token_type":"Bearer", "refresh_token":"refresh"}`)
			}),
			timeout:          5 * time.Second,
			ctx:              context.Background(),
			refreshToken:     "123",
			redirectURI:      "abc",
			expectedResponse: &UserToken{AccessToken: "access", RefreshToken: "refresh", ExpiresIn: 600, TokenType: "Bearer"},
		},
		{
			name: "exceptional path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `invalid_request`)
			}),
			timeout:         5 * time.Second,
			ctx:             context.Background(),
			refreshToken:    "123",
			redirectURI:     "abc",
			expectedErrCode: ErrorAPI,
		},
		{
			name: "exceptional path - timeout",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Millisecond)
				fmt.Fprint(w, `{"access_token": "access","expires_in": 600,"token_type": "Bearer","refresh_token": "refresh"}`)
			}),
			timeout:         1 * time.Millisecond,
			ctx:             context.Background(),
			refreshToken:    "123",
			redirectURI:     "abc",
			expectedErrCode: client.ErrorRequestError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c, ts := testClient(tc.handler, tc.timeout)
			defer ts.Close()
			ret, err := c.RefreshUser(tc.ctx, tc.refreshToken, tc.redirectURI)
			if tc.expectedErrCode != "" || err != nil {
				if tc.expectedErrCode == "" {
					t.Fatalf("Unexpected error occurred (%#v)", err)
				}
				if err == nil {
					t.Fatalf("Expected error did not occur")
				}
				if err.Code() != tc.expectedErrCode {
					t.Fatalf("Actual error (%#v) did not match expected (%#v)", err.Code(), tc.expectedErrCode)
				}
				if tc.expectedResponse != nil && ret != nil {
					tc.expectedResponse.ExpireTime = ret.ExpireTime
				}
				if !reflect.DeepEqual(tc.expectedResponse, ret) {
					t.Fatalf("Actual response (%#v) did not match expected (%#v)", ret, tc.expectedResponse)
				}
			}
		})
	}
}

func TestUnit_GetDevices(t *testing.T) {

	type testcase struct {
		name             string
		handler          http.HandlerFunc
		timeout          time.Duration
		ctx              context.Context
		accessToken      string
		startDate        time.Time
		endDate          time.Time
		expectedErrCode  string
		expectedResponse *DeviceResponse
	}

	testcases := []testcase{
		{
			name: "base path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"devices": [{"model": "G5 Mobile App","lastUploadDate": "2016-08-15T00:00:00","alertSettings": [{"alertName": "high","value": 200,"unit": "mg/dL","snooze": 30,"delay": 0,"enabled": true,"systemTime": "2016-08-15T00:00:00","displayTime": "2016-08-15T00:00:00"}]}]}`)
			}),
			timeout:          5 * time.Second,
			ctx:              context.Background(),
			accessToken:      "123",
			startDate:        time.Now(),
			endDate:          time.Now(),
			expectedResponse: &DeviceResponse{Devices: []Device{Device{Model: "G5 Mobile App", LastUploadDate: "2016-08-15T00:00:00", AlertSettings: []AlertSetting{AlertSetting{AlertName: "high", Value: 200, Unit: "mg/dl", Snooze: 30, Delay: 0, Enabled: true, SystemTime: "2016-08-15T00:00:00", DisplayTime: "2016-08-15T00:00:00"}}}}},
		},
		{
			name: "exceptional path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `invalid_request`)
			}),
			timeout:         5 * time.Second,
			ctx:             context.Background(),
			accessToken:     "123",
			startDate:       time.Now(),
			endDate:         time.Now(),
			expectedErrCode: ErrorAPI,
		},
		{
			name: "exceptional path - timeout",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Millisecond)
				fmt.Fprint(w, `foo`)
			}),
			timeout:         1 * time.Millisecond,
			ctx:             context.Background(),
			accessToken:     "123",
			startDate:       time.Now(),
			endDate:         time.Now(),
			expectedErrCode: client.ErrorRequestError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c, ts := testClient(tc.handler, tc.timeout)
			defer ts.Close()
			ret, err := c.GetDevices(tc.ctx, tc.accessToken, tc.startDate, tc.endDate)
			if tc.expectedErrCode != "" || err != nil {
				if tc.expectedErrCode == "" {
					t.Fatalf("Unexpected error occurred (%#v)", err)
				}
				if err == nil {
					t.Fatalf("Expected error did not occur")
				}
				if err.Code() != tc.expectedErrCode {
					t.Fatalf("Actual error (%#v) did not match expected (%#v)", err.Code(), tc.expectedErrCode)
				}
				if !reflect.DeepEqual(tc.expectedResponse, ret) {
					t.Fatalf("Actual response (%#v) did not match expected (%#v)", ret, tc.expectedResponse)
				}
			}
		})
	}
}

func TestUnit_GetEGVs(t *testing.T) {

	type testcase struct {
		name             string
		handler          http.HandlerFunc
		timeout          time.Duration
		ctx              context.Context
		accessToken      string
		expectedErrCode  string
		expectedResponse *EGVResponse
	}

	testcases := []testcase{
		{
			name: "base path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"unit": "mg/dL","rateUnit": "mg/dL/min","egvs": [{"systemTime": "2017-06-16T15:40:00","displayTime": "2017-06-16T07:40:00","value": 119,"status": null,"trend": "fortyFiveDown","trendRate": -1.3}]}`)
			}),
			timeout:          5 * time.Second,
			ctx:              context.Background(),
			accessToken:      "123",
			expectedResponse: &EGVResponse{Unit: "mg/dL", RateUnit: "mg/dL/min", EGVs: []EGV{EGV{SystemTime: "2017-06-16T15:40:00", DisplayTime: "2017-06-16T07:40:00", Value: 119, Trend: makeStrPtr("fortyFiveDown"), TrendRate: makeFloat64Ptr(-1.3)}}},
		},
		{
			name: "exceptional path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `invalid_request`)
			}),
			timeout:         5 * time.Second,
			ctx:             context.Background(),
			accessToken:     "123",
			expectedErrCode: ErrorAPI,
		},
		{
			name: "exceptional path - timeout",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Millisecond)
				fmt.Fprint(w, `foo`)
			}),
			timeout:         1 * time.Millisecond,
			ctx:             context.Background(),
			accessToken:     "123",
			expectedErrCode: client.ErrorRequestError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c, ts := testClient(tc.handler, tc.timeout)
			defer ts.Close()
			ret, err := c.GetEGVs(tc.ctx, tc.accessToken)
			if tc.expectedErrCode != "" || err != nil {
				if tc.expectedErrCode == "" {
					t.Fatalf("Unexpected error occurred (%#v)", err)
				}
				if err == nil {
					t.Fatalf("Expected error did not occur")
				}
				if err.Code() != tc.expectedErrCode {
					t.Fatalf("Actual error (%#v) did not match expected (%#v)", err.Code(), tc.expectedErrCode)
				}
				if !reflect.DeepEqual(tc.expectedResponse, ret) {
					t.Fatalf("Actual response (%#v) did not match expected (%#v)", ret, tc.expectedResponse)
				}
			}
		})
	}
}

func TestUnit_GetEvents(t *testing.T) {

	type testcase struct {
		name             string
		handler          http.HandlerFunc
		timeout          time.Duration
		ctx              context.Context
		accessToken      string
		startDate        time.Time
		endDate          time.Time
		expectedErrCode  string
		expectedResponse *EventResponse
	}

	testcases := []testcase{
		{
			name: "base path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"events": [{"systemTime": "2017-06-16T19:45:00","displayTime": "2017-06-16T11:45:00","eventType": "exercise","eventSubType": "medium","value": 42,"unit": "minutes"}]}`)
			}),
			timeout:          5 * time.Second,
			ctx:              context.Background(),
			accessToken:      "123",
			startDate:        time.Now(),
			endDate:          time.Now(),
			expectedResponse: &EventResponse{Events: []Event{Event{SystemTime: "2017-06-16T19:45:00", DisplayTime: "2017-06-16T11:45:00", EventType: "exercise", EventSubType: "medium", Value: 42, Unit: "minutes"}}},
		},
		{
			name: "exceptional path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `invalid_request`)
			}),
			timeout:         5 * time.Second,
			ctx:             context.Background(),
			accessToken:     "123",
			startDate:       time.Now(),
			endDate:         time.Now(),
			expectedErrCode: ErrorAPI,
		},
		{
			name: "exceptional path - timeout",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Millisecond)
				fmt.Fprint(w, `foo`)
			}),
			timeout:         1 * time.Millisecond,
			ctx:             context.Background(),
			accessToken:     "123",
			startDate:       time.Now(),
			endDate:         time.Now(),
			expectedErrCode: client.ErrorRequestError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c, ts := testClient(tc.handler, tc.timeout)
			defer ts.Close()
			ret, err := c.GetEvents(tc.ctx, tc.accessToken, tc.startDate, tc.endDate)
			if tc.expectedErrCode != "" || err != nil {
				if tc.expectedErrCode == "" {
					t.Fatalf("Unexpected error occurred (%#v)", err)
				}
				if err == nil {
					t.Fatalf("Expected error did not occur")
				}
				if err.Code() != tc.expectedErrCode {
					t.Fatalf("Actual error (%#v) did not match expected (%#v)", err.Code(), tc.expectedErrCode)
				}
				if !reflect.DeepEqual(tc.expectedResponse, ret) {
					t.Fatalf("Actual response (%#v) did not match expected (%#v)", ret, tc.expectedResponse)
				}
			}
		})
	}
}

func TestUnit_GetStatistics(t *testing.T) {

	type testcase struct {
		name             string
		handler          http.HandlerFunc
		timeout          time.Duration
		ctx              context.Context
		accessToken      string
		startDate        time.Time
		endDate          time.Time
		stats            map[string][]StatRequest
		expectedErrCode  string
		expectedResponse *Statistics
	}

	testcases := []testcase{
		{
			name: "base path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"hypoglycemiaRisk": "minimal","min": 39,"max": 287,"mean": 131.20452051788453,"median": 121,"variance": 1836.660387728187,"stdDev": 42.85627594329898,"sum": 597899,"q1": 100,"q2": 121,"q3": 155,"utilizationPercent": 98.89322916666666,"meanDailyCalibrations": 2,"nDays": 16,"nValues": 4557,"nBelowRange": 185,"nWithinRange": 3605,"nAboveRange": 767,"percentBelowRange": 4.0596883914856265,"percentWithinRange": 79.10906298003071,"percentAboveRange": 16.831248628483653}`)
			}),
			timeout:     5 * time.Second,
			ctx:         context.Background(),
			accessToken: "123",
			startDate:   time.Now(),
			endDate:     time.Now(),
			stats:       map[string][]StatRequest{"foo": []StatRequest{StatRequest{Name: "foo", StartTime: time.Now(), EndTime: time.Now(), EGVRange: MinMax{Min: 0, Max: 100}}}},
			expectedResponse: &Statistics{HypoglycemiaRisk: "minimal", Min: 39, Max: 287, Mean: 131.20452051788453, Median: 121, Variance: 1836.660387728187, StdDev: 42.85627594329898, Sum: 597899, Q1: 100, Q2: 121, Q3: 155,
				UtilizationPercent: 98.89322916666666, MeanDailyCalibrations: 2, NDays: 16, NValues: 4557, NBelowRange: 185, NAboveRange: 767, NWithinRange: 3605, PercentBelowRange: 4.0596883914856265, PercentWithinRange: 79.10906298003071, PercentAboveRange: 16.831248628483653},
		},
		{
			name: "exceptional path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `invalid_request`)
			}),
			timeout:         5 * time.Second,
			ctx:             context.Background(),
			accessToken:     "123",
			startDate:       time.Now(),
			endDate:         time.Now(),
			expectedErrCode: ErrorAPI,
		},
		{
			name: "exceptional path - timeout",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Millisecond)
				fmt.Fprint(w, `foo`)
			}),
			timeout:         1 * time.Millisecond,
			ctx:             context.Background(),
			accessToken:     "123",
			startDate:       time.Now(),
			endDate:         time.Now(),
			expectedErrCode: client.ErrorRequestError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c, ts := testClient(tc.handler, tc.timeout)
			defer ts.Close()
			ret, err := c.GetStatistics(tc.ctx, tc.accessToken, tc.startDate, tc.endDate, tc.stats)
			if tc.expectedErrCode != "" || err != nil {
				if tc.expectedErrCode == "" {
					t.Fatalf("Unexpected error occurred (%#v)", err)
				}
				if err == nil {
					t.Fatalf("Expected error did not occur")
				}
				if err.Code() != tc.expectedErrCode {
					t.Fatalf("Actual error (%#v) did not match expected (%#v)", err.Code(), tc.expectedErrCode)
				}
				if !reflect.DeepEqual(tc.expectedResponse, ret) {
					t.Fatalf("Actual response (%#v) did not match expected (%#v)", ret, tc.expectedResponse)
				}
			}
		})
	}
}
