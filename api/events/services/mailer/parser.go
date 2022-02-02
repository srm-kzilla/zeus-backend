package mailer

import (
	"bytes"
	"io/ioutil"
	"log"
	"text/template"
)

func getHTMLTemplate(name string, templateName string) string {
	var templateBuffer bytes.Buffer

	type EmailData struct {
		Name string
	}

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
		log.Println(err)
		return ""
	}

	return templateBuffer.String()
}
