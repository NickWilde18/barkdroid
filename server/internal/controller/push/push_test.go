package push

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"barkdroid/internal/model"
	"barkdroid/internal/provider"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// --- mock store ---

type mockStore struct {
	devices map[string]*model.Device
}

func newMockStore() *mockStore {
	return &mockStore{devices: make(map[string]*model.Device)}
}

func (m *mockStore) RegisterDevice(platform, pushProvider, registrationID string) (*model.Device, error) {
	for _, d := range m.devices {
		if d.RegistrationID == registrationID {
			return d, nil
		}
	}
	dev := &model.Device{
		ID:             "id-" + registrationID,
		Key:            "test-key",
		Platform:       platform,
		PushProvider:   pushProvider,
		RegistrationID: registrationID,
	}
	m.devices[dev.Key] = dev
	return dev, nil
}

func (m *mockStore) GetDeviceByKey(key string) (*model.Device, error) {
	dev, ok := m.devices[key]
	if !ok {
		return nil, fmt.Errorf("device not found: %s", key)
	}
	return dev, nil
}

func (m *mockStore) ListDevices() ([]*model.Device, error) {
	var list []*model.Device
	for _, d := range m.devices {
		list = append(list, d)
	}
	return list, nil
}

func (m *mockStore) Close() error { return nil }

// --- mock provider ---

type mockProvider struct {
	name   string
	pushed []*provider.PushMessage
}

func newMockProvider(name string) *mockProvider {
	return &mockProvider{name: name}
}

func (p *mockProvider) Name() string { return p.name }
func (p *mockProvider) Push(msg *provider.PushMessage) error {
	p.pushed = append(p.pushed, msg)
	return nil
}

// --- test server ---

var (
	testServer  *ghttp.Server
	testBaseURL string
)

func setupTestServer(t *testing.T, providers map[string]provider.Provider, st *mockStore) {
	t.Helper()

	if testServer != nil {
		testServer.Shutdown()
		time.Sleep(100 * time.Millisecond)
	}

	ctrl := New(st, providers)

	name := fmt.Sprintf("test-%s", t.Name())
	s := g.Server(name)
	s.SetAddr(":0")

	s.BindHandler("/:key/:title/:body", ctrl.PushTitleBody)
	s.BindHandler("/:key/:body", ctrl.PushBody)
	s.BindHandler("POST:/push", ctrl.PushPost)
	s.BindHandler("POST:/register", ctrl.RegisterDevice)
	s.BindHandler("GET:/health", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{"status": "ok"})
	})

	go s.Run()
	time.Sleep(300 * time.Millisecond)

	testServer = s
	testBaseURL = "http://127.0.0.1" + s.GetListenedAddress()
	t.Logf("test server running at %s", testBaseURL)
}

func doGet(t *testing.T, path string) (int, string) {
	t.Helper()
	resp, err := http.Get(testBaseURL + path)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(raw)
}

func doPost(t *testing.T, path, body string) (int, string) {
	t.Helper()
	resp, err := http.Post(testBaseURL+path, "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(raw)
}

func parseJSON(raw string) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal([]byte(raw), &m)
	return m
}

func codeVal(v interface{}) int {
	if n, ok := v.(float64); ok {
		return int(n)
	}
	return 0
}

// --- tests ---

func TestHealthCheck(t *testing.T) {
	setupTestServer(t, nil, newMockStore())
	defer testServer.Shutdown()

	code, raw := doGet(t, "/health")
	t.Logf("health: %d %s", code, raw)

	result := parseJSON(raw)
	if result["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", result)
	}
}

func TestRegister(t *testing.T) {
	setupTestServer(t, nil, newMockStore())
	defer testServer.Shutdown()

	code, raw := doPost(t, "/register",
		`{"platform":"android","push_provider":"jpush","registration_id":"reg-test"}`)
	t.Logf("register: %d %s", code, raw)

	result := parseJSON(raw)
	if codeVal(result["code"]) != 200 {
		t.Errorf("expected code 200, got %v: %s", result["code"], result["message"])
	}
	if data, _ := result["data"].(map[string]interface{}); data == nil || data["key"] != "test-key" {
		t.Errorf("expected key=test-key, got %v", result["data"])
	}
}

