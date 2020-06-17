package oktalib

import (
	"fmt"
)

// NewInput returns an OktaClient object, it's for initialization purposes.
type NewInput struct {
	Org                 string
	IdentityProviderArn string
	SamlURI             string
}

// New returns Oktaclient with checks
func New(o *NewInput) (OktaClient, error) {
	if o.Org == "" {
		return OktaClient{}, fmt.Errorf("Please specify organization")
	}
	if o.IdentityProviderArn == "" {
		return OktaClient{}, fmt.Errorf("Please specify identify provider arn")
	}
	if o.SamlURI == "" {
		return OktaClient{}, fmt.Errorf("Please specify samlURI")
	}
	oktaBase := "https://" + o.Org + ".okta.com"
	return OktaClient{
		OktaUrl:    oktaBase + "/api/v1/authn",
		AwsSamlUrl: oktaBase + o.SamlURI,
		Principle:  o.IdentityProviderArn,
	}, nil
}
