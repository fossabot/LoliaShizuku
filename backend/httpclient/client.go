package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"loliashizuku/backend/models"
	"loliashizuku/backend/version"
)

var ErrUnauthorized = errors.New("center api unauthorized")

type APIError struct {
	Path       string
	StatusCode int
	Code       int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api %s failed: status=%d code=%d msg=%s", e.Path, e.StatusCode, e.Code, e.Message)
}

type Options struct {
	BaseURL        string
	HTTPClient     *http.Client
	UserAgent      string
	GetAccessToken func(ctx context.Context) (string, error)
	OnUnauthorized func(ctx context.Context) error
}

type Client struct {
	baseURL        string
	httpClient     *http.Client
	userAgent      string
	getAccessToken func(ctx context.Context) (string, error)
	onUnauthorized func(ctx context.Context) error
}

func ResolveUserAgent(userAgent string) string {
	resolved := strings.TrimSpace(userAgent)
	if resolved == "" {
		resolved = strings.TrimSpace(os.Getenv("LOLIA_HTTP_USER_AGENT"))
	}
	if resolved == "" {
		resolved = version.UserAgent()
	}
	return resolved
}

func New(options Options) *Client {
	userAgent := ResolveUserAgent(options.UserAgent)

	return &Client{
		baseURL:        strings.TrimRight(strings.TrimSpace(options.BaseURL), "/"),
		httpClient:     options.HTTPClient,
		userAgent:      userAgent,
		getAccessToken: options.GetAccessToken,
		onUnauthorized: options.OnUnauthorized,
	}
}

func (c *Client) DoJSON(
	ctx context.Context,
	method, path string,
	query map[string]string,
	body any,
	dest any,
) error {
	requestURL, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("build request url: %w", err)
	}

	queryValues := requestURL.Query()
	for key, value := range query {
		if strings.TrimSpace(value) != "" {
			queryValues.Set(key, value)
		}
	}
	requestURL.RawQuery = queryValues.Encode()

	var reqBody io.Reader
	if body != nil {
		payload, marshalErr := json.Marshal(body)
		if marshalErr != nil {
			return fmt.Errorf("marshal request body for %s: %w", path, marshalErr)
		}
		reqBody = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL.String(), reqBody)
	if err != nil {
		return fmt.Errorf("build request for %s: %w", path, err)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.getAccessToken != nil {
		accessToken, tokenErr := c.getAccessToken(ctx)
		if tokenErr != nil {
			return tokenErr
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	payload, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return fmt.Errorf("read response body for %s: %w", path, readErr)
	}

	var envelope models.APIEnvelope
	if len(payload) > 0 {
		if err := json.Unmarshal(payload, &envelope); err != nil {
			return fmt.Errorf("decode response for %s: %w", path, err)
		}
	}

	businessCode := envelope.Code
	if businessCode == 0 {
		businessCode = envelope.Status
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden ||
		businessCode == http.StatusUnauthorized || businessCode == http.StatusForbidden {
		if c.onUnauthorized != nil {
			_ = c.onUnauthorized(ctx)
		}
		apiErr := &APIError{
			Path:       path,
			StatusCode: resp.StatusCode,
			Code:       businessCode,
			Message:    envelope.Msg,
		}
		return errors.Join(ErrUnauthorized, apiErr)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{
			Path:       path,
			StatusCode: resp.StatusCode,
			Code:       businessCode,
			Message:    envelope.Msg,
		}
	}

	if businessCode != 0 && businessCode != 200 {
		return &APIError{
			Path:       path,
			StatusCode: resp.StatusCode,
			Code:       businessCode,
			Message:    envelope.Msg,
		}
	}

	if dest == nil {
		return nil
	}
	if len(envelope.Data) == 0 || string(envelope.Data) == "null" {
		return nil
	}
	if err := json.Unmarshal(envelope.Data, dest); err != nil {
		return fmt.Errorf("decode response data for %s: %w", path, err)
	}
	return nil
}
