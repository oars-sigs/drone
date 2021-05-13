package main

import (
	"net/http"
	"os"

	spec "github.com/drone/drone/cmd/drone-server/config"
	"github.com/drone/drone/core"
	"github.com/drone/drone/handler/api"
	"github.com/drone/drone/handler/web"
	"github.com/drone/drone/metric"
	"github.com/drone/drone/plugin/config"
	"github.com/drone/drone/store/shared/db"
	"github.com/drone/go-login/login"
	"github.com/drone/go-login/login/gitlab"
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/transport/oauth2"
	"github.com/go-chi/chi"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"

	"github.com/oars-sigs/drone/handler/extendv1"
	"github.com/oars-sigs/drone/model"
	"github.com/oars-sigs/drone/pkg/scm/driver/gitee"
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

// provideBitbucketClient is a Wire provider function that
// returns a Source Control Management client based on the
// environment configuration.
func provideClient(config spec.Config) *scm.Client {
	switch {
	case config.Bitbucket.ClientID != "":
		return provideBitbucketClient(config)
	case config.Github.ClientID != "":
		return provideGithubClient(config)
	case config.Gitea.Server != "":
		return provideGiteaClient(config)
	case config.GitLab.ClientID != "":
		return provideGitlabClient(config)
	case config.Gogs.Server != "":
		return provideGogsClient(config)
	case config.Stash.ConsumerKey != "":
		return provideStashClient(config)
	//@+++
	case os.Getenv("DRONE_GITEE_CLIENT_ID") != "":
		return provideGiteeClient(config)
		//@+++
	}
	logrus.Fatalln("main: source code management system not configured")
	return nil
}

// provideGiteeClient is a Wire provider function that returns
// a Gitee client based on the environment configuration.
func provideGiteeClient(config spec.Config) *scm.Client {
	server := "https://gitee.com"
	logrus.WithField("server", server).
		WithField("skip_verify", false).
		Debugln("main: creating the Gitee client")

	client, err := gitee.New(server)
	if err != nil {
		logrus.WithError(err).
			Fatalln("main: cannot create the Gitee client")
	}

	client.Client = &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.ContextTokenSource(),
			Base:   defaultTransport(false),
		},
	}
	return client
}

// provideGiteeLogin is a Wire provider function that returns
// a Gitee authenticator based on the environment configuration.
func provideGiteeLogin(config spec.Config) login.Middleware {
	clientID := os.Getenv("DRONE_GITEE_CLIENT_ID")
	clientSecret := os.Getenv("DRONE_GITEE_CLIENT_SECRET")
	server := "https://gitee.com"
	return &gitlab.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  config.Server.Addr + "/login",
		Server:       server,
		Client:       defaultClient(false),
	}
}

// provideLogin is a Wire provider function that returns an
// authenticator based on the environment configuration.
func provideLogin(config spec.Config) login.Middleware {
	switch {
	case config.Bitbucket.ClientID != "":
		return provideBitbucketLogin(config)
	case config.Github.ClientID != "":
		return provideGithubLogin(config)
	case config.Gitea.Server != "":
		return provideGiteaLogin(config)
	case config.GitLab.ClientID != "":
		return provideGitlabLogin(config)
	case config.Gogs.Server != "":
		return provideGogsLogin(config)
	case config.Stash.ConsumerKey != "":
		return provideStashLogin(config)
	//@+++
	case os.Getenv("DRONE_GITEE_CLIENT_ID") != "":
		return provideGiteeLogin(config)
		//@+++
	}
	logrus.Fatalln("main: source code management system not configured")
	return nil
}
