package qr

/***********************
Generates a QR Code of any string data.
***********************/
func GenerateQRCode(data string) string {
	qrCode := "https://chart.apis.google.com/chart?cht=qr&chs=256x256&chl=" + data
	return qrCode
}
