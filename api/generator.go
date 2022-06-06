package api

import (
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
)

//go:generate go run github.com/go-swagger/go-swagger/cmd/swagger generate client -A podman --with-flatten=full --skip-tag-packages -f swagger-v4.1.yaml
