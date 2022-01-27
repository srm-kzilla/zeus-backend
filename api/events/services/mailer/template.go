package mailer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/srm-kzilla/events/api/events/model"
)

func GenerateSESTemplate(input model.SESInput) (template *ses.SendEmailInput) {

	html := getHTMLTemplate(input.Name, input.TemplateName)

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
