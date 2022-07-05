package validators

import (
	"github.com/go-playground/validator/v10"
	authModel "github.com/srm-kzilla/events/api/auth/model"
	eventModel "github.com/srm-kzilla/events/api/events/model"
	inEventModel "github.com/srm-kzilla/events/api/inEvent/model"
	userModel "github.com/srm-kzilla/events/api/users/model"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

/***********************
Validates add user request.
***********************/
func ValidateUser(user userModel.User) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

/***********************
Validates add events request.
***********************/
func ValidateEvents(event eventModel.Event) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(event)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

/***********************
Validates add speaker request.
***********************/
func ValidateSpeaker(speaker eventModel.Speaker) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(speaker)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

/***********************
Validates user registration for event request.
***********************/
func ValidateRegisterUserReq(reqBody userModel.RegisterUserReq) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(reqBody)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

/***********************
Validates user RSVP for event request.
***********************/
func ValidateRsvpUserReq(reqBody userModel.RsvpUserReq) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(reqBody)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

/***********************
Validates register new admin request.
***********************/
func ValidateAdminUser(reqBody authModel.User) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(reqBody)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

/***********************
Validates user attendance for event request.
***********************/
func ValidateAttendanceQuery(reqBody inEventModel.AttendanceQuery) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(reqBody)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
