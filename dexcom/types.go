package dexcom

import "time"

// DeviceResponse holds the response from the GET /devices endpoint
type DeviceResponse struct {
	Devices []Device `json:"devices"`
}

// Device holds device information and alert settings for that device
type Device struct {
	Model          string         `json:"model"`
	LastUploadDate string         `json:"lastUploadDate"`
	AlertSettings  []AlertSetting `json:"alertSettings"`
}

// AlertSetting describes the settings for a particular alert
type AlertSetting struct {
	AlertName   string  `json:"alertName"`
	Value       float64 `json:"value"`
	Unit        string  `json:"unit"`
	Snooze      int64   `json:"snooze"`
	Delay       int64   `json:"delay"`
	Enabled     bool    `json:"enabled"`
	SystemTime  string  `json:"systemTime"`
	DisplayTime string  `json:"displayTime"`
}

// EGVResponse holds the response to GET /egvs
type EGVResponse struct {
	Unit     string `json:"unit"`
	RateUnit string `json:"rateUnit"`
	EGVs     []EGV  `json:"egvs"`
}

// EGV estimated glucose value
type EGV struct {
	SystemTime  string   `json:"systemTime"`
	DisplayTime string   `json:"displayTime"`
	Value       float64  `json:"value"`
	Status      *string  `json:"status"`
	Trend       *string  `json:"trend"`
	TrendRate   *float64 `json:"trendRate"`
}

// EventResponse holds the response to GET /events
type EventResponse struct {
	Events []Event `json:"events"`
}

// Event is a user's event record
type Event struct {
	SystemTime   string  `json:"systemTime"`
	DisplayTime  string  `json:"displayTime"`
	EventType    string  `json:"eventType"`
	EventSubType string  `json:"eventSubType"`
	Value        float64 `json:"value"`
	Unit         string  `json:"unit"`
}

// StatRequest is used to fetch statistics
type StatRequest struct {
	Name      string    `json:"name"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	EGVRange  MinMax    `json:"egvrange"`
}

// MinMax holds a min and max value
type MinMax struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// Statistics holds the response to POST /statistics
type Statistics struct {
	HypoglycemiaRisk      string  `json:"hypoglycemiaRisk"`
	Min                   float64 `json:"min"`
	Max                   float64 `json:"max"`
	Mean                  float64 `json:"mean"`
	Median                float64 `json:"median"`
	Variance              float64 `json:"variance"`
	StdDev                float64 `json:"stdDev"`
	Sum                   float64 `json:"sum"`
	Q1                    float64 `json:"q1"`
	Q2                    float64 `json:"q2"`
	Q3                    float64 `json:"q3"`
	UtilizationPercent    float64 `json:"utilizationPercent"`
	MeanDailyCalibrations float64 `json:"meanDailyCalibrations"`
	NDays                 int64   `json:"nDays"`
	NValues               int64   `json:"nValues"`
	NBelowRange           int64   `json:"nBelowRange"`
	NWithinRange          int64   `json:"nWithinRange"`
	NAboveRange           int64   `json:"nAboveRange"`
	PercentBelowRange     float64 `json:"percentBelowRange"`
	PercentWithinRange    float64 `json:"percentWithinRange"`
	PercentAboveRange     float64 `json:"percentAboveRange"`
}

// UserToken holds the authorization info necessary to access user data
type UserToken struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresIn    int64      `json:"expires_in"`
	TokenType    string     `json:"token_type"`
	ExpireTime   *time.Time `json:"expire_time"`
}
