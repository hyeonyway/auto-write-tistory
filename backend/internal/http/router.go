package http

import (
	"net/http"
	"strings"

	"devlog-studio/backend/internal/config"
	"devlog-studio/backend/internal/post"
	"devlog-studio/backend/internal/publisher/tistory"
	"devlog-studio/backend/internal/settings"
	apptemplate "devlog-studio/backend/internal/template"
	tistoryapi "devlog-studio/backend/internal/tistory"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(cfg config.Config, pool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	templateRepo := apptemplate.NewRepository(pool)
	renderer := apptemplate.NewRenderer(templateRepo)
	postRepo := post.NewRepository(pool)
	settingsRepo := settings.NewRepository(pool)
	settingsSvc := settings.NewService(settingsRepo)
	publisher := tistory.NewPublisher("")
	postSvc := post.NewService(postRepo, renderer, settingsSvc, publisher)
	tistorySvc := tistoryapi.NewService(publisher, settingsSvc)

	postHandler := post.NewHandler(postSvc, writeJSON, writeError, decodeJSON)
	settingsHandler := settings.NewHandler(settingsSvc, writeJSON, writeError, decodeJSON)
	tistoryHandler := tistoryapi.NewHandler(tistorySvc, writeJSON, writeError, decodeJSON)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("POST /api/posts/preview", postHandler.Preview)
	mux.HandleFunc("POST /api/posts", postHandler.Create)
	mux.HandleFunc("GET /api/posts", postHandler.List)
	mux.HandleFunc("GET /api/posts/{id}", postHandler.Get)
	mux.HandleFunc("PUT /api/posts/{id}", postHandler.Update)
	mux.HandleFunc("DELETE /api/posts/{id}", postHandler.Delete)
	mux.HandleFunc("POST /api/posts/{id}/publish/tistory", postHandler.PublishTistory)

	mux.HandleFunc("GET /api/settings", settingsHandler.Get)
	mux.HandleFunc("PUT /api/settings", settingsHandler.Update)
	mux.HandleFunc("POST /api/tistory/session/start", tistoryHandler.StartSession)
	mux.HandleFunc("POST /api/tistory/session/confirm", tistoryHandler.ConfirmSession)
	mux.HandleFunc("GET /api/tistory/session/status", tistoryHandler.SessionStatus)
	mux.HandleFunc("DELETE /api/tistory/session", tistoryHandler.DeleteSession)
	mux.HandleFunc("POST /api/tistory/categories/fetch", tistoryHandler.FetchCategories)

	return cors(cfg.CORSAllowedOrigins, mux)
}

func cors(allowedOrigins []string, next http.Handler) http.Handler {
	allowed := map[string]bool{}
	for _, origin := range allowedOrigins {
		allowed[origin] = true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && (allowed[origin] || len(allowed) == 0) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		}
		if strings.EqualFold(r.Method, http.MethodOptions) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
