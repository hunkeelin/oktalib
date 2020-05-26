package main

import (
	"fmt"
	"github.com/hunkeelin/oktalib"
	"github.com/hunkeelin/userprompt"
	"net/http/cookiejar"
	"os"
)

func main() {
	o, err := oktalib.New(&oktalib.NewInput{
		Org:                 "dev-815627",
		IdentityProviderArn: "arn:aws:iam::216228501626:saml-provider/Okta_2",
		SamlURI:             "/app/amazon_aws/exkawa67iQIlhKIxE4x6",
	})
	if err != nil {
		panic(err)
	}
	currentUser := os.Getenv("USER")
	userName, err := userprompt.UserPromptWithDefault("Enter Okta Username ("+currentUser+")", currentUser, false)
	if err != nil {
		panic(err)
	}

	pass, err := userprompt.UserPrompt("Enter Okta Password", true)
	if err != nil {
		panic(err)
	}
	cJar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	o.Username = userName
	o.Password = pass
	o.CookieJar = cJar
	err = o.LdapLogin()
	if err != nil {
		panic(err)
	}

	if len(o.UserAuth.Embedded.Factors) == 0 || len(o.UserAuth.Embedded.Factors) < 1 {
		panic(fmt.Errorf("Extra verification must be enabled in Okta. Visit https://varomoney.okta.com/enduser/settings."))
	}

	switch {
	case searchAuthMethod(o.UserAuth.Embedded.Factors, oktalib.YubiKey):
		fmt.Println("Congrats on your shiny new Yubikey")
		code, err := userprompt.UserPrompt("Please give it a squeeze", false)
		if err != nil {
			panic(err)
		}
		err = o.OktaMfa(oktalib.YubiKey, code)
		if err != nil {
			panic(err)
		}
	case searchAuthMethod(o.UserAuth.Embedded.Factors, oktalib.MfaPush):
		err = o.OktaMfa(oktalib.MfaPush, "")
		if err != nil {
			panic(err)
		}
	case searchAuthMethod(o.UserAuth.Embedded.Factors, oktalib.MfaCode):
		passcode, err := userprompt.UserPrompt("Enter a token from your mobile authenticator app", false)
		if err != nil {
			panic(err)
		}
		err = o.OktaMfa(oktalib.MfaCode, passcode)
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Errorf("No recongized mfa method found please contact your administrator."))
	}
	out, err := o.GetAwsCredentials(oktalib.GetAwsCredentialsInput{
		RoleArn:    "arn:aws:iam::216228501626:role/devops-admin-role",
		Expiration: 28800,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	return
}
func searchAuthMethod(sep []oktalib.OktaUserAuthnFactor, s string) bool {
	for _, i := range sep {
		if i.FactorType == s {
			return true
		}
	}
	return false
}
