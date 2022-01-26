package mailer

import (
	"bytes"
	"io/ioutil"
	"log"
	"text/template"
)

func getHTMLTemplate(name string, templateName string) string {
	var templateBuffer bytes.Buffer

	// can use this type to add some additional data to the email template
	type EmailData struct {
		Name string
	}

	// You can bind custom data here as per requirements.

	data := EmailData{
		Name: name,
	}

	htmlData, err := ioutil.ReadFile("./templates/" + templateName)
	if err != nil {
		log.Print("Error Reading the file")
	}

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", data)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	return templateBuffer.String()
}
