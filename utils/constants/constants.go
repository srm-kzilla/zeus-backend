package constants

import "reflect"

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