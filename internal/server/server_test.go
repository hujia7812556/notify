package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"notify/internal/config"
	"notify/internal/dispatcher"
	"notify/internal/sender"
)

func TestHandleHealth(t *testing.T) {
	// 创建测试服务器
	cfg := config.ServerConfig{
		Port: 8081,
		Mode: "debug",
	}
	senderMgr := sender.NewManager()
	disp := dispatcher.New(100, 10, senderMgr)
	srv := New(cfg, disp)
	srv.registerRoutes()

	// 创建测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	srv.engine.ServeHTTP(w, req)

	// 检查响应
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}
}

func TestHandleNotify(t *testing.T) {
	// 创建测试服务器
	cfg := config.ServerConfig{
		Port: 8081,
		Mode: "debug",
	}
	senderMgr := sender.NewManager()
	disp := dispatcher.New(100, 10, senderMgr)
	srv := New(cfg, disp)
	srv.registerRoutes()

	// 测试用例
	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid wxpusher message",
			payload: map[string]interface{}{
				"platform": "wechat",
				"content":  "test message",
				"extra": map[string]interface{}{
					"user_id": "UID_xxx",
				},
			},
			wantStatus: http.StatusAccepted,
		},
		{
			name: "invalid platform",
			payload: map[string]interface{}{
				"platform": "invalid",
				"content":  "test message",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing content",
			payload: map[string]interface{}{
				"platform": "wechat",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/notify", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			srv.engine.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
