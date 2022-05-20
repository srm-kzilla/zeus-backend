package helpers

import (
	"fmt"
	"reflect"
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

func ExistsInArray(Array interface{}, item interface{}) bool {
arr := reflect.ValueOf(Array)

if arr.Kind() != reflect.Array {
	println("Invalid data type")
}
for i := 0; i < arr.Len(); i++ {
	if arr.Index(i).Interface() == item {
		return true
	}
}

return false
}