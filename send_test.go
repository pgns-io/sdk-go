package sdk

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSend(t *testing.T) {
	var capturedHeaders http.Header
	var capturedBody []byte
	var capturedPath string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header.Clone()
		capturedBody, _ = io.ReadAll(r.Body)
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "pgn_test123", "status": "received", "destinations": 1,
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL)
	resp, err := client.Send(context.Background(), SendOptions{
		RoostID:       "rst_abc",
		EventType:     "order.created",
		Payload:       map[string]string{"order_id": "123"},
		SigningSecret: "test-secret",
	})
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if resp.ID != "pgn_test123" {
		t.Errorf("expected id pgn_test123, got %s", resp.ID)
	}
	if resp.Status != "received" {
		t.Errorf("expected status received, got %s", resp.Status)
	}
	if resp.Destinations != 1 {
		t.Errorf("expected destinations 1, got %d", resp.Destinations)
	}

	// Verify path
	if capturedPath != "/r/rst_abc" {
		t.Errorf("expected path /r/rst_abc, got %s", capturedPath)
	}

	// Verify headers
	sig := capturedHeaders.Get("X-Pigeon-Signature")
	if !strings.HasPrefix(sig, "sha256=") {
		t.Errorf("expected sha256= prefix, got %s", sig)
	}
	ts := capturedHeaders.Get("X-Pigeon-Timestamp")
	if ts == "" {
		t.Fatal("missing X-Pigeon-Timestamp header")
	}
	eventType := capturedHeaders.Get("X-Pigeon-Event-Type")
	if eventType != "order.created" {
		t.Errorf("expected event type order.created, got %s", eventType)
	}
	ct := capturedHeaders.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected content-type application/json, got %s", ct)
	}

	// Verify legacy HMAC signature
	signedPayload := ts + "." + string(capturedBody)
	mac := hmac.New(sha256.New, []byte("test-secret"))
	mac.Write([]byte(signedPayload))
	expectedSig := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	if sig != expectedSig {
		t.Errorf("legacy signature mismatch:\n  got:  %s\n  want: %s", sig, expectedSig)
	}

	// Verify Standard Webhooks headers
	msgID := capturedHeaders.Get("Webhook-Id")
	if !strings.HasPrefix(msgID, "msg_") {
		t.Errorf("expected webhook-id with msg_ prefix, got %s", msgID)
	}
	whTs := capturedHeaders.Get("Webhook-Timestamp")
	if whTs != ts {
		t.Errorf("webhook-timestamp %s != X-Pigeon-Timestamp %s", whTs, ts)
	}
	whSig := capturedHeaders.Get("Webhook-Signature")
	if !strings.HasPrefix(whSig, "v1,") {
		t.Errorf("expected webhook-signature with v1, prefix, got %s", whSig)
	}

	// Verify Standard Webhooks HMAC
	stdPayload := msgID + "." + ts + "." + string(capturedBody)
	stdMac := hmac.New(sha256.New, []byte("test-secret"))
	stdMac.Write([]byte(stdPayload))
	expectedStdSig := "v1," + base64.StdEncoding.EncodeToString(stdMac.Sum(nil))
	if whSig != expectedStdSig {
		t.Errorf("std webhooks signature mismatch:\n  got:  %s\n  want: %s", whSig, expectedStdSig)
	}
}

func TestSendURLEncoding(t *testing.T) {
	var capturedPath string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.RawPath
		if capturedPath == "" {
			capturedPath = r.URL.Path
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "pgn_1", "status": "received", "destinations": 0,
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL)
	_, err := client.Send(context.Background(), SendOptions{
		RoostID:       "rst/special",
		EventType:     "test",
		Payload:       map[string]any{},
		SigningSecret: "secret",
	})
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	if capturedPath != "/r/rst%2Fspecial" {
		t.Errorf("expected encoded path /r/rst%%2Fspecial, got %s", capturedPath)
	}
}
