package oktalib

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

// OktaMfa: This function serve as the mfa part of okta's authentication chain
func (o *OktaClient) OktaMfa(FactorType string, PassCode string) error {
	var (
		body     []byte
		res      *http.Response
		err      error
		factor   *int
		toreturn OktaUserAuthn
	)

	for i, e := range o.UserAuth.Embedded.Factors {
		if e.FactorType == FactorType {
			factor = addIntPtr(factor, i)
		}
	}

	if factor == nil {
		return fmt.Errorf("there's no factor avalible for " + FactorType)
	}
	// generate data
	p := oktaStateToken{
		StateToken: o.UserAuth.StateToken,
		PassCode:   PassCode,
	}
	//setting transport
	j := &doRequestInput{
		Dest:    o.OktaUrl + "/factors/" + o.UserAuth.Embedded.Factors[*factor].Id + "/verify",
		Payload: p,
		Headers: map[string]string{
			"Accept":       "application/json",
			"Content-Type": "application/json",
		},
		Method:    "POST",
		CookieJar: o.CookieJar,
	}
	res, err = doRequest(j)
	if err != nil {
		return fmt.Errorf("Unable to send payload %v", err)
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Unable to read resBody %v", err)
	}
	err = json.Unmarshal(body, &toreturn)
	if err != nil {
		return err
	}
	if FactorType == "push" {
		counter := 15
		fmt.Println("⏱ Waiting to hear from Okta...")
		for toreturn.Status != "SUCCESS" {
			if counter < 0 {
				return fmt.Errorf("Too long of a wait exiting")
			}
			res, err = doRequest(j)
			if err != nil {
				return fmt.Errorf("Unable to send payload %v", err)
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return fmt.Errorf("Unable to read respBody %v", err)
			}
			err = json.Unmarshal(body, &toreturn)
			if err != nil {
				return err
			}
			time.Sleep(2 * time.Second)
			counter--
		}
		color.Green("✅ Push received")
	}
	switch {
	case res.StatusCode == 403:
		return errors.Errorf("Invalid MFA code.")
	case res.StatusCode == 200:
	default:
		return fmt.Errorf("Status code not 200: " + res.Status + " " + string(body))
	}
	o.UserAuth.SessionToken = toreturn.SessionToken
	return nil
}
func addIntPtr(g *int, toadd int) *int {
	if g == nil {
		return &toadd
	} else {
		gg := *g + toadd
		return &gg
	}
}
