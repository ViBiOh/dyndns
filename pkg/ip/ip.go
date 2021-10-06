package ip

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"

	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/httputils/v4/pkg/request"
)

// Get returns current IP
func Get(url, wantedNetwork string) (string, error) {
	req, err := request.New().Get(url).Build(context.Background(), nil)
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

	for i := 0; i < 3; i++ {
		response, err := request.DoWithClient(&httpClient, req)
		if err != nil {
			logger.Error("attempt #%d failed with error: %s", i+1, err)
			time.Sleep(time.Second)
			continue
		}

		content, err := request.ReadBodyResponse(response)
		if err != nil {
			return "", err
		}

		return strings.TrimSpace(string(content)), nil
	}

	return "", errors.New("unable to get current IP")
}
