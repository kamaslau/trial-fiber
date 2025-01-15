package handlers

import "github.com/gofiber/fiber/v3"

var HTTPStatus = map[uint16]string{
	200: "OK",
	201: "Created",
	204: "No Content",
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	422: "Unprocessable Entity",
	500: "Internal Server Error",
}

func GetHTTPStatus(code uint16) map[string]any {
	return fiber.Map{
		"status":  code,
		"message": HTTPStatus[code],
	}
}
