package ip

import (
	"context"
	"strings"

	"github.com/ViBiOh/httputils/v3/pkg/request"
)

// Get returns current IP
func Get() (string, error) {
	response, err := request.New().Get("https://ifconfig.co").Send(context.Background(), nil)
	if err != nil {
		return "", err
	}

	content, err := request.ReadBodyResponse(response)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}
