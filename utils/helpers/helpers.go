package helpers

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid"
)

func GenerateNanoID(size int) string {
var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	nanoID, err := gonanoid.Generate(alphabet, size)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return nanoID
}