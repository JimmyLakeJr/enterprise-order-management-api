package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"enterprise-order-management-api/internal/config"
)

type ZaloPayClient interface {
	CreateOrder(ctx context.Context, req ZaloPayCreateOrderRequest) (*ZaloPayCreateOrderResponse, string, error)
	QueryOrder(ctx context.Context, appTransID string) (*ZaloPayQueryOrderResponse, string, error)
	VerifyCallback(data string, mac string) bool
}

type ZaloPayCreateOrderRequest struct {
	AppUser       string
	AppTransID    string
	Amount        int64
	Description   string
	Items         []map[string]any
	RedirectURL   string
	CallbackURL   string
	DefaultBank   string
	ExpiresAfter  time.Duration
	ReferenceData map[string]any
}

type ZaloPayCreateOrderResponse struct {
	ReturnCode       int    `json:"return_code"`
	ReturnMessage    string `json:"return_message"`
	SubReturnCode    int    `json:"sub_return_code"`
	SubReturnMessage string `json:"sub_return_message"`
	ZPTransToken     string `json:"zp_trans_token"`
	OrderURL         string `json:"order_url"`
	OrderToken       string `json:"order_token"`
	QRCode           string `json:"qr_code"`
}

type ZaloPayQueryOrderResponse struct {
	ReturnCode       int    `json:"return_code"`
	ReturnMessage    string `json:"return_message"`
	SubReturnCode    int    `json:"sub_return_code"`
	SubReturnMessage string `json:"sub_return_message"`
	IsProcessing     bool   `json:"is_processing"`
	Amount           int64  `json:"amount"`
	ZPTransID        int64  `json:"zp_trans_id"`
	ServerTime       int64  `json:"server_time"`
	DiscountAmount   int64  `json:"discount_amount"`
}

type zaloPayClient struct {
	httpClient *http.Client
	cfg        config.ZaloPayConfig
}

func NewZaloPayClient(cfg config.ZaloPayConfig) ZaloPayClient {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &zaloPayClient{
		httpClient: &http.Client{Timeout: timeout},
		cfg:        cfg,
	}
}

func (c *zaloPayClient) CreateOrder(ctx context.Context, req ZaloPayCreateOrderRequest) (*ZaloPayCreateOrderResponse, string, error) {
	now := time.Now()
	itemsJSON, err := json.Marshal(req.Items)
	if err != nil {
		return nil, "", err
	}

	embedData := map[string]any{
		"redirecturl": req.RedirectURL,
	}
	for key, value := range req.ReferenceData {
		embedData[key] = value
	}
	embedDataJSON, err := json.Marshal(embedData)
	if err != nil {
		return nil, "", err
	}

	payload := map[string]any{
		"app_id":                  c.cfg.AppID,
		"app_user":                req.AppUser,
		"app_trans_id":            req.AppTransID,
		"app_time":                now.UnixMilli(),
		"amount":                  req.Amount,
		"description":             req.Description,
		"callback_url":            req.CallbackURL,
		"item":                    string(itemsJSON),
		"embed_data":              string(embedDataJSON),
		"bank_code":               req.DefaultBank,
		"expire_duration_seconds": maxInt64(int64(req.ExpiresAfter/time.Second), 300),
	}

	hmacInput := fmt.Sprintf(
		"%d|%s|%s|%d|%v|%s|%s",
		c.cfg.AppID,
		req.AppTransID,
		req.AppUser,
		req.Amount,
		payload["app_time"],
		payload["embed_data"],
		payload["item"],
	)
	payload["mac"] = signHMACSHA256(hmacInput, c.cfg.Key1)

	rawRequest, err := json.Marshal(payload)
	if err != nil {
		return nil, "", err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.CreateOrderEndpoint, bytes.NewReader(rawRequest))
	if err != nil {
		return nil, "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpRes, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, "", err
	}
	defer httpRes.Body.Close()

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, "", err
	}

	var response ZaloPayCreateOrderResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, string(body), err
	}

	return &response, string(body), nil
}

func (c *zaloPayClient) QueryOrder(ctx context.Context, appTransID string) (*ZaloPayQueryOrderResponse, string, error) {
	form := url.Values{}
	form.Set("app_id", fmt.Sprintf("%d", c.cfg.AppID))
	form.Set("app_trans_id", appTransID)
	form.Set("mac", signHMACSHA256(fmt.Sprintf("%d|%s|%s", c.cfg.AppID, appTransID, c.cfg.Key1), c.cfg.Key1))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.QueryOrderEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, "", err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpRes, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, "", err
	}
	defer httpRes.Body.Close()

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, "", err
	}

	var response ZaloPayQueryOrderResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, string(body), err
	}

	return &response, string(body), nil
}

func (c *zaloPayClient) VerifyCallback(data string, mac string) bool {
	return secureEqual(signHMACSHA256(data, c.cfg.Key2), mac)
}

func signHMACSHA256(message string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func secureEqual(expected string, actual string) bool {
	return hmac.Equal([]byte(strings.ToLower(expected)), []byte(strings.ToLower(actual)))
}

func maxInt64(value int64, fallback int64) int64 {
	if value > fallback {
		return value
	}
	return fallback
}
