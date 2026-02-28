package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func writeJSON(w http.ResponseWriter, v any) {
	_ = json.NewEncoder(w).Encode(v)
}

func readJSON(r *http.Request, v any) {
	_ = json.NewDecoder(r.Body).Decode(v)
}

func setupServer(t *testing.T) (*httptest.Server, *Client) {
	t.Helper()
	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc("POST /v1/auth/signup", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		writeJSON(w, AuthTokens{AccessToken: "new_token", TokenType: "Bearer", ExpiresIn: 3600})
	})
	mux.HandleFunc("POST /v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, AuthTokens{AccessToken: "jwt_token", TokenType: "Bearer", ExpiresIn: 3600})
	})
	mux.HandleFunc("POST /v1/auth/magic-link", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, MagicLinkResponse{Message: "check your email"})
	})
	mux.HandleFunc("POST /v1/auth/magic-link/verify", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, AuthTokens{AccessToken: "magic_token", TokenType: "Bearer", ExpiresIn: 3600})
	})
	mux.HandleFunc("POST /v1/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, AuthTokens{AccessToken: "refreshed", TokenType: "Bearer", ExpiresIn: 3600})
	})
	mux.HandleFunc("POST /v1/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Roosts
	mux.HandleFunc("GET /v1/roosts", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, []Roost{{
			ID: "r1", Name: "Test Roost", Description: "desc",
			IsActive: true, CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z",
		}})
	})
	mux.HandleFunc("GET /v1/roosts/{id}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, Roost{
			ID: r.PathValue("id"), Name: "Test Roost", Description: "desc",
			IsActive: true, CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z",
		})
	})
	mux.HandleFunc("POST /v1/roosts", func(w http.ResponseWriter, r *http.Request) {
		var req CreateRoost
		readJSON(r, &req)
		w.WriteHeader(http.StatusCreated)
		writeJSON(w, Roost{
			ID: "r_new", Name: req.Name, Description: "",
			IsActive: true, CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z",
		})
	})
	mux.HandleFunc("PATCH /v1/roosts/{id}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, Roost{
			ID: r.PathValue("id"), Name: "Updated", Description: "desc",
			IsActive: true, CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z",
		})
	})
	mux.HandleFunc("DELETE /v1/roosts/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Pigeons
	mux.HandleFunc("GET /v1/pigeons", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, PaginatedPigeons{
			Data: []Pigeon{{
				ID: "p1", RoostID: "r1", SourceIP: "127.0.0.1",
				RequestMethod: "POST", ContentType: "application/json",
				Headers: map[string]any{}, DeliveryStatus: "delivered",
				ReceivedAt: "2024-01-01T00:00:00Z",
			}},
			HasMore: false,
		})
	})
	mux.HandleFunc("GET /v1/pigeons/{id}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, Pigeon{
			ID: r.PathValue("id"), RoostID: "r1", SourceIP: "127.0.0.1",
			RequestMethod: "POST", ContentType: "application/json",
			Headers: map[string]any{}, DeliveryStatus: "delivered",
			ReceivedAt: "2024-01-01T00:00:00Z",
		})
	})
	mux.HandleFunc("GET /v1/pigeons/{id}/deliveries", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, PaginatedDeliveryAttempts{Data: []DeliveryAttempt{}, HasMore: false})
	})
	mux.HandleFunc("POST /v1/pigeons/{id}/replay", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, ReplayResponse{Replayed: true, PigeonID: r.PathValue("id"), DeliveryAttempts: 2})
	})

	// Destinations
	mux.HandleFunc("GET /v1/roosts/{id}/destinations", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, []Destination{{
			ID: "d1", RoostID: r.PathValue("id"), DestinationType: "url",
			Config: map[string]any{"url": "https://example.com"}, RetryMax: 5,
			RetryDelayMs: 1000, RetryMultiplier: 2.0,
			CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z",
		}})
	})
	mux.HandleFunc("GET /v1/destinations/{id}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, Destination{
			ID: r.PathValue("id"), RoostID: "r1", DestinationType: "url",
			Config: map[string]any{}, RetryMax: 5, RetryDelayMs: 1000, RetryMultiplier: 2.0,
			CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z",
		})
	})
	mux.HandleFunc("POST /v1/roosts/{id}/destinations", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		writeJSON(w, Destination{
			ID: "d_new", RoostID: r.PathValue("id"), DestinationType: "url",
			Config: map[string]any{}, RetryMax: 5, RetryDelayMs: 1000, RetryMultiplier: 2.0,
			CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z",
		})
	})
	mux.HandleFunc("PATCH /v1/destinations/{id}/pause", func(w http.ResponseWriter, r *http.Request) {
		var input PauseInput
		readJSON(r, &input)
		writeJSON(w, PauseResponse(input))
	})
	mux.HandleFunc("DELETE /v1/destinations/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// API Keys
	mux.HandleFunc("GET /v1/api-keys", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, []ApiKeyResponse{{ID: "k1", KeyPrefix: "pk_live_test1234", Name: "Default", CreatedAt: "2024-01-01T00:00:00Z"}})
	})
	mux.HandleFunc("GET /v1/api-keys/{id}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, ApiKeyResponse{ID: r.PathValue("id"), KeyPrefix: "pk_live_test1234", Name: "Default", CreatedAt: "2024-01-01T00:00:00Z"})
	})
	mux.HandleFunc("POST /v1/api-keys", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		writeJSON(w, ApiKeyCreatedResponse{ID: "k_new", Key: "pk_live_fullkey", KeyPrefix: "pk_live_fullk", Name: "Default", CreatedAt: "2024-01-01T00:00:00Z"})
	})
	mux.HandleFunc("PATCH /v1/api-keys/{id}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, ApiKeyResponse{ID: r.PathValue("id"), KeyPrefix: "pk_live_test1234", Name: "Renamed", CreatedAt: "2024-01-01T00:00:00Z"})
	})
	mux.HandleFunc("DELETE /v1/api-keys/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Templates
	mux.HandleFunc("GET /v1/templates", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, []Template{{ID: "t1", Name: "Test", Description: "desc", Body: "{{ body }}", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"}})
	})
	mux.HandleFunc("GET /v1/templates/{id}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, Template{ID: r.PathValue("id"), Name: "Test", Description: "desc", Body: "{{ body }}", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"})
	})
	mux.HandleFunc("POST /v1/templates/preview", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, PreviewTemplateResponse{Rendered: "hello world"})
	})
	mux.HandleFunc("POST /v1/templates", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		writeJSON(w, Template{ID: "t_new", Name: "New", Description: "", Body: "", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"})
	})
	mux.HandleFunc("PATCH /v1/templates/{id}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, Template{ID: r.PathValue("id"), Name: "Updated", Description: "desc", Body: "{{ body }}", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"})
	})
	mux.HandleFunc("DELETE /v1/templates/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Billing
	mux.HandleFunc("POST /v1/billing/checkout", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, CheckoutResponse{CheckoutURL: "https://checkout.stripe.com/session"})
	})
	mux.HandleFunc("POST /v1/billing/portal", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, PortalResponse{PortalURL: "https://billing.stripe.com/portal"})
	})
	mux.HandleFunc("GET /v1/billing/status", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, BillingStatus{
			Plan: "pro", SubscriptionStatus: "active", UsageCount: 42,
			Limits: BillingLimits{PigeonsPerMonth: 10000, MaxRoosts: 50, APIPerMinute: 1000, InboundPerSecond: 100, EmailPerMonth: 500},
		})
	})

	// User
	mux.HandleFunc("GET /v1/me", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, User{ID: "u1", Email: "test@example.com", Name: "Test User", Plan: "free", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"})
	})
	mux.HandleFunc("PATCH /v1/me", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, User{ID: "u1", Email: "test@example.com", Name: "Updated", Plan: "free", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"})
	})
	mux.HandleFunc("GET /v1/stats", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, DashboardStats{TotalRoosts: 5, ActiveRoosts: 3, TotalPigeons: 100, PigeonsToday: 10, Delivered: 90, Failed: 5})
	})

	server := httptest.NewServer(mux)
	client := NewClient(server.URL, WithAPIKey("pk_live_test"))
	return server, client
}

func TestAuthHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer pk_live_test" {
			t.Errorf("expected Bearer pk_live_test, got %s", auth)
		}
		writeJSON(w, []Roost{})
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAPIKey("pk_live_test"))
	_, err := client.ListRoosts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestSignup(t *testing.T) {
	server, _ := setupServer(t)
	defer server.Close()

	unauthClient := NewClient(server.URL)
	tokens, err := unauthClient.Signup(context.Background(), SignupRequest{
		Email: "test@example.com", Password: "secret", TOSAccepted: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if tokens.AccessToken != "new_token" {
		t.Errorf("expected new_token, got %s", tokens.AccessToken)
	}
}

func TestLogin(t *testing.T) {
	server, _ := setupServer(t)
	defer server.Close()

	client := NewClient(server.URL)
	tokens, err := client.Login(context.Background(), LoginRequest{
		Email: "test@example.com", Password: "secret",
	})
	if err != nil {
		t.Fatal(err)
	}
	if tokens.AccessToken != "jwt_token" {
		t.Errorf("expected jwt_token, got %s", tokens.AccessToken)
	}
}

func TestLogout(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	if err := client.Logout(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestListRoosts(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	roosts, err := client.ListRoosts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(roosts) != 1 {
		t.Fatalf("expected 1 roost, got %d", len(roosts))
	}
	if roosts[0].Name != "Test Roost" {
		t.Errorf("expected Test Roost, got %s", roosts[0].Name)
	}
}

func TestGetRoost(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	roost, err := client.GetRoost(context.Background(), "r1")
	if err != nil {
		t.Fatal(err)
	}
	if roost.ID != "r1" {
		t.Errorf("expected r1, got %s", roost.ID)
	}
}

func TestCreateRoost(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	roost, err := client.CreateRoost(context.Background(), CreateRoost{Name: "New Roost"})
	if err != nil {
		t.Fatal(err)
	}
	if roost.Name != "New Roost" {
		t.Errorf("expected New Roost, got %s", roost.Name)
	}
}

func TestDeleteRoost(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	if err := client.DeleteRoost(context.Background(), "r1"); err != nil {
		t.Fatal(err)
	}
}

func TestListPigeons(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	result, err := client.ListPigeons(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 pigeon, got %d", len(result.Data))
	}
}

func TestReplayPigeon(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	result, err := client.ReplayPigeon(context.Background(), "p1")
	if err != nil {
		t.Fatal(err)
	}
	if !result.Replayed {
		t.Error("expected replayed to be true")
	}
	if result.DeliveryAttempts != 2 {
		t.Errorf("expected 2 delivery attempts, got %d", result.DeliveryAttempts)
	}
}

func TestListDestinations(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	dests, err := client.ListDestinations(context.Background(), "r1")
	if err != nil {
		t.Fatal(err)
	}
	if len(dests) != 1 {
		t.Fatalf("expected 1 destination, got %d", len(dests))
	}
	if dests[0].DestinationType != "url" {
		t.Errorf("expected url, got %s", dests[0].DestinationType)
	}
}

func TestPauseDestination(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	resp, err := client.PauseDestination(context.Background(), "d1", true)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsPaused {
		t.Error("expected is_paused to be true")
	}
}

func TestDeleteDestination(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	if err := client.DeleteDestination(context.Background(), "d1"); err != nil {
		t.Fatal(err)
	}
}

func TestListApiKeys(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	keys, err := client.ListApiKeys(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}
}

func TestCreateApiKey(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	result, err := client.CreateApiKey(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(result.Key, "pk_live_") {
		t.Errorf("expected key to start with pk_live_, got %s", result.Key)
	}
}

func TestListTemplates(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	templates, err := client.ListTemplates(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) != 1 {
		t.Fatalf("expected 1 template, got %d", len(templates))
	}
}

func TestPreviewTemplate(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	resp, err := client.PreviewTemplate(context.Background(), PreviewTemplateRequest{
		Body: "{{ body }}", PigeonID: "p1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Rendered != "hello world" {
		t.Errorf("expected 'hello world', got %s", resp.Rendered)
	}
}

func TestBillingStatus(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	status, err := client.GetBillingStatus(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if status.Plan != "pro" {
		t.Errorf("expected pro, got %s", status.Plan)
	}
	if status.Limits.PigeonsPerMonth != 10000 {
		t.Errorf("expected 10000, got %d", status.Limits.PigeonsPerMonth)
	}
}

func TestGetMe(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	user, err := client.GetMe(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected test@example.com, got %s", user.Email)
	}
}

func TestGetStats(t *testing.T) {
	server, client := setupServer(t)
	defer server.Close()

	stats, err := client.GetStats(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if stats.TotalRoosts != 5 {
		t.Errorf("expected 5, got %d", stats.TotalRoosts)
	}
	if stats.Delivered != 90 {
		t.Errorf("expected 90, got %d", stats.Delivered)
	}
}

func TestErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		writeJSON(w, map[string]string{"error": "not found"})
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAPIKey("pk_live_test"))
	_, err := client.GetRoost(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found, got %v", err)
	}
}

func TestNonJSONError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = fmt.Fprint(w, "Bad Gateway")
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAPIKey("pk_live_test"))
	_, err := client.ListRoosts(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var pe *PigeonsError
	if !errors.As(err, &pe) {
		t.Fatalf("expected PigeonsError, got %T", err)
	}
	if pe.StatusCode != 502 {
		t.Errorf("expected 502, got %d", pe.StatusCode)
	}
}

func Test401RefreshRetry(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if r.URL.Path == "/v1/auth/refresh" {
			writeJSON(w, AuthTokens{AccessToken: "refreshed", TokenType: "Bearer", ExpiresIn: 3600})
			return
		}
		if callCount == 1 {
			w.WriteHeader(http.StatusUnauthorized)
			writeJSON(w, map[string]string{"error": "token expired"})
			return
		}
		writeJSON(w, []Roost{{
			ID: "r1", Name: "Test", IsActive: true,
			CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z",
		}})
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAccessToken("expired"))
	roosts, err := client.ListRoosts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(roosts) != 1 {
		t.Fatalf("expected 1 roost, got %d", len(roosts))
	}
}

func TestSSEEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/events" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "text/event-stream" {
			t.Errorf("expected Accept: text/event-stream")
		}
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("expected http.Flusher")
		}
		_, _ = fmt.Fprintln(w, "data: {\"id\":\"p1\",\"roost_id\":\"r1\"}")
		flusher.Flush()
		_, _ = fmt.Fprintln(w, "data: {\"id\":\"p2\",\"roost_id\":\"r1\"}")
		flusher.Flush()
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAPIKey("pk_live_test"))
	ctx, cancel := context.WithCancel(context.Background())

	var events []string
	go func() {
		_ = client.ListenEvents(ctx, func(data string) {
			events = append(events, data)
			if len(events) >= 2 {
				cancel()
			}
		})
	}()

	<-ctx.Done()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if !strings.Contains(events[0], "p1") {
		t.Errorf("expected event to contain p1, got %s", events[0])
	}
}
