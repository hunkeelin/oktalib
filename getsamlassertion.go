package oktalib

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

type samlresp struct {
	OriData string
	RawData []byte
	SamlAss string
	Sresp   *Response
}

// This is the third step of the authentication chain where we obtain the SAML assertion.
func (o *OktaClient) GetSamlAssertion() error {
	var toreturn []byte
	notAssignedPattern := regexp.MustCompile(`(?i)not\s+assigned`)
	j := &doRequestInput{
		Dest:      o.AwsSamlUrl + "/?onetimetoken=" + o.UserAuth.SessionToken,
		Method:    "GET",
		CookieJar: o.CookieJar,
	}
	res, err := doRequest(j)
	if err != nil {
		return fmt.Errorf("Unable to send payload %v", err)
	}
	defer res.Body.Close()
	toreturn, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Unable to read resBody %v", err)
	}
	err = o.parseSAML(toreturn)
	if err != nil {
		if notAssignedPattern.Match(toreturn) {
			return errors.New("user not assigned the AWS app in Okta")
		} else if err == io.EOF {
			return errors.New("EOF encountered parsing Okta SAML response")
		} else {
			return err
		}
	}
	return nil
}

// Parsing the SAML assertion.
func (o *OktaClient) parseSAML(body []byte) error {
	var val string
	var data []byte
	var doc *html.Node
	var resp *Response
	var toreturn samlresp

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	val, _ = getNode(doc, "SAMLResponse")
	if val != "" {
		toreturn.RawData = []byte(val)
		val = strings.Replace(val, "&#x2b;", "+", -1)
		val = strings.Replace(val, "&#x3d;", "=", -1)
		data, err = base64.StdEncoding.DecodeString(val)
		if err != nil {
			return err
		}
	}
	err = xml.Unmarshal(data, &resp)
	if err != nil {
		return err
	}
	toreturn.SamlAss = val
	toreturn.Sresp = resp
	toreturn.OriData = string(body)
	o.SamlData = &toreturn
	return nil
}

func getNode(n *html.Node, name string) (val string, node *html.Node) {
	var isMatch bool
	if n.Type == html.ElementNode && n.Data == "input" {
		for _, a := range n.Attr {
			if a.Key == "name" && a.Val == name {
				isMatch = true
			}
			if a.Key == "value" && isMatch {
				val = a.Val
			}
		}
	}
	if node == nil || val == "" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			val, node = getNode(c, name)
			if val != "" {
				return
			}
		}
	}
	return
}
