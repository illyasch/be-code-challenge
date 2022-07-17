// Package handlers manages the API.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/illyasch/be-code-challenge/pkg/business/calc"
	"github.com/illyasch/be-code-challenge/pkg/data/database"
)

// APIConfig contains all the mandatory systems required by handlers.
type APIConfig struct {
	Log *zap.SugaredLogger
	DB  *sqlx.DB
}

type errorResponse struct {
	Error string `json:"error"`
}

// Router constructs a http.Handler with all application routes defined.
func (cfg APIConfig) Router() http.Handler {
	store := calc.New(cfg.DB)

	router := mux.NewRouter()
	router.HandleFunc("/hourly", cfg.handleHourly(store)).Methods(http.MethodGet)
	router.HandleFunc("/readiness", cfg.handleReadiness).Methods(http.MethodGet)
	router.HandleFunc("/liveness", cfg.handleLiveness).Methods(http.MethodGet)

	return router
}

// handleHourly handler returns the amount of fees being paid for Ethereum transactions per hour.
func (cfg APIConfig) handleHourly(store calc.Calc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		hours, err := store.Hourly(r.Context())
		if err == nil {
			cfg.respond(w, http.StatusOK, hours)
			cfg.Log.Infow("hourly", "statusCode", http.StatusOK, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
			return
		}

		status := http.StatusInternalServerError
		cfg.respond(w, status, errorResponse{
			Error: http.StatusText(status),
		})
		cfg.Log.Errorw("hourly", "ERROR", fmt.Errorf("calc: %w", err))
		return
	}
}

// handleReadiness checks if the database is ready and if not will return a 500 status if it's not.
func (cfg APIConfig) handleReadiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	status := "ok"
	statusCode := http.StatusOK
	if err := database.StatusCheck(ctx, cfg.DB); err != nil {
		status = "db not ready"
		statusCode = http.StatusInternalServerError
		cfg.Log.Errorw("readiness", "ERROR", fmt.Errorf("status check: %w", err))
	}

	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	cfg.respond(w, statusCode, data)
	cfg.Log.Infow("readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
}

// handleLiveness returns simple status info if the service is alive. If the
// app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
func (cfg APIConfig) handleLiveness(w http.ResponseWriter, r *http.Request) {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := struct {
		Status    string `json:"status,omitempty"`
		Build     string `json:"build,omitempty"`
		Host      string `json:"host,omitempty"`
		Pod       string `json:"pod,omitempty"`
		PodIP     string `json:"podIP,omitempty"`
		Node      string `json:"node,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	}{
		Status:    "up",
		Host:      host,
		Pod:       os.Getenv("KUBERNETES_PODNAME"),
		PodIP:     os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
		Node:      os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
	}

	statusCode := http.StatusOK
	cfg.respond(w, statusCode, data)
	cfg.Log.Infow("liveness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
}

func (cfg APIConfig) respond(w http.ResponseWriter, statusCode int, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		cfg.Log.Errorw("respond", "ERROR", fmt.Errorf("json marshal: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		cfg.Log.Errorw("respond", "ERROR", fmt.Errorf("write output: %w", err))
		return
	}
}
