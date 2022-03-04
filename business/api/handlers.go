package api

import "net/http"

type HttpMappings interface {
	GetMappings() map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error
}
