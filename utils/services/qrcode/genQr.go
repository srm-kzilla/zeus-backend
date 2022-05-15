package qr

import (
	"log"

	qrcode "github.com/skip2/go-qrcode"
)

func GenerateQRCode(data string) {

	err := qrcode.WriteFile(data, qrcode.Medium, 256, "qr.svg")
	// png, err := qrcode.Encode("https://example.org", qrcode.Medium, 256)

	if err != nil {
		log.Println("Error: ", err)
	}

}
func GenerateQR(data string) string {
	qrCode := "https://chart.apis.google.com/chart?cht=qr&chs=256x256&chl=" + data
	return qrCode
}
