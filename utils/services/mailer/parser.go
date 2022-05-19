package mailer

import (
	"bytes"
	"io/ioutil"
	"log"
	"text/template"
)

// var TEMPLATES = map[string]string{
// 	"rsvpHtmlTemplate": templates.RsvpHtmlTemplate,
// 	"newUserHtmlTemplate": templates.NewUserHtmlTemplate,
// }

// var Template_Names = TemplateNames{
// RsvpTemplate: "rsvpHtmlTemplate",
// NewUserTemplate: "newUserHtmlTemplate",
// }

var TEMPLATES = TemplateNames{
	RsvpTemplate: "rsvpEmailTemplate.html",
	NewUserTemplate: "newUserTemplate.html",
}

func getHTMLTemplate(name string, templateName string, embedData interface{}) string {
	var templateBuffer bytes.Buffer


	htmlData, err := ioutil.ReadFile("./templates/" + templateName)
	if err != nil {
		log.Print("Error Reading the file")
	}

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))
	// htmlTemplate := template.Must(template.New("email.html").Parse(TEMPLATES[templateName]))

	error := htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", embedData)


	if error != nil {
		log.Println(error)
		return ""
	}

	return templateBuffer.String()
}
