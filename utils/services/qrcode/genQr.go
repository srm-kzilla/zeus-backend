package qr

import (
	"log"

	qrcode "github.com/skip2/go-qrcode"
)

func GenerateQRCode(data string) {
	err := qrcode.WriteFile(data, qrcode.Medium, 256, "qr.png")
	
	if err != nil {
		log.Println("Error: ", err)
	}

}