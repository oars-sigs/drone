module github.com/oars-sigs/drone

go 1.14

replace (
	github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20200309214505-aa6a9891b09c+incompatible
	github.com/drone/drone/ui => ./ui
	github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.14
)

require (
	github.com/buildkite/yaml v2.1.0+incompatible // indirect
	github.com/drone/drone v1.9.2
	github.com/drone/drone-runtime v1.1.1-0.20200623162453-61e33e2cab5d
	github.com/drone/go-login v1.0.4-0.20190311170324-2a4df4f242a2
	github.com/drone/go-scm v1.7.2-0.20201028160627-427b8a85897c
	github.com/drone/signal v1.0.0
	github.com/go-chi/chi v3.3.3+incompatible
	github.com/go-chi/cors v1.0.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/uuid v1.1.2
	github.com/google/wire v0.4.0
	github.com/jmoiron/sqlx v0.0.0-20180614180643-0dae4fefe7c0
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.8.0
	github.com/mattn/go-sqlite3 v1.14.5
	github.com/sirupsen/logrus v1.7.0
	github.com/unrolled/secure v1.0.8
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
)
