package scene

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/scene"
	scene2 "com.fha.gocan/business/data/store/scene"
	"com.fha.gocan/foundation/web"
	"net/http"
)

type Handlers struct {
	App   app.Core
	Scene scene.Core
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

type ContextRouteData interface {
	Route() string
	Params() map[string]string
}

type applicationDto struct {
	Id string `json:"id"`
}

func (h *Handlers) QueryById(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	id := params["id"]
	s, err := h.Scene.QueryById(id)

	if err != nil {
		return err
	}

	apps, _ := h.App.QueryBySceneName(s.Name)
	if err != nil {
		return err
	}

	appDtos := []applicationDto{}
	for _, a := range apps {
		appDtos = append(appDtos, applicationDto{
			Id: a.Id,
		})
	}
	data := struct {
		Id           string           `json:"id"`
		Name         string           `json:"name"`
		Applications []applicationDto `json:"applications"`
	}{
		Id:           s.Id,
		Name:         s.Name,
		Applications: appDtos,
	}

	return web.Respond(w, data, 200)
}
