package mailer

import (
	"bytes"
	"log"
	"text/template"

	"github.com/srm-kzilla/events/templates"
)

var TEMPLATES = map[string]string{
	"rsvpHtmlTemplate": templates.RsvpHtmlTemplate,
	"newUserHtmlTemplate": templates.NewUserHtmlTemplate,
}

var Template_Names = TemplateNames{
RsvpTemplate: "rsvpHtmlTemplate",
NewUserTemplate: "newUserHtmlTemplate",
}

func getHTMLTemplate(name string, templateName string, embedData interface{}) string {
	var templateBuffer bytes.Buffer

	// type EmailData struct {
	// 	Name string
	// }

	// data := EmailData{
	// 	Name: name,
	// }

	// htmlData, err := ioutil.ReadFile("./templates/" + templateName)
	// if err != nil {
	// 	log.Print("Error Reading the file")
	// }

	// htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))
	htmlTemplate := template.Must(template.New("email.html").Parse(TEMPLATES[templateName]))

	err := htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", embedData)


	if err != nil {
		log.Println(err)
		return ""
	}

	return templateBuffer.String()
}
