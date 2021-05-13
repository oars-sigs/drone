package main

import (
	spec "github.com/drone/drone/cmd/drone-server/config"
	"github.com/drone/drone/core"
	"github.com/drone/drone/handler/api"
	"github.com/drone/drone/handler/web"
	"github.com/drone/drone/metric"
	"github.com/drone/drone/plugin/config"
	"github.com/drone/drone/store/shared/db"
	"github.com/drone/go-scm/scm"
	"github.com/go-chi/chi"
	"github.com/google/wire"

	"github.com/oars-sigs/drone/handler/extendv1"
	"github.com/oars-sigs/drone/model"
	"github.com/oars-sigs/drone/services/git"
	"github.com/oars-sigs/drone/store/pipelines"
	extdb "github.com/oars-sigs/drone/store/shared/db"
	"github.com/oars-sigs/drone/store/templates"
)

var extendSet = wire.NewSet(
	templates.New,
	pipelines.New,
	extendv1.New,
	git.New,
)

// provideRouter is a Wire provider function that returns a
// router that is serves the provided handlers.
func provideRouter(api api.Server, apiext extendv1.Server, web web.Server, rpcv1 rpcHandlerV1, rpcv2 rpcHandlerV2, healthz healthzHandler, metrics *metric.Server, pprof pprofHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/healthz", healthz)
	r.Mount("/metrics", metrics)
	r.Mount("/api", api.Handler())
	r.Mount("/rpc/v2", rpcv2)
	r.Mount("/rpc", rpcv1)
	r.Mount("/", web.Handler())
	r.Mount("/debug", pprof)
	//@+++
	r.Mount("/extend", apiext.Handler())
	//@+++
	return r
}

// provideConfigPlugin is a Wire provider function that returns
// a yaml configuration plugin based on the environment
// configuration.
func provideConfigPlugin(client *scm.Client, contents core.FileService, pipeStore model.PipelineStore, conf spec.Config) core.ConfigService {
	return config.Combine(
		config.Memoize(
			config.Global(
				conf.Yaml.Endpoint,
				conf.Yaml.Secret,
				conf.Yaml.SkipVerify,
				conf.Yaml.Timeout,
			),
		),
		pipeStore,
	)
}

// provideDatabase is a Wire provider function that provides a
// database connection, configured from the environment.
func provideDatabase(config spec.Config) (*db.DB, error) {
	return extdb.Connect(
		config.Database.Driver,
		config.Database.Datasource,
	)
}
