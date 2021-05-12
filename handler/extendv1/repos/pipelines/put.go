package pipelines

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/google/uuid"

	"github.com/oars-sigs/drone/model"

	"github.com/drone/drone/core"
	"github.com/drone/drone/handler/api/render"
	"github.com/sirupsen/logrus"
)

// HandlePutPipeline returns an http.HandlerFunc that processes http
// requests to put a pipeline for the specified slug.
func HandlePutPipeline(
	repos core.RepositoryStore,
	pipelineStore model.PipelineStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx       = r.Context()
			name      = chi.URLParam(r, "name")
			namespace = chi.URLParam(r, "owner")
			branch    = r.FormValue("branch")
		)

		in, err := ioutil.ReadAll(r.Body)
		if err != nil {
			render.BadRequest(w, err)
			return
		}
		configPath := ".drone.yml"
		ref := "refs/heads/" + branch
		if branch == "" {
			ref = "default"
		}
		slug := namespace + "/" + name
		pipe, isExist, err := pipelineStore.GetPipeline(ctx, slug, ref, configPath)
		if err != nil {
			logrus.Error(err)
			render.InternalError(w, err)
			return
		}
		if !isExist {
			pipe := &model.Pipeline{
				UUID:       uuid.New().String(),
				Slug:       slug,
				Ref:        ref,
				ConfigPath: configPath,
				Content:    string(in),
				Created:    time.Now().Unix(),
				Updated:    time.Now().Unix(),
			}
			err = pipelineStore.CreatePipeline(ctx, pipe)
			if err != nil {
				logrus.Error(err)
				render.InternalError(w, err)
				return
			}
			w.WriteHeader(204)
		}
		pipe.Content = string(in)
		pipe.Updated = time.Now().Unix()
		err = pipelineStore.UpdatePipeline(ctx, pipe)
		if err != nil {
			logrus.Error("update", err)
			render.InternalError(w, err)
			return
		}
		w.WriteHeader(204)
	}
}
