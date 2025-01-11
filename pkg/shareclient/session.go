package shareclient

import (
	"regexp"
	"strconv"
	"time"

	"github.com/havardelnan/dexcomsharego/pkg/glucose"
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
	WT    string `json:"WT"`
	ST    string `json:"ST"`
	DT    string `json:"DT"`
	Value int    `json:"Value"`
	Trend string `json:"Trend"`
}

func GlucoseValueFromAPI(apiValue int) glucose.GlucoseValue {
	return glucose.GlucoseValue(glucose.MgDlTommolL(apiValue))
}

func timeFromApiString(apiTime string) time.Time {
	regexp := regexp.MustCompile(`Date\((\d+)([+-]\d+\))`)
	timeparts := regexp.FindStringSubmatch(apiTime)
	unixtime, _ := strconv.Atoi(timeparts[1])
	return time.Unix(int64(unixtime)/1000, 0)
}

func (r ApiGlucoseReading) NewGlucoseReading() glucose.GlucoseReading {

	return glucose.GlucoseReading{
		Time:  timeFromApiString(r.DT),
		Value: GlucoseValueFromAPI(r.Value),
		Trend: r.Trend,
	}
}

func (r ApiGlucoseReadings) NewGlucoseReadings() glucose.GlucoseReadings {
	var readings glucose.GlucoseReadings
	for _, reading := range r {
		readings = append(readings, reading.NewGlucoseReading())
	}
	return readings
}

func (s *Sharesession) GetGlucoseReading() glucose.GlucoseReading {
	readings := s.GetGlucoseReadings(1440, 1)
	return readings[0]
}

func (s *Sharesession) GetGlucoseReadings(minutes int, maxcount int) glucose.GlucoseReadings {
	var apireadings ApiGlucoseReadings

	s.client.PostJSON("Publisher/ReadPublisherLatestGlucoseValues", map[string]string{
		"applicationId": s.AuthConfig.ApplicationId,
		"sessionId":     s.SessionId,
		"minutes":       minutesToApiString(minutes),
		"maxCount":      maxcountToApiString(maxcount),
	}, &apireadings)

	readings := apireadings.NewGlucoseReadings()

	return readings
}

func maxcountToApiString(maxcount int) string {
	if maxcount > 288 {
		return "288"
	}
	if maxcount < 1 {
		return "1"
	}
	return strconv.Itoa(maxcount)
}

func minutesToApiString(minutes int) string {
	if minutes > 1440 {
		return "1440"
	}
	if minutes < 1 {
		return "1"
	}
	return strconv.Itoa(minutes)
}
