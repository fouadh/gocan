package scene

import (
	scene2 "com.fha.gocan/business/data/store/scene"
	"com.fha.gocan/foundation/web"
	"net/http"
)

type Handlers struct {
	Scene Core
}

func (h *Handlers) QueryAll(w http.ResponseWriter, r *http.Request) error {
	scenes, err := h.Scene.QueryAll()
	if err != nil {
		return err
	}
	list := struct {
		Scenes []scene2.Scene `json:"scenes"`
	}{
		Scenes: scenes,
	}
	return web.Respond(w, list, 200)
}