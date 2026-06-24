package handler

import (
	"net/http"
	"net/url"
	"time"

	"enterprise-order-management-api/internal/config"
	"enterprise-order-management-api/internal/dto"
	appmiddleware "enterprise-order-management-api/internal/middleware"
	"enterprise-order-management-api/internal/oauth"
	"enterprise-order-management-api/internal/pkg/response"
	"enterprise-order-management-api/internal/service"

	"github.com/labstack/echo/v4"
)

const googleOAuthStateCookie = "google_oauth_state"

type AuthHandler struct {
	auth                service.AuthService
	frontendCallbackURL string
	stateSecret         string
}

func NewAuthHandler(auth service.AuthService, cfg config.Config) *AuthHandler {
	return &AuthHandler{
		auth:                auth,
		frontendCallbackURL: cfg.FrontendAuthCallbackURL,
		stateSecret:         cfg.OAuthStateSecret,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.auth.Register(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return response.Created(c, res)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.auth.Login(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	loginURL, state, err := h.auth.BeginGoogleLogin(c.Request().Context())
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, h.buildFrontendCallbackURL(url.Values{
			"status": {"error"},
			"error":  {err.Error()},
		}))
	}

	c.SetCookie(&http.Cookie{
		Name:     googleOAuthStateCookie,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((10 * time.Minute).Seconds()),
	})

	return c.Redirect(http.StatusTemporaryRedirect, loginURL)
}

func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	if googleError := c.QueryParam("error"); googleError != "" {
		return c.Redirect(http.StatusTemporaryRedirect, h.buildFrontendCallbackURL(url.Values{
			"status": {"error"},
			"error":  {googleError},
		}))
	}

	state := c.QueryParam("state")
	stateCookie, err := c.Cookie(googleOAuthStateCookie)
	if err != nil || state == "" || stateCookie.Value == "" || stateCookie.Value != state {
		h.clearOAuthStateCookie(c)
		return c.Redirect(http.StatusTemporaryRedirect, h.buildFrontendCallbackURL(url.Values{
			"status": {"error"},
			"error":  {"Invalid OAuth state"},
		}))
	}
	if _, err := oauth.ParseState(state, h.stateSecret, oauth.GoogleProvider); err != nil {
		h.clearOAuthStateCookie(c)
		return c.Redirect(http.StatusTemporaryRedirect, h.buildFrontendCallbackURL(url.Values{
			"status": {"error"},
			"error":  {"Invalid OAuth state"},
		}))
	}

	code := c.QueryParam("code")
	if code == "" {
		h.clearOAuthStateCookie(c)
		return c.Redirect(http.StatusTemporaryRedirect, h.buildFrontendCallbackURL(url.Values{
			"status": {"error"},
			"error":  {"Missing OAuth code"},
		}))
	}

	authRes, err := h.auth.CompleteGoogleLogin(c.Request().Context(), code)
	h.clearOAuthStateCookie(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, h.buildFrontendCallbackURL(url.Values{
			"status": {"error"},
			"error":  {err.Error()},
		}))
	}

	return c.Redirect(http.StatusTemporaryRedirect, h.buildFrontendCallbackURL(url.Values{
		"status":        {"success"},
		"access_token":  {authRes.AccessToken},
		"refresh_token": {authRes.RefreshToken},
	}))
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	var req dto.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.auth.Refresh(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	var req dto.LogoutRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	if err := h.auth.Logout(c.Request().Context(), req.RefreshToken); err != nil {
		return err
	}
	return response.Message(c, http.StatusOK, "Logged out successfully")
}

func (h *AuthHandler) Me(c echo.Context) error {
	res, err := h.auth.Me(c.Request().Context(), appmiddleware.CurrentUserID(c))
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *AuthHandler) buildFrontendCallbackURL(values url.Values) string {
	callbackURL, err := url.Parse(h.frontendCallbackURL)
	if err != nil {
		return h.frontendCallbackURL
	}
	callbackURL.Fragment = values.Encode()
	return callbackURL.String()
}

func (h *AuthHandler) clearOAuthStateCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     googleOAuthStateCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}