func TestPushBody(t *testing.T) {
	mp := newMockProvider("jpush")
	setupTestServer(t, map[string]provider.Provider{"jpush": mp}, newMockStore())
	defer testServer.Shutdown()

	doPost(t, "/register", `{"platform":"android","push_provider":"jpush","registration_id":"reg-push"}`)

	code, raw := doGet(t, "/test-key/HelloWorld")
	t.Logf("push: %d %s", code, raw)

	result := parseJSON(raw)
	if codeVal(result["code"]) != 200 {
		t.Errorf("expected code 200, got %v: %s", result["code"], raw)
	}

	if len(mp.pushed) != 1 {
		t.Fatalf("expected 1 push, got %d", len(mp.pushed))
	}
	if mp.pushed[0].Body != "HelloWorld" {
		t.Errorf("expected 'HelloWorld', got %q", mp.pushed[0].Body)
	}
}

func TestPushWithTitle(t *testing.T) {
	mp := newMockProvider("jpush")
	setupTestServer(t, map[string]provider.Provider{"jpush": mp}, newMockStore())
	defer testServer.Shutdown()

	doPost(t, "/register", `{"platform":"android","push_provider":"jpush","registration_id":"reg-title"}`)
	doGet(t, "/test-key/ServerAlert/CPUHigh")

	if len(mp.pushed) != 1 {
		t.Fatalf("expected 1 push, got %d", len(mp.pushed))
	}
	if mp.pushed[0].Title != "ServerAlert" {
		t.Errorf("expected 'ServerAlert', got %q", mp.pushed[0].Title)
	}
}

func TestPushWithURL(t *testing.T) {
	mp := newMockProvider("jpush")
	setupTestServer(t, map[string]provider.Provider{"jpush": mp}, newMockStore())
	defer testServer.Shutdown()

	doPost(t, "/register", `{"platform":"android","push_provider":"jpush","registration_id":"reg-url"}`)
	doGet(t, "/test-key/Alert/Body?url=https://example.com")

	if len(mp.pushed) != 1 {
		t.Fatalf("expected 1 push, got %d", len(mp.pushed))
	}
	if mp.pushed[0].URL != "https://example.com" {
		t.Errorf("expected URL, got %q", mp.pushed[0].URL)
	}
}

func TestPostPush(t *testing.T) {
	mp := newMockProvider("jpush")
	setupTestServer(t, map[string]provider.Provider{"jpush": mp}, newMockStore())
	defer testServer.Shutdown()

	doPost(t, "/register", `{"platform":"android","push_provider":"jpush","registration_id":"reg-post"}`)
	doPost(t, "/push", `{"key":"test-key","title":"Test","body":"Hello POST","url":"https://grafana.example.com"}`)

	if mp.pushed[0].URL != "https://grafana.example.com" {
		t.Errorf("expected URL, got %q", mp.pushed[0].URL)
	}
}

func TestUnknownKey(t *testing.T) {
	setupTestServer(t, nil, newMockStore())
	defer testServer.Shutdown()

	_, raw := doGet(t, "/unknown/body")
	result := parseJSON(raw)
	if codeVal(result["code"]) != 404 {
		t.Errorf("expected 404, got %v: %s", result["code"], raw)
	}
}

func TestMissingProvider(t *testing.T) {
	st := newMockStore()
	st.RegisterDevice("android", "unknown", "reg-x")
	setupTestServer(t, nil, st)
	defer testServer.Shutdown()

	_, raw := doGet(t, "/test-key/body")
	result := parseJSON(raw)
	if codeVal(result["code"]) != 500 {
		t.Errorf("expected 500, got %v: %s", result["code"], raw)
	}
}

func TestBarkForward(t *testing.T) {
	mp := newMockProvider("bark")
	st := newMockStore()
	st.RegisterDevice("ios", "bark", "ios-token")
	setupTestServer(t, map[string]provider.Provider{"bark": mp}, st)
	defer testServer.Shutdown()

	doGet(t, "/test-key/Test/HelloIOS")

	if len(mp.pushed) != 1 {
		t.Fatalf("expected 1 push, got %d", len(mp.pushed))
	}
}
