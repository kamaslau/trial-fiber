package handlers

import (
	"log"
	"net/http"
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
	if c.Query("filter") == "" {
		return nil
	}

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

var httpMessages = map[int]string{
	http.StatusNoContent:           "No Content",
	http.StatusBadRequest:          "Bad Request",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusForbidden:           "Forbidden",
	http.StatusNotFound:            "Not Found",
	http.StatusUnprocessableEntity: "Unprocessable Entity",
	http.StatusInternalServerError: "Internal Server Error",
}

func GetHTTPMsg(code int) map[string]any {
	message, exists := httpMessages[code]
	if !exists {
		message = http.StatusText(code)
	}
	return fiber.Map{
		"code":    code,
		"message": message,
	}
}
