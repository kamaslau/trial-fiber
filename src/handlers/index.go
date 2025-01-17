package handlers

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// Pager Paging segments
type Pager struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

// ComposeFilter Compose filter map
func ComposeFilter(c fiber.Ctx, filter *map[string]any) error {
	var filters = strings.Split(c.Query("filter"), ",")
	// log.Printf("filters: %#v\n", filters)

	for _, condition := range filters {
		var item = strings.Split(condition, ":")
		// log.Printf("item: %#v\n", item)

		(*filter)[item[0]] = item[len(item)-1]
	}

	log.Printf("filter: %#v\n", filter)
	return nil
}

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
