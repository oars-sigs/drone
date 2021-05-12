package templates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/drone/drone/handler/api/render"
	"github.com/go-chi/chi"
	"github.com/oars-sigs/drone/model"
	"github.com/sirupsen/logrus"
)

// HandleGetTemp returns an http.HandlerFunc that processes http
// requests to get the templates.
func HandleGetTemp(
	tmpls model.TemplateStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templates, err := tmpls.GetTemplate()

		if err != nil {
			render.InternalError(w, err)
		}

		ts := make([]model.Template, 0)
		for _, r := range templates {
			t := model.Template{
				UUID:    r.UUID,
				Name:    r.Name,
				Format:  r.Format,
				Type:    r.Type,
				Content: r.Content,
			}
			ts = append(ts, t)
		}
		res := model.TemplateJSON{
			Total:     len(ts),
			Templates: ts,
		}
		logrus.Debug(res)
		render.JSON(w, res, 200)
	}
}

// HandleCreateTemp returns an http.HandlerFunc that processes http
// requests to create the templates.
func HandleCreateTemp(
	tmpls model.TemplateStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx = r.Context()
		)

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		var tmps model.Template
		err = json.Unmarshal(data, &tmps)
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		err = tmpls.CreateTemplate(ctx, &tmps)

		if err != nil {
			render.InternalError(w, err)
			return
		}
		render.JSON(w, nil, 200)
	}
}

// HandlePutTemp returns an http.HandlerFunc that processes http
// requests to modify the templates.
func HandlePutTemp(
	tmpls model.TemplateStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx = r.Context()
		)

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		var tmps model.Template
		err = json.Unmarshal(data, &tmps)
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		err = tmpls.PutTemplate(ctx, &tmps)

		if err != nil {
			render.InternalError(w, err)
			return
		}
		render.JSON(w, nil, 200)
	}
}

// HandleDeleteTemp returns an http.HandlerFunc that processes http
// requests to delete the templates.
func HandleDeleteTemp(
	tmpls model.TemplateStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx  = r.Context()
			uuid = chi.URLParam(r, "uuid")
		)

		err := tmpls.DeleteTemplate(ctx, uuid)

		if err != nil {
			render.InternalError(w, err)
		}
		render.JSON(w, nil, 200)
	}
}

// HandleFindTemp returns an http.HandlerFunc that processes http
// requests to find a templates.
func HandleFindTemp(
	tmpls model.TemplateStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx  = r.Context()
			uuid = chi.URLParam(r, "uuid")
		)

		out, _, err := tmpls.FindTemplate(ctx, uuid)
		res := model.Template{
			UUID:    out.UUID,
			Name:    out.Name,
			Format:  out.Format,
			Type:    out.Type,
			Content: out.Content,
		}

		if err != nil {
			render.InternalError(w, err)
		}
		render.JSON(w, res, 200)
	}
}
