# -*- mode: Python -*-

load('ext://restart_process', 'docker_build_with_restart')
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
  'spanner', 
  "go run spanner/script/main.go", 
  deps=["spanner/script", "spanner/schema"],
  dir=".",
  resource_deps=["gcp-spanner-emulator"],
  labels=["spanner"],
)
