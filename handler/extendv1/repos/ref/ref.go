package ref

import (
	"net/http"

	"github.com/drone/drone/core"
	"github.com/drone/drone/handler/api/render"
	"github.com/go-chi/chi"

	"github.com/oars-sigs/drone/model"
)

// HandleFindTags returns an http.HandlerFunc that processes http
// requests to a project all tags for the specified slug.
func HandleFindTags(
	repos core.RepositoryStore,
	gits model.GitService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx       = r.Context()
			name      = chi.URLParam(r, "name")
			namespace = chi.URLParam(r, "owner")
		)
		slug := namespace + "/" + name
		tags, err := gits.FindTags(ctx, nil, slug)
		if err != nil {
			render.InternalError(w, err)
			return
		}
		render.JSON(w, tags, 200)
	}
}

// HandleFindBranches returns an http.HandlerFunc that processes http
// requests to a project all branches for the specified slug.
func HandleFindBranches(
	repos core.RepositoryStore,
	gits model.GitService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx       = r.Context()
			name      = chi.URLParam(r, "name")
			namespace = chi.URLParam(r, "owner")
		)
		slug := namespace + "/" + name
		branches, err := gits.FindBranches(ctx, nil, slug)
		if err != nil {
			render.InternalError(w, err)
			return
		}
		render.JSON(w, branches, 200)
	}
}
