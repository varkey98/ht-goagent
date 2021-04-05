module github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/jackc/hyperpgx

go 1.15

require (
	github.com/hypertrace/goagent v0.0.0-00010101000000-000000000000
	github.com/jackc/pgconn v1.8.1
	github.com/jackc/pgtype v1.7.0
	github.com/jackc/pgx/v4 v4.11.0
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel/trace v0.18.0
)

replace github.com/hypertrace/goagent => ../../../../../