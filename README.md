# Concourse Resource Boilerplate  [![Go Reference](https://pkg.go.dev/badge/github.com/crhntr/resource.svg)](https://pkg.go.dev/github.com/crhntr/resource)

A package to reduce boilerplate when writing [custom Concourse Resource Types](https://concourse-ci.org/implementing-resource-types.html).

## Start with this

```go

package main

import (
	"context"
	"log"
	"os"

	"github.com/crhntr/resource"
)

func main() {
	cmd := resource.Run(get, put, check)
	if err := cmd(os.Stdout, os.Stderr, os.Stdin, os.Args); err != nil {
		log.Fatal(err)
	}
}

type (
	ResourceParams struct{}
	GetParams      struct{}
	PutParams      struct{}
	Version        struct{}
)

func get(context.Context, *log.Logger, ResourceParams, GetParams, Version, string) ([]resource.MetadataField, error) {
	panic("not implemented")
}

func put(context.Context, *log.Logger, ResourceParams, PutParams, string) (Version, []resource.MetadataField, error) {
	panic("not implemented")
}

func check(context.Context, *log.Logger, ResourceParams, Version) ([]Version, error) {
	panic("not implemented")
}

```
