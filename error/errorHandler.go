package errors

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

//HTTPError ::::
type HTTPError struct {
	Status int
	Err    error
}

//NewHTTPError ::::
func NewHTTPError(status int, err error) HTTPError {
	return HTTPError{Status: status, Err: err}
}

func (he HTTPError) Error() string {
	return fmt.Sprintf("error %d %s (%s)", he.Status, he.Err.Error(), http.StatusText(he.Status))
}

//ErrorHandler ::::
func ErrorHandler() fiber.Handler {
	return (func(c *fiber.Ctx) error {
		err := c.Context().Err()
		if err == nil { // no error
			c.Next()
			return nil
		}
		status := http.StatusInternalServerError //default error status
		if e, ok := err.(HTTPError); ok {        // it's a custom error, so use the status in the error
			status = e.Status
		}
		msg := map[string]interface{}{
			"status":      status,
			"status_text": http.StatusText(status),
			"error_msg":   err.Error(),
		}

		c.Set("Content-Type", "application/json")
		c.Status(status).JSON(msg)

		return nil
	})
}
