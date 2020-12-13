package utils

import (
	"fmt"
	"net/http"
)

func HTTPStatusCheck(response *http.Response) error {
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("Invalid status code in HTTP response: %d", response.StatusCode)
	}
	return nil
}
