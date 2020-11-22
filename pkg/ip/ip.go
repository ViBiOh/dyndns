package ip

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"

	"github.com/ViBiOh/httputils/v3/pkg/request"
)

// Get returns current IP
func Get(wantedNetwork string) (string, error) {
	req, err := request.New().Get("https://ifconfig.co").Build(context.Background(), nil)
	if err != nil {
		return "", err
	}

	httpClient := http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   15 * time.Second,
				KeepAlive: 15 * time.Second,
				DualStack: false,
				Control: func(network, address string, c syscall.RawConn) error {
					if network == wantedNetwork {
						return nil
					}

					return fmt.Errorf("only want %s", wantedNetwork)
				},
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       15 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	response, err := request.DoWithClientAndRetry(httpClient, req, 3)
	if err != nil {
		return "", err
	}

	content, err := request.ReadBodyResponse(response)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}
