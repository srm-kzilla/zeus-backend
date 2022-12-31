package constants

import (
	userModel "github.com/srm-kzilla/events/api/users/model"
)

var Animations = userModel.Animation{
	RsvpSuccess:       "https://assets2.lottiefiles.com/packages/lf20_znxedwj6.json",
	EventCompleted:    "https://assets8.lottiefiles.com/packages/lf20_rbbibjz5.json",
	EventDoesNotExist: "https://assets8.lottiefiles.com/packages/lf20_rbbibjz5.json",
	AlreadyRsvpd:      "https://assets2.lottiefiles.com/packages/lf20_znxedwj6.json",
	FullyBooked:       "https://assets8.lottiefiles.com/packages/lf20_rbbibjz5.json",
}
