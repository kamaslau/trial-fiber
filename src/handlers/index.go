package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// Pager Paging segments
type Pager struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

// TODO ComposeFilter Compose filter map
// Operators:
// COMPARE: n for not, nn for not null, e for equal, l for less, g for greater, le for less or equal, ge for greater or equal
// TEXT: tft for full text, ts for start, te for end, t for exact
func ComposeFilter(c fiber.Ctx, filter *map[string]any) error {
	var filters = strings.Split(c.Query("filter"), ",")
	// log.Printf("filters: %#v\n", filters)

	var filterMap = map[string]any{}

	for _, condition := range filters {
		var item = strings.Split(condition, ":")
		// log.Printf("item: %#v\n", item)

		var parts = strings.Split(item[len(item)-1], "_")
		var rule = fmt.Sprintf("%s %s", parts[0], parts[len(parts)-1])

		// log.Printf("item %s parts: %#v %s\n", item[0], parts, rule)
		// TODO Interpret the rule
		filterMap[item[0]] = rule
	}

	log.Printf("filterMap: %#v\n", filterMap)
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
