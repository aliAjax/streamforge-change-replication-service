#!/usr/bin/env sh
set -eu
base="${BASE_URL:-http://127.0.0.1:8092}"
auth="Authorization: Bearer ${STREAMFORGE_API_KEY:-dev-key}"
curl -fsS "$base/healthz"
curl -fsS "$base/readyz"
curl -fsS -H "$auth" "$base/metrics"
curl -fsS -X POST -H "$auth" -H 'Content-Type: application/json' "$base/api/v1/sources" -d '{"id":"sim-postgres","name":"simulated PostgreSQL","kind":"postgres","endpoint":"simulated://postgres","database":"app"}'
curl -fsS -X POST -H "$auth" -H 'Content-Type: application/json' "$base/api/v1/pipelines" -d '{"id":"pipeline-17","name":"demo","source_id":"sim-postgres","strategy":"compatible","destination":"memory"}'
curl -fsS -X POST -H "$auth" -H 'Content-Type: application/json' "$base/api/v1/simulate/transaction" -d '{"source":"sim-postgres","operation":"insert","after":{"id":1,"name":"alpha"},"primary_key":{"id":1}}'
curl -fsS -H "$auth" "$base/api/v1/events?cursor=0&limit=10"
