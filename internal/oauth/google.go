package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const GoogleProvider = "google"

type GoogleUserInfo struct {
	ProviderUserID string
	Email          string
	Name           string
	AvatarURL      string
	EmailVerified  bool
}

type GoogleProviderClient interface {
	Enabled() bool
	AuthCodeURL(state string) string
	Exchange(ctx context.Context, code string) (*GoogleUserInfo, error)
}

type disabledGoogleProvider struct{}

func (disabledGoogleProvider) Enabled() bool { return false }

func (disabledGoogleProvider) AuthCodeURL(string) string { return "" }

func (disabledGoogleProvider) Exchange(context.Context, string) (*GoogleUserInfo, error) {
	return nil, fmt.Errorf("google oauth is not configured")
}

type googleProvider struct {
	config     *oauth2.Config
	httpClient *http.Client
}

func NewGoogleProvider(clientID string, clientSecret string, redirectURL string) GoogleProviderClient {
	if strings.TrimSpace(clientID) == "" || strings.TrimSpace(clientSecret) == "" || strings.TrimSpace(redirectURL) == "" {
		return disabledGoogleProvider{}
	}

	return &googleProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
		httpClient: http.DefaultClient,
	}
}

func (p *googleProvider) Enabled() bool {
	return p != nil && p.config != nil
}

func (p *googleProvider) AuthCodeURL(state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (p *googleProvider) Exchange(ctx context.Context, code string) (*GoogleUserInfo, error) {
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := p.config.Client(ctx, token)
	if p.httpClient != nil {
		client = &http.Client{
			Transport: client.Transport,
			Timeout:   p.httpClient.Timeout,
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://openidconnect.googleapis.com/v1/userinfo", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo returned status %d", resp.StatusCode)
	}

	var payload struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &GoogleUserInfo{
		ProviderUserID: payload.Sub,
		Email:          strings.TrimSpace(strings.ToLower(payload.Email)),
		Name:           strings.TrimSpace(payload.Name),
		AvatarURL:      strings.TrimSpace(payload.Picture),
		EmailVerified:  payload.EmailVerified,
	}, nil
}
