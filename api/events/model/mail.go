package model

// SESInput ... contains TemplateName, Subject, Email (Reciever Email), Name (Reciever Name)
type SESInput struct {
	TemplateName  string `json:"templateName"`
	Subject       string `json:"subject"`
	RecieverEmail string `json:"recieverEmail"`
	SenderEmail   string `json:"senderEmail"`
	Name          string `json:"name"`
}
