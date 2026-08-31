package router

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"mewsoproxy/server/config"
	"mewsoproxy/server/webui"
)

func TestServeEmbeddedSPA(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := New(config.Load(), nil, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET / status = %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "/assets/") {
		t.Fatalf("GET / body is not SPA index.html: %s", w.Body.String())
	}
}

func TestHealthz(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := New(config.Load(), nil, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET /healthz status = %d", w.Code)
	}
}

func TestServeEmbeddedAsset(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := New(config.Load(), nil, nil)

	assets, err := fs.Sub(webui.Dist, "dist/assets")
	if err != nil {
		t.Fatalf("sub assets: %v", err)
	}
	entries, _ := fs.ReadDir(assets, ".")
	var name string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".js") {
			name = e.Name()
			break
		}
	}
	if name == "" {
		t.Fatal("no js asset found in embedded dist")
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/assets/"+name, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET /assets/%s status = %d", name, w.Code)
	}
}
