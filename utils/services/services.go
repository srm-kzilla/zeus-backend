package services

import (
	"os"
	"strconv"
)

func GenerateCollegeYearRegistrationPrefix() (firstYearPrefix, secondYearPrefix, thirdYearPrefix, fourthYearPrefix string, err error) {

	currentFirstYear, err := strconv.Atoi(os.Getenv("FIRST_YEAR"))
	if err != nil {
		return "", "", "", "", err
	}
	firstYearPrefix = "RA" + strconv.Itoa(currentFirstYear)
	secondYearPrefix = "RA" + strconv.Itoa(currentFirstYear-1)
	thirdYearPrefix = "RA" + strconv.Itoa(currentFirstYear-2)
	fourthYearPrefix = "RA" + strconv.Itoa(currentFirstYear-3)
	return firstYearPrefix, secondYearPrefix, thirdYearPrefix, fourthYearPrefix, nil
}
