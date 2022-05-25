package mailer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

func GenerateSESTemplate(input SESInput) (template *ses.SendEmailInput) {

	html := getHTMLTemplate(input.Name, input.TemplateName, input.EmbedData)
	title := input.Subject
	template = &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(input.RecieverEmail),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("utf-8"),
					Data:    aws.String(html),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    aws.String(title),
			},
		},
		Source: aws.String("SRMKZILLA <" + input.SenderEmail + ">"),
	}

	return
}
