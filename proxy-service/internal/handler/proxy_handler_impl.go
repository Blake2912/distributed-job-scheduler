package handler

import (
	"io"
	"log/slog"
	"maps"
	"net/http"

	"github.com/Blake2912/distributed-job-scheduler/proxy-service/internal/resolver"
)

type proxyHandler struct {
	resolver resolver.LeaderResolver
}

func NewProxyHandler(resolver resolver.LeaderResolver) ProxyHandler {
	return &proxyHandler{
		resolver: resolver,
	}
}

func (p proxyHandler) ProxyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		leader, err := p.resolver.GetLeader(r.Context())
		if err != nil || leader == "" {
			http.Error(w, "Leader unavailable", http.StatusServiceUnavailable)
			return
		}

		target := leader + r.URL.Path

		slog.Info("Hitting: ", "target", target)

		req, err := http.NewRequestWithContext(
			r.Context(),
			r.Method,
			target,
			r.Body,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		req.Header = r.Header.Clone()
		
		// Fix proxy headers
		req.Host = ""
		req.Header.Set("X-Forwarded-For", r.RemoteAddr)
		req.Header.Set("X-Forwarded-Host", r.Host)
		req.Header.Set("X-Forwarded-Proto", "http")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "Failed to reach leader", 502)
			return
		}
		defer resp.Body.Close()

		maps.Copy(w.Header(), resp.Header)

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)

		slog.Info("Endpoint hit", "target", target, "response code", resp.Status)
	}
}
