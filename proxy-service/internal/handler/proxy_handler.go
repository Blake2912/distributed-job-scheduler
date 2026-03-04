package handler

import "net/http"

type ProxyHandler interface {
	ProxyHandler() http.HandlerFunc
}
