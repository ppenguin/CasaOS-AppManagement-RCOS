package route

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ReCasaOS/CasaOS-AppManagement/pkg/config"
	"github.com/ReCasaOS/CasaOS-Common/external"
)

// an image-managed host: every docker-location route 404s, none reaches a container by id
func TestDockerLocationImageManaged(t *testing.T) {
	old := config.ServerInfo.ImageManaged
	config.ServerInfo.ImageManaged = true
	t.Cleanup(func() { config.ServerInfo.ImageManaged = old })
	runtime := t.TempDir()
	previous := config.CommonInfo.RuntimePath
	config.CommonInfo.RuntimePath = runtime
	t.Cleanup(func() { config.CommonInfo.RuntimePath = previous })
	if err := os.WriteFile(filepath.Join(runtime, external.InternalSecretFilename), []byte("s3cret"), 0o600); err != nil {
		t.Fatal(err)
	}

	e := InitV1Router()
	for _, r := range []struct{ method, path string }{
		{http.MethodGet, "/v1/container/docker-location"},
		{http.MethodGet, "/v1/container/docker-location/check?default=true"},
		{http.MethodPost, "/v1/container/docker-location"},
		{http.MethodDelete, "/v1/container/docker-location/old"},
	} {
		req := httptest.NewRequest(r.method, r.path, nil)
		req.RemoteAddr = "127.0.0.1:40000"
		req.Header.Set("Authorization", "Internal s3cret")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s %s: %d, want 404", r.method, r.path, rec.Code)
		}
	}
}
