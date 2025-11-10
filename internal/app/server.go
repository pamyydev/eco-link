package app

import (
	"compress/gzip"
	"encoding/json"

	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pamelamiranda/eco-link/internal/domain"
)

// StartServer inicia um servidor HTTP simples que atua como BFF para o frontend.
// Rotas:
// - GET /api/analyze?url={url} -> retorna domain.Report em JSON
// - GET /api/greencheck?url={url} -> retorna GreenCheck em JSON
// - static files servidos de FRONTEND_DIST (./frontend/dist por padrão)
func (a *Application) StartServer() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/analyze", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Query().Get("url")
		if urlParam == "" {
			http.Error(w, "missing url parameter", http.StatusBadRequest)
			return
		}

		report, err := a.service.Analyze(urlParam)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			if jerr := json.NewEncoder(w).Encode(report); jerr != nil {
				a.logger.Printf("error encoding report json: %v", jerr)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Log original analysis error and continue with fallback assembly
		a.logger.Printf("service.Analyze error for %s: %v; building fallback report", urlParam, err)

		// Fallback: tentar montar um relatório parcial usando fetchGreen (mock ou API)
		site, serr := domain.NewSite(urlParam)
		if serr != nil {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}

		// tentar obter greencheck (pode usar mock)
		gi, gerr := a.fetchGreen(site.URL)
		if gerr != nil {
			// se tudo falhar, usar valores padrões
			gi = GreenCheck{Carbon: 0.0, Green: false, Bytes: 0}
		}

		metrics := domain.Metrics{
			CarbonPerVisit: gi.Carbon,
			BytesPerVisit:  int(gi.Bytes),
			GreenHost:      gi.Green,
		}

		// If carbon is zero but we have bytes, estimate CO2 using same formula as WebsiteCarbonAdapter
		if metrics.CarbonPerVisit == 0 && metrics.BytesPerVisit > 0 {
			// estimation: CO2 = (pageSizeKB) * 0.0005 per KB
			pageSizeKB := float64(metrics.BytesPerVisit) / 1024.0
			metrics.CarbonPerVisit = pageSizeKB * 0.0005
		}

		fallbackReport := domain.NewReport(*site, metrics)
		// mark that this is a fallback/partial report so frontend can show a note
		fallbackReport.Warning = "partial report: website carbon metrics unavailable; values may be estimated or from mock"

		w.Header().Set("Content-Type", "application/json")
		// Ensure we explicitly return 200 with a JSON body (fallback)
		w.WriteHeader(http.StatusOK)
		if jerr := json.NewEncoder(w).Encode(fallbackReport); jerr != nil {
			a.logger.Printf("error encoding fallback report json: %v", jerr)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/api/greencheck", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Query().Get("url")
		if urlParam == "" {
			http.Error(w, "missing url parameter", http.StatusBadRequest)
			return
		}
		// extrair hostname
		parsed, err := url.Parse(urlParam)
		if err != nil {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		host := parsed.Hostname()
		gi, err := a.fetchGreen(host)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(gi)
	})

	// Servir arquivos estáticos do frontend (build)
	dist := os.Getenv("FRONTEND_DIST")
	if dist == "" {
		dist = filepath.Join("frontend", "dist")
	}
	// Se a pasta existir, servir-a na raiz com cache-control e gzip
	if fi, err := os.Stat(dist); err == nil && fi.IsDir() {
		fs := http.FileServer(http.Dir(dist))
		// Wrap to set caching headers and allow gzip
		fileHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Cache static assets for 1h (can be tuned)
			w.Header().Set("Cache-Control", "public, max-age=3600")
			fs.ServeHTTP(w, r)
		})
		mux.Handle("/", withGzip(fileHandler))
		a.logger.Printf("serving frontend static files from %s", dist)
	} else {
		a.logger.Printf("no frontend dist found at %s; API only mode", dist)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	addr := ":" + port

	// Wrap with gzip + CORS middleware allowing frontend dev to call /api
	handler := withCORS(withGzip(mux))

	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	a.logger.Printf("starting BFF server on %s", addr)
	return srv.ListenAndServe()
}

// withCORS adiciona headers CORS para rotas /api e responde OPTIONS.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/api") {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// withGzip compresses responses when the client supports gzip.
func withGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Don't compress for OPTIONS
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// We will create a gzip writer only when Write is called to avoid unnecessary work

		// custom response writer
		gz := &gzipResponseWriter{ResponseWriter: w, writer: nil}
		// ensure Content-Type is sniffed by underlying handler
		next.ServeHTTP(gz, r)
		if gz.writer != nil {
			gz.writer.Close()
		}
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	if g.writer == nil {
		// set header now that we will compress
		g.Header().Set("Content-Encoding", "gzip")
		// delete content-length because it will change
		g.Header().Del("Content-Length")
		gw := gzip.NewWriter(g.ResponseWriter)
		g.writer = gw
	}
	return g.writer.Write(b)
}
