package helpers

import (
	"fmt"
	"net/http"
	"strconv"
)

// ParsePathInt64 extracts and parses an int64 from URL path parameters
func ParsePathInt64(r *http.Request, key string) (int64, error) {
	value := r.PathValue(key)
	if value == "" {
		return 0, fmt.Errorf("missing path parameter: %s", key)
	}
	
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: must be a number", key)
	}
	
	return intValue, nil
}

