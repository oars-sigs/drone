package pipelines

import (
	"net/http"

	"github.com/oars-sigs/drone/model"
	"github.com/sirupsen/logrus"

	"github.com/drone/drone/core"
	"github.com/drone/drone/handler/api/render"
	"github.com/go-chi/chi"
)

// HandleFindPipelines returns an http.HandlerFunc that processes http
// requests to a project all pipelines for the specified slug.
func HandleFindPipelines(
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
		ref := "refs/heads/" + branch
		if branch == "" {
			ref = "default"
		}
		configPath := ".drone.yml"
		slug := namespace + "/" + name
		pipe, isExist, err := pipelineStore.GetPipeline(ctx, slug, ref, configPath)
		if err != nil {
			logrus.Error(err)
			render.InternalError(w, err)
			return
		}
		if !isExist {
			render.JSON(w, map[string]string{"data": ""}, 200)
			return
		}
		render.JSON(w, map[string]string{"data": pipe.Content}, 200)
		return
	}
}
