package resource

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"path/filepath"
)

type MetadataField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type (
	Get[ResourceParams, GetParams, Version any] func(context.Context, *log.Logger, ResourceParams, GetParams, Version, string) ([]MetadataField, error)
	Put[ResourceParams, PutParams, Version any] func(context.Context, *log.Logger, ResourceParams, PutParams, string) (Version, []MetadataField, error)
	Check[ResourceParams, Version any]          func(context.Context, *log.Logger, ResourceParams, Version) ([]Version, error)
)

// Run calls the given Get, Put, and Check functions based on the command name.
// You probably want to call this in your main function like so:
//
//		func main() {
//		  cmd := resource.Run(get, put, check)
//	      if err := cmd(os.Stdout, os.Stderr, os.Stdin, os.Args); err != nil {
//	        log.Fatal(err)
//	      }
//		}
func Run[ResourceParams, GetParams, PutParams, Version any](
	in Get[ResourceParams, GetParams, Version],
	out Put[ResourceParams, PutParams, Version],
	check Check[ResourceParams, Version],
) func(stdout, stderr io.Writer, stdin io.Reader, args []string) error {
	return RunWithCustomization(Customization{
		LoggerPrefix: log.Default().Prefix(),
		LoggerFlags:  log.Default().Flags(),
	}, in, out, check)
}

// Customization allows you to configure the behavior of the Run function.
// In particular the zero value Customization can be helpful for getting deterministic logging output in testing (without any timestamps).
type Customization struct {
	LoggerPrefix          string
	LoggerFlags           int
	DisallowUnknownFields bool
}

// RunWithCustomization calls the given Get, Put, and Check functions based on the command name.
func RunWithCustomization[ResourceParams, GetParams, PutParams, Version any](
	customization Customization,
	in Get[ResourceParams, GetParams, Version],
	out Put[ResourceParams, PutParams, Version],
	check Check[ResourceParams, Version],
) func(stdout, stderr io.Writer, stdin io.Reader, args []string) error {
	return func(stdout io.Writer, stderr io.Writer, stdin io.Reader, args []string) error {
		ctx := context.Background()
		stderrLogger := log.New(stderr, customization.LoggerPrefix, customization.LoggerFlags)
		var err error
		switch filepath.Base(args[0]) {
		case "in":
			err = handleJSON(ctx, customization, stdout, stderrLogger, stdin, args[1:], in.run)
		case "out":
			err = handleJSON(ctx, customization, stdout, stderrLogger, stdin, args[1:], out.run)
		case "check":
			err = handleJSON(ctx, customization, stdout, stderrLogger, stdin, args[1:], check.run)
		}
		return err
	}
}

func handleJSON[Req, Res any](ctx context.Context, bc Customization, stdout io.Writer, log *log.Logger, stdin io.Reader, args []string, run func(context.Context, *log.Logger, Req, []string) (Res, error)) error {
	var req Req
	dec := json.NewDecoder(stdin)
	if bc.DisallowUnknownFields {
		dec.DisallowUnknownFields()
	}
	err := dec.Decode(&req)
	if err != nil {
		return err
	}
	res, err := run(ctx, log, req, args)
	if err != nil {
		return err
	}
	return json.NewEncoder(stdout).Encode(res)
}

type inRequest[ResourceParams, InParams, Version any] struct {
	Source  ResourceParams `json:"source"`
	Params  InParams       `json:"params"`
	Version Version        `json:"version"`
}

type inResponse[Version any] struct {
	Version         Version         `json:"version"`
	VersionMetadata []MetadataField `json:"metadata,omitempty"`
}

func (in Get[ResourceParams, GetParams, Version]) run(ctx context.Context, log *log.Logger, req inRequest[ResourceParams, GetParams, Version], args []string) (inResponse[Version], error) {
	m, err := in(ctx, log, req.Source, req.Params, req.Version, args[0])
	return inResponse[Version]{Version: req.Version, VersionMetadata: m}, err
}

type outRequest[ResourceParams, PutParams, Version any] struct {
	Source ResourceParams `json:"source"`
	Params PutParams      `json:"params"`
}

type outResponse[Version any] struct {
	Version         Version         `json:"version"`
	VersionMetadata []MetadataField `json:"metadata"`
}

func (out Put[ResourceParams, PutParams, Version]) run(ctx context.Context, log *log.Logger, req outRequest[ResourceParams, PutParams, Version], args []string) (outResponse[Version], error) {
	v, m, err := out(ctx, log, req.Source, req.Params, args[0])
	return outResponse[Version]{Version: v, VersionMetadata: m}, err
}

type checkRequest[ResourceParams, Version any] struct {
	Source  ResourceParams `json:"source"`
	Version Version        `json:"version"`
}

type checkResponse[Version any] []Version

func (fn Check[ResourceParams, Version]) run(ctx context.Context, log *log.Logger, req checkRequest[ResourceParams, Version], _ []string) (checkResponse[Version], error) {
	return fn(ctx, log, req.Source, req.Version)
}
