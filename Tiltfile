# -*- mode: Python -*-

load('ext://helm_resource', 'helm_resource', 'helm_repo')

compile_opt = 'GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 '

# Spin up spanner
helm_resource(
    'gcp-spanner-emulator',
    'spanner/infra',
    port_forwards=["9010:9010", "9020:9020"],
    flags=[
    '-f',
    'spanner/infra/values-spanner-dev.yaml',
  ],
  labels='spanner',
)

# Create spanner db
local_resource(
  'spanner-init', 
  "go run spanner/script/main.go", 
  deps=["spanner/script", "spanner/schema"],
  resource_deps=["gcp-spanner-emulator"],
  labels=["spanner"],
)

# Run pagination
local_resource(
  'pagination', 
  "go run pagination/main.go", 
  deps=["pagination"],
  resource_deps=["spanner-init"],
  labels=["pagination"],
)
