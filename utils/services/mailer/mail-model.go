package mailer

type SESInput struct {
	TemplateName  string `json:"templateName"`
	Subject       string `json:"subject"`
	RecieverEmail string `json:"recieverEmail"`
	SenderEmail   string `json:"senderEmail"`
	Name          string `json:"name"`
	EmbedData	  interface{} `json:"embedData"`
}
type TemplateNames struct {
	RsvpTemplate	string
	NewUserTemplate	string
}