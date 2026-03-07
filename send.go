package sdk

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// SendOptions configures a Send call.
type SendOptions struct {
	// RoostID is the target roost.
	RoostID string
	// EventType is set as the X-Pigeon-Event-Type header.
	EventType string
	// Payload is the JSON body to send.
	Payload any
	// SigningSecret is the HMAC signing secret for computing signatures.
	SigningSecret string
}

// SendResponse is returned by Send.
type SendResponse struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Destinations int    `json:"destinations"`
}

// Send sends a signed webhook to a roost.
//
// It computes dual signatures: legacy X-Pigeon-Signature (hex) and
// Standard Webhooks webhook-signature (base64). Both header sets are
// sent for backward compatibility.
func (c *Client) Send(ctx context.Context, opts SendOptions) (SendResponse, error) {
	body, err := json.Marshal(opts.Payload)
	if err != nil {
		return SendResponse{}, fmt.Errorf("pgns: marshal payload: %w", err)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	msgID := "msg_" + generateUUID4()

	// Decode signing key: whsec_ prefix → base64-decode, otherwise raw bytes
	keyBytes := decodeSigningKey(opts.SigningSecret)

	// Legacy signature: HMAC-SHA256("{timestamp}.{body}") → hex
	legacyMac := hmac.New(sha256.New, keyBytes)
	legacyMac.Write([]byte(timestamp + "." + string(body)))
	legacySig := "sha256=" + hex.EncodeToString(legacyMac.Sum(nil))

	// Standard Webhooks signature: HMAC-SHA256("{msg_id}.{timestamp}.{body}") → base64
	stdMac := hmac.New(sha256.New, keyBytes)
	stdMac.Write([]byte(msgID + "." + timestamp + "." + string(body)))
	stdSig := "v1," + base64.StdEncoding.EncodeToString(stdMac.Sum(nil))

	path := "/r/" + url.PathEscape(opts.RoostID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return SendResponse{}, fmt.Errorf("pgns: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Pigeon-Signature", legacySig)
	req.Header.Set("X-Pigeon-Timestamp", timestamp)
	req.Header.Set("X-Pigeon-Event-Type", opts.EventType)
	req.Header.Set("webhook-id", msgID)
	req.Header.Set("webhook-timestamp", timestamp)
	req.Header.Set("webhook-signature", stdSig)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SendResponse{}, fmt.Errorf("pgns: send: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var result SendResponse
	if err := handleResponse(resp, &result); err != nil {
		return SendResponse{}, err
	}
	return result, nil
}

// decodeSigningKey decodes a signing secret from any supported format.
func decodeSigningKey(secret string) []byte {
	if strings.HasPrefix(secret, "whsec_") {
		decoded, err := base64.StdEncoding.DecodeString(secret[6:])
		if err == nil {
			return decoded
		}
	}
	if len(secret) == 64 {
		decoded, err := hex.DecodeString(secret)
		if err == nil {
			return decoded
		}
	}
	return []byte(secret)
}

// generateUUID4 creates a random UUID v4 without external dependencies.
func generateUUID4() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 2
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
