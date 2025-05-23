package yandexcloud

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/context_key"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func TestGetHTTPClient_Caching(t *testing.T) {
	c1 := GetHTTPClient(10)
	c2 := GetHTTPClient(10)
	if c1 != c2 {
		t.Error("GetHTTPClient should cache clients by timeout")
	}
	c3 := GetHTTPClient(20)
	if c1 == c3 {
		t.Error("GetHTTPClient should return different clients for different timeouts")
	}
}

func TestHandleHTTPError(t *testing.T) {
	r := httptest.NewRecorder()
	r.WriteHeader(404)
	r.Write([]byte("not found"))
	err := HandleHTTPError(r.Result(), 404)
	if err != nil {
		t.Errorf("Expected nil for ignored code, got %v", err)
	}

	r = httptest.NewRecorder()
	r.WriteHeader(500)
	r.Write([]byte("fail"))
	err = HandleHTTPError(r.Result())
	if err == nil {
		t.Error("Expected error for 500, got nil")
	}
}

func TestDoWithRetry(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 3 {
			w.WriteHeader(500)
			w.Write([]byte("fail"))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	var resp *http.Response
	var err error
	reqFactory := func() *http.Request {
		req, _ := http.NewRequest("GET", ts.URL, nil)
		return req
	}
	// Inject logger into context using context_key.Logger
	logger := hclog.NewNullLogger()
	ctx := context.WithValue(context.Background(), context_key.Logger, logger)
	resp, err = DoWithRetry(ctx, ts.Client(), reqFactory, 5, 2)
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}
	b, _ := io.ReadAll(resp.Body)
	if string(b) != "ok" {
		t.Errorf("Expected 'ok', got '%s'", string(b))
	}
	if calls != 3 {
		t.Errorf("Expected 3 attempts, got %d", calls)
	}
}

func TestApplyRequestOptions(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/path", nil)
	ua := "test-agent"
	epo := "https://override.com"
	ApplyRequestOptions(req, &ua, &epo)
	if req.Header.Get("User-Agent") != ua {
		t.Errorf("User-Agent not set")
	}
	if req.URL.Host != "override.com" || req.URL.Scheme != "https" {
		t.Errorf("Endpoint override failed: %v", req.URL)
	}
}

func TestBillingAccountNameUpperTransform(t *testing.T) {
	acc := &BillingAccount{Id: "id1", Name: "myAccount"}
	val, err := billingAccountNameUpperTransform(context.Background(), &transform.TransformData{HydrateItem: acc})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "MYACCOUNT" {
		t.Errorf("expected MYACCOUNT, got %v", val)
	}
}

func TestGetYandexBillingAccount_NoToken(t *testing.T) {
	d := &plugin.QueryData{Connection: &plugin.Connection{Config: &Config{}}}
	logger := hclog.NewNullLogger()
	ctx := context.WithValue(context.Background(), context_key.Logger, logger)
	item, err := getYandexBillingAccount(ctx, d, nil)
	if err == nil || item != nil {
		t.Errorf("expected error for missing token, got item=%v, err=%v", item, err)
	}
}
