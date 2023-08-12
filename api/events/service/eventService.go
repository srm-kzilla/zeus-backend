package eventService

import (
	"strings"

	userModel "github.com/srm-kzilla/events/api/users/model"
	"github.com/srm-kzilla/events/utils/services"
)

func GetUsersByCollegeYear(users []userModel.User) (totalRegistrations, firstYears, secondYears, thirdYears, fourthYears int, err error) {
	totalRegistrations = len(users)
	firstYearPrefix, secondYearPrefix, thirdYearPrefix, fourthYearPrefix, err := services.GenerateCollegeYearRegistrationPrefix()
	for _, user := range users {
		regNumber := strings.ToUpper(user.RegNumber)
		switch {
		case strings.HasPrefix(regNumber, fourthYearPrefix):
			fourthYears++
		case strings.HasPrefix(regNumber, thirdYearPrefix):
			thirdYears++
		case strings.HasPrefix(regNumber, secondYearPrefix):
			secondYears++
		case strings.HasPrefix(regNumber, firstYearPrefix):
			firstYears++
		}
	}
	return totalRegistrations, firstYears, secondYears, thirdYears, fourthYears, err
}
