package mailer

type SESInput struct {
	TemplateName  string `json:"templateName"`
	Subject       string `json:"subject"`
	RecieverEmail string `json:"recieverEmail"`
	SenderEmail   string `json:"senderEmail"`
	Name          string `json:"name"`
}