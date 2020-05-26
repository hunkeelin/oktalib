package oktalib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// GetAwsCredentialsInput : The input for the function
type GetAwsCredentialsInput struct {
	RoleArn    string // RoleArn The Role ARN to assume. The user needs to have permission to assume the role in okta
	Expiration int64  // Expiration in seconds
}

// GetAwsCredentialsOutput : The output for the function
type GetAwsCredentialsOutput struct {
	AwsAccessKeyId     string // AwsAccessKeyId
	AwsSecretAccessKey string // AwsSecretAccessKey
	AwsSessionToken    string // AwsSessionToken
}

// GetAwsCredentials : Returns the secret,access and session token
func (o *OktaClient) GetAwsCredentials(i GetAwsCredentialsInput) (GetAwsCredentialsOutput, error) {
	// Get the saml assertion first.
	err := o.GetSamlAssertion()
	if err != nil {
		return GetAwsCredentialsOutput{}, err
	}
	samlSess := session.Must(session.NewSession())
	svc := sts.New(samlSess)
	samlParams := &sts.AssumeRoleWithSAMLInput{
		PrincipalArn:    aws.String(o.Principle),
		RoleArn:         aws.String(i.RoleArn),
		SAMLAssertion:   aws.String(string(o.SamlData.RawData)),
		DurationSeconds: aws.Int64(i.Expiration),
	}

	samlResp, err := svc.AssumeRoleWithSAML(samlParams)
	if err != nil {
		return GetAwsCredentialsOutput{}, err
	}
	return GetAwsCredentialsOutput{
		AwsAccessKeyId:     *samlResp.Credentials.AccessKeyId,
		AwsSecretAccessKey: *samlResp.Credentials.SecretAccessKey,
		AwsSessionToken:    *samlResp.Credentials.SessionToken,
	}, nil
}
