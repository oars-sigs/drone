package extendv1

import (
	"net/http"

	"github.com/oars-sigs/drone/handler/extendv1/repos/builds"
	"github.com/oars-sigs/drone/handler/extendv1/repos/pipelines"
	"github.com/oars-sigs/drone/handler/extendv1/repos/ref"
	"github.com/oars-sigs/drone/handler/extendv1/templates"
	"github.com/oars-sigs/drone/model"
	"github.com/oars-sigs/drone/ui/dist"

	"github.com/drone/drone/core"
	"github.com/drone/drone/handler/api/acl"
	"github.com/drone/drone/handler/api/auth"
	"github.com/drone/drone/logger"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

var corsOpts = cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300,
}

func New(
	builds core.BuildStore,
	commits core.CommitService,
	cron core.CronStore,
	events core.Pubsub,
	globals core.GlobalSecretStore,
	hooks core.HookService,
	logs core.LogStore,
	license *core.License,
	licenses core.LicenseService,
	perms core.PermStore,
	repos core.RepositoryStore,
	repoz core.RepositoryService,
	scheduler core.Scheduler,
	secrets core.SecretStore,
	stages core.StageStore,
	steps core.StepStore,
	status core.StatusService,
	session core.Session,
	stream core.LogStream,
	syncer core.Syncer,
	system *core.System,
	triggerer core.Triggerer,
	users core.UserStore,
	webhook core.WebhookSender,
	tmpls model.TemplateStore,
	pipelineStore model.PipelineStore,
	gits model.GitService,
) Server {
	return Server{
		Builds:    builds,
		Cron:      cron,
		Commits:   commits,
		Events:    events,
		Globals:   globals,
		Hooks:     hooks,
		Logs:      logs,
		License:   license,
		Licenses:  licenses,
		Perms:     perms,
		Repos:     repos,
		Repoz:     repoz,
		Scheduler: scheduler,
		Secrets:   secrets,
		Stages:    stages,
		Steps:     steps,
		Status:    status,
		Session:   session,
		Stream:    stream,
		Syncer:    syncer,
		System:    system,
		Triggerer: triggerer,
		Users:     users,
		Webhook:   webhook,

		Tmpls:         tmpls,
		PipelineStore: pipelineStore,
		gits:          gits,
	}
}

// Server is a http.Handler which exposes drone functionality over HTTP.
type Server struct {
	Builds    core.BuildStore
	Cron      core.CronStore
	Commits   core.CommitService
	Events    core.Pubsub
	Globals   core.GlobalSecretStore
	Hooks     core.HookService
	Logs      core.LogStore
	License   *core.License
	Licenses  core.LicenseService
	Perms     core.PermStore
	Repos     core.RepositoryStore
	Repoz     core.RepositoryService
	Scheduler core.Scheduler
	Secrets   core.SecretStore
	Stages    core.StageStore
	Steps     core.StepStore
	Status    core.StatusService
	Session   core.Session
	Stream    core.LogStream
	Syncer    core.Syncer
	System    *core.System
	Triggerer core.Triggerer
	Users     core.UserStore
	Webhook   core.WebhookSender

	Tmpls         model.TemplateStore
	PipelineStore model.PipelineStore
	gits          model.GitService
}

// Handler returns an http.Handler
func (s Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)
	r.Use(logger.Middleware)
	r.Use(auth.HandleAuthentication(s.Session))
	r.Use(acl.AuthorizeUser)

	cors := cors.New(corsOpts)
	r.Use(cors.Handler)

	r.Route("/templates", func(r chi.Router) {
		r.Get("/", templates.HandleGetTemp(s.Tmpls))
		r.Post("/", templates.HandleCreateTemp(s.Tmpls))
		r.Route("/{uuid}", func(r chi.Router) {
			r.Get("/", templates.HandleFindTemp(s.Tmpls))
			r.Put("/", templates.HandlePutTemp(s.Tmpls))
			r.Delete("/", templates.HandleDeleteTemp(s.Tmpls))
		})

	})

	r.Route("/{owner}/{name}", func(r chi.Router) {
		r.Get("/branches", ref.HandleFindBranches(s.Repos, s.gits))
		r.Get("/tags", ref.HandleFindTags(s.Repos, s.gits))
		r.Route("/builds", func(r chi.Router) {
			r.Post("/", builds.HandleCreate(s.Users, s.Repos, s.Commits, s.Triggerer))
		})
		r.Route("/pipelines", func(r chi.Router) {
			r.Get("/", pipelines.HandleFindPipelines(s.Repos, s.PipelineStore))
			r.Put("/", pipelines.HandlePutPipeline(s.Repos, s.PipelineStore))
		})

	})
	return r
}

func (s Server) Web() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)
	r.Use(logger.Middleware)
	h := http.FileServer(dist.New())
	r.Handle("/*filepath", h)
	return r
}
