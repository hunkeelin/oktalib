package oktalib

import (
	"net/http"

	"encoding/xml"
)

// OktaClient struct for methods
type OktaClient struct {
	Principle      string
	Role           string
	SamlData       *samlresp
	Username       string
	Password       string
	UserAuth       *OktaUserAuthn
	OktaAwsSAMLUrl string
	CookieJar      http.CookieJar
	OktaUrl        string
	AwsSamlUrl     string
}

type oktaUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type oktaStateToken struct {
	StateToken string `json:"stateToken"`
	PassCode   string `json:"passCode"`
}

// OktaUserAuthn a struct
type OktaUserAuthn struct {
	StateToken   string                `json:"stateToken"`
	SessionToken string                `json:"sessionToken"`
	ExpiresAt    string                `json:"expiresAt"`
	Status       string                `json:"status"`
	Embedded     OktaUserAuthnEmbedded `json:"_embedded"`
	FactorResult string                `json:"factorResult"`
	CookieJar    http.CookieJar
}

// OktaUserAuthnEmbedded struct
type OktaUserAuthnEmbedded struct {
	Factors []OktaUserAuthnFactor `json:"factors"`
	Factor  OktaUserAuthnFactor   `json:"factor"`
}

// OktaUserAuthnFactor struct
type OktaUserAuthnFactor struct {
	Id         string                      `json:"id"`
	FactorType string                      `json:"factorType"`
	Provider   string                      `json:"provider"`
	Embedded   OktaUserAuthnFactorEmbedded `json:"_embedded"`
}

// OktaUserAuthnFactorEmbedded struct
type OktaUserAuthnFactorEmbedded struct {
	Verification OktaUserAuthnFactorEmbeddedVerification `json:"verification"`
	Challange    OktaUserAuthnFactorEmbeddedChallenge    `json:"challenge"`
}

// OktaUserAuthnFactorEmbeddedChallenge
type OktaUserAuthnFactorEmbeddedChallenge struct {
	CorrectAnswer int `json:"correctAnswer"`
}

// OktaUserAuthnFactorEmbeddedVerification struct
type OktaUserAuthnFactorEmbeddedVerification struct {
	Host         string                                       `json:"host"`
	Signature    string                                       `json:"signature"`
	FactorResult string                                       `json:"factorResult"`
	Links        OktaUserAuthnFactorEmbeddedVerificationLinks `json:"_links"`
}

// OktaUserAuthnFactorEmbeddedVerificationLinks struct
type OktaUserAuthnFactorEmbeddedVerificationLinks struct {
	Complete OktaUserAuthnFactorEmbeddedVerificationLinksComplete `json:"complete"`
}

// OktaUserAuthnFactorEmbeddedVerificationLinksComplete struct
type OktaUserAuthnFactorEmbeddedVerificationLinksComplete struct {
	Href string `json:"href"`
}

// Response struct
type Response struct {
	XMLName      xml.Name
	SAMLP        string `xml:"xmlns:samlp,attr"`
	SAML         string `xml:"xmlns:saml,attr"`
	SAMLSIG      string `xml:"xmlns:samlsig,attr"`
	Destination  string `xml:"Destination,attr"`
	ID           string `xml:"ID,attr"`
	Version      string `xml:"Version,attr"`
	IssueInstant string `xml:"IssueInstant,attr"`
	InResponseTo string `xml:"InResponseTo,attr"`

	Assertion Assertion `xml:"Assertion"`
	Status    Status    `xml:"Status"`

	originalString string
}

// Assertion struct
type Assertion struct {
	XMLName            xml.Name
	ID                 string `xml:"ID,attr"`
	Version            string `xml:"Version,attr"`
	XS                 string `xml:"xmlns:xs,attr"`
	XSI                string `xml:"xmlns:xsi,attr"`
	SAML               string `xml:"saml,attr"`
	IssueInstant       string `xml:"IssueInstant,attr"`
	Subject            Subject
	Conditions         Conditions
	AttributeStatement AttributeStatement
}

// Conditions struct
type Conditions struct {
	XMLName      xml.Name
	NotBefore    string `xml:",attr"`
	NotOnOrAfter string `xml:",attr"`
}

// Subject struct
type Subject struct {
	XMLName             xml.Name
	NameID              NameID
	SubjectConfirmation SubjectConfirmation
}

// SubjectConfirmation struct
type SubjectConfirmation struct {
	XMLName                 xml.Name
	Method                  string `xml:",attr"`
	SubjectConfirmationData SubjectConfirmationData
}

// Status struct
type Status struct {
	XMLName    xml.Name
	StatusCode StatusCode `xml:"StatusCode"`
}

// SubjectConfirmationData struct
type SubjectConfirmationData struct {
	InResponseTo string `xml:",attr"`
	NotOnOrAfter string `xml:",attr"`
	Recipient    string `xml:",attr"`
}

// NameID struct
type NameID struct {
	XMLName xml.Name
	Format  string `xml:",attr"`
	Value   string `xml:",innerxml"`
}

// StatusCode struct
type StatusCode struct {
	XMLName xml.Name
	Value   string `xml:",attr"`
}

// AttributeValue struct
type AttributeValue struct {
	XMLName xml.Name
	Type    string `xml:"xsi:type,attr"`
	Value   string `xml:",innerxml"`
}

// Attribute struct
type Attribute struct {
	XMLName         xml.Name
	Name            string           `xml:",attr"`
	FriendlyName    string           `xml:",attr"`
	NameFormat      string           `xml:",attr"`
	AttributeValues []AttributeValue `xml:"AttributeValue"`
}

// AttributeStatement struct
type AttributeStatement struct {
	XMLName    xml.Name
	Attributes []Attribute `xml:"Attribute"`
}
