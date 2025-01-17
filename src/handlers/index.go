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

var FilterOps = map[string]string{
	"n":  "NOT",
	"eq":  "=", // equal to
	"lt":  "<", // less than
	"gt":  ">", // greater than
	"lte": "<=", // less than or equal to
	"gte": ">=", // greater than or equal to
	"tft": "FULLTEXT", // full text search
	"ts":  "START", // start with
	"te":  "END", // end with
	"tx":   "EXACT", // exact match
}

// TODO ComposeFilter Compose filter map
func ComposeFilter(c fiber.Ctx, filter *map[string]any) error {
	var filters = strings.Split(c.Query("filter"), ",")
	// log.Printf("filters: %#v\n", filters)

	var filterMap = map[string]any{}

	for _, condition := range filters {
		var item = strings.Split(condition, ":")
		// log.Printf("item: %#v\n", item)

		var parts = strings.Split(item[len(item)-1], "_")
		var rule = fmt.Sprintf("%s %s", FilterOps[parts[0]], parts[len(parts)-1])

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
