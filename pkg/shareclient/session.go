package shareclient

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"
)

type ShareAuthConfig struct {
	ApplicationId string
	Username      string
	Password      string
}

type Sharesession struct {
	client     *ShareClient
	AuthConfig ShareAuthConfig
	AccountId  string
	SessionId  string
}

func NewSharesession(config ShareAuthConfig) (*Sharesession, error) {

	client := NewShareClient()
	s := Sharesession{
		client:     client,
		AuthConfig: config,
	}
	err := s.getAccountId()
	if err != nil {
		return nil, err
	}

	err = s.getSessionId()
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *Sharesession) getAccountId() error {
	accountid := ""
	s.client.PostJSON("General/AuthenticatePublisherAccount", map[string]string{
		"applicationId": s.AuthConfig.ApplicationId,
		"accountName":   s.AuthConfig.Username,
		"password":      s.AuthConfig.Password,
	}, &accountid)

	s.AccountId = accountid
	return nil
}

func (s *Sharesession) getSessionId() error {
	sessionid := ""
	s.client.PostJSON("General/LoginPublisherAccountById", map[string]string{
		"applicationId": s.AuthConfig.ApplicationId,
		"accountId":     s.AccountId,
		"password":      s.AuthConfig.Password,
	}, &sessionid)

	s.SessionId = sessionid
	return nil
}

type ApiGlucoseReadings []ApiGlucoseReading

type ApiGlucoseReading struct {
	// {
	// 	"WT": "Date(1736610842000)",
	// 	"ST": "Date(1736610842000)",
	// 	"DT": "Date(1736610842000+0100)",
	// 	"Value": 187,
	// 	"Trend": "Flat"
	//   }
	WT    string `json:"WT"`
	ST    string `json:"ST"`
	DT    string `json:"DT"`
	Value int    `json:"Value"`
	Trend string `json:"Trend"`
}

type GlucoseReadings []GlucoseReading

type GlucoseReading struct {
	Time  time.Time
	Value GlucoseValue
	Trend string
}

type GlucoseValue float64

// String returns the string representation of the GlucoseValue in mmol/L
func mgDlTommolL(val int) float64 {
	mmol := float64(val) * 0.0555
	return (math.Round(mmol*10) / 10)
}

func (v GlucoseValue) String() string {
	return fmt.Sprintf("%.1f mmol/L", v)
}

func GlucoseValueFromAPI(apiValue int) GlucoseValue {
	return GlucoseValue(mgDlTommolL(apiValue))
}

func timeFromApiString(apiTime string) time.Time {
	regexp := regexp.MustCompile(`Date\((\d+)([+-]\d+\))`)
	timeparts := regexp.FindStringSubmatch(apiTime)
	unixtime, _ := strconv.Atoi(timeparts[1])
	return time.Unix(int64(unixtime)/1000, 0)
}

func (r ApiGlucoseReading) NewGlucoseReading() GlucoseReading {

	return GlucoseReading{
		Time:  timeFromApiString(r.DT),
		Value: GlucoseValueFromAPI(r.Value),
		Trend: r.Trend,
	}
}

func (r ApiGlucoseReadings) NewGlucoseReadings() GlucoseReadings {
	var readings GlucoseReadings
	for _, reading := range r {
		readings = append(readings, reading.NewGlucoseReading())
	}
	return readings
}

func (s *Sharesession) GetGlucoseReading() GlucoseReading {
	var apireadings ApiGlucoseReadings
	s.client.PostJSON("Publisher/ReadPublisherLatestGlucoseValues", map[string]string{
		"applicationId": s.AuthConfig.ApplicationId,
		"sessionId":     s.SessionId,
		"minutes":       "1440",
		"maxCount":      "2",
	}, &apireadings)

	readings := apireadings.NewGlucoseReadings()

	return readings[0]
}
