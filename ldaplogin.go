package oktalib

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

// LdapLogin takes in user and password and authenticate against okta
func (o *OktaClient) LdapLogin() error {
	var (
		res      *http.Response
		body     []byte
		toreturn OktaUserAuthn
		err      error
	)
	p := oktaUser{
		Username: o.Username,
		Password: o.Password,
	}
	j := &doRequestInput{
		Dest:   o.OktaUrl,
		Method: "POST",
		Headers: map[string]string{
			"Accept":       "application/json",
			"Content-Type": "application/json",
		},
		CookieJar: o.CookieJar,
		Payload:   p,
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
	switch {
	case res.StatusCode == 401:
		return errors.Errorf("Authentication failed. Invalid username or password.")
	case res.StatusCode == 200:
	default:
		return fmt.Errorf("Status code not 200: " + string(res.Status) + " " + string(body))
	}
	o.UserAuth = &toreturn
	return nil
}
