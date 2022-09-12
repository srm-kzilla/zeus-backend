package mailer

import (
	"bytes"
	"io/ioutil"
	"log"
	"text/template"
)

var TEMPLATES = TemplateNames{
	RsvpTemplate:    "rsvpEmailTemplate.html",
	NewUserTemplate: "newUserTemplate.html",
}

/**************************
Gets the the HTML template.
**************************/
func getHTMLTemplate(templateName string, embedData interface{}) string {
	var templateBuffer bytes.Buffer

	htmlData, err := ioutil.ReadFile("./templates/" + templateName)
	if err != nil {
		log.Print("Error Reading the file")
	}

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))

	error := htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", embedData)

	if error != nil {
		log.Println(error)
		return ""
	}

	return templateBuffer.String()
}
