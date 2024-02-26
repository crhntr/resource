package resource_test

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/crhntr/resource"
	"github.com/crhntr/resource/fakes"
	"github.com/crhntr/resource/internal/example"
)

//go:generate counterfeiter -generate
//counterfeiter:generate -o ./fakes/get.go --fake-name Get . getFunc
//counterfeiter:generate -o ./fakes/put.go --fake-name Put . putFunc
//counterfeiter:generate -o ./fakes/check.go --fake-name Check . checkFunc

type (
	// I need to put the config structures in another package "example" so counterfeiter can import them.
	// In a real implementation, you can put them in the same place as the function implementations.

	getFunc   resource.Get[example.Resource, example.GetParams, example.Version]
	putFunc   resource.Put[example.Resource, example.PutParams, example.Version]
	checkFunc resource.Check[example.Resource, example.Version]
)

var (
	// these assignments are helpful to ensure fakes have the correct Spy signature.
	// they also make the linter happier (the func types are flagged as unused by the linter even though they are used by counterfeiter)

	_ getFunc   = fakes.Get{}.Spy
	_ putFunc   = fakes.Put{}.Spy
	_ checkFunc = fakes.Check{}.Spy
)

const (
	checkStdin = `{
  "source": {
    "uri": "git://some-uri",
    "branch": "develop"
  },
  "version": { "ref": "pear" }
}`
	putStdin = `{
  "source": {
    "uri": "git://some-uri",
    "branch": "develop"
  },
  "params": {
	"ensure_checksum": true
  }
}`
	getStdin = `{
  "source": {
    "uri": "git://some-uri",
    "branch": "develop"
  },
  "params": {
	"include_zip": true
  },
  "version": {
  	"ref": "peach"
  }
}`
)

func TestRun_check(t *testing.T) {
	customization := resource.Customization{LoggerPrefix: "", LoggerFlags: 0, DisallowUnknownFields: true}

	get := new(fakes.Get)
	put := new(fakes.Put)
	check := new(fakes.Check)

	mux := resource.RunWithCustomization(customization, get.Spy, put.Spy, check.Spy)

	// language=json
	stdin := bytes.NewBufferString(checkStdin)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	err := mux(stdout, stderr, stdin, []string{"/some/absolute-path/check"})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if got := check.CallCount(); got != 1 {
		t.Fatalf("expected check to be called once, but it was called %d times", got)
	}
	if got := get.CallCount(); got != 0 {
		t.Fatalf("expected get not to be called, but it was called %d times", got)
	}
	if got := put.CallCount(); got != 0 {
		t.Fatalf("expected put not to be called, but it was called %d times", got)
	}

	_, resourceParamsArg, versionArg, _ := check.ArgsForCall(0)

	if exp := "develop"; resourceParamsArg.Branch != exp {
		t.Errorf("expected %q got %q", exp, resourceParamsArg.Branch)
	}
	if exp := "git://some-uri"; resourceParamsArg.URI != exp {
		t.Errorf("expected %q got %q", exp, resourceParamsArg.URI)
	}
	if exp := "pear"; versionArg.Ref != exp {
		t.Errorf("expected %q got %q", exp, versionArg.Ref)
	}
}

func TestRun_get(t *testing.T) {
	customization := resource.Customization{LoggerPrefix: "", LoggerFlags: 0, DisallowUnknownFields: true}

	get := new(fakes.Get)
	put := new(fakes.Put)
	check := new(fakes.Check)

	mux := resource.RunWithCustomization(customization, get.Spy, put.Spy, check.Spy)

	// language=json
	stdin := bytes.NewBufferString(getStdin)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	err := mux(stdout, stderr, stdin, []string{"/some/absolute-path/in"})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if got := check.CallCount(); got != 0 {
		t.Fatalf("expected check not to be called, but it was called %d times", got)
	}
	if got := get.CallCount(); got != 1 {
		t.Fatalf("expected get to be called once, but it was called %d times", got)
	}
	if got := put.CallCount(); got != 0 {
		t.Fatalf("expected put not to be called, but it was called %d times", got)
	}

	_, resourceParamsArg, getParamArg, versionArg, _ := get.ArgsForCall(0)

	if !getParamArg.IncludeZip {
		t.Errorf("expected get param args to be parsed")
	}
	if exp := "develop"; resourceParamsArg.Branch != exp {
		t.Errorf("expected %q got %q", exp, resourceParamsArg.Branch)
	}
	if exp := "git://some-uri"; resourceParamsArg.URI != exp {
		t.Errorf("expected %q got %q", exp, resourceParamsArg.URI)
	}
	if exp := "peach"; versionArg.Ref != exp {
		t.Errorf("expected %q got %q", exp, versionArg.Ref)
	}
}

func TestRun_put(t *testing.T) {
	customization := resource.Customization{LoggerPrefix: "", LoggerFlags: 0, DisallowUnknownFields: true}

	get := new(fakes.Get)
	put := new(fakes.Put)
	check := new(fakes.Check)

	mux := resource.RunWithCustomization(customization, get.Spy, put.Spy, check.Spy)

	// language=json
	stdin := bytes.NewBufferString(putStdin)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	err := mux(stdout, stderr, stdin, []string{"/some/absolute-path/out"})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if got := check.CallCount(); got != 0 {
		t.Fatalf("expected check not to be called, but it was called %d times", got)
	}
	if got := get.CallCount(); got != 0 {
		t.Fatalf("expected get not to be called, but it was called %d times", got)
	}
	if got := put.CallCount(); got != 1 {
		t.Fatalf("expected put to be called once, but it was called %d times", got)
	}

	_, resourceParamsArg, getParamArg, _ := put.ArgsForCall(0)

	if !getParamArg.EnsureChecksum {
		t.Errorf("expected get param args to be parsed")
	}
	if exp := "develop"; resourceParamsArg.Branch != exp {
		t.Errorf("expected %q got %q", exp, resourceParamsArg.Branch)
	}
	if exp := "git://some-uri"; resourceParamsArg.URI != exp {
		t.Errorf("expected %q got %q", exp, resourceParamsArg.URI)
	}
}

func TestRun_failure_cases(t *testing.T) {
	t.Run("write to stdout fails", func(t *testing.T) {
		customization := resource.Customization{LoggerPrefix: "", LoggerFlags: 0, DisallowUnknownFields: true}

		get := new(fakes.Get)
		put := new(fakes.Put)
		check := new(fakes.Check)

		mux := resource.RunWithCustomization(customization, get.Spy, put.Spy, check.Spy)

		stdin := bytes.NewBufferString(putStdin)
		stdout := errWriter{}
		stderr := new(bytes.Buffer)
		err := mux(stdout, stderr, stdin, []string{"/some/absolute-path/out"})

		if exp := "write banana"; err == nil || !strings.Contains(err.Error(), exp) {
			t.Fatalf("expected error containing %q got %s", exp, err)
		}
	})

	t.Run("read from stdin fails", func(t *testing.T) {
		customization := resource.Customization{LoggerPrefix: "", LoggerFlags: 0, DisallowUnknownFields: true}

		get := new(fakes.Get)
		put := new(fakes.Put)
		check := new(fakes.Check)

		mux := resource.RunWithCustomization(customization, get.Spy, put.Spy, check.Spy)

		stdin := iotest.ErrReader(fmt.Errorf("read banana"))
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		err := mux(stdout, stderr, stdin, []string{"/some/absolute-path/out"})

		if exp := "read banana"; err == nil || !strings.Contains(err.Error(), exp) {
			t.Fatalf("expected error containing %q got %s", exp, err)
		}
	})

	t.Run("in fails", func(t *testing.T) {
		customization := resource.Customization{LoggerPrefix: "", LoggerFlags: 0, DisallowUnknownFields: true}

		get := new(fakes.Get)
		put := new(fakes.Put)
		check := new(fakes.Check)

		get.Returns(nil, fmt.Errorf("get banana"))

		mux := resource.RunWithCustomization(customization, get.Spy, put.Spy, check.Spy)

		stdin := strings.NewReader(getStdin)
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		err := mux(stdout, stderr, stdin, []string{"/some/absolute-path/in"})

		if exp := "get banana"; err == nil || !strings.Contains(err.Error(), exp) {
			t.Fatalf("expected error containing %q got %s", exp, err)
		}
	})

	t.Run("out fails", func(t *testing.T) {
		customization := resource.Customization{LoggerPrefix: "", LoggerFlags: 0, DisallowUnknownFields: true}

		get := new(fakes.Get)
		put := new(fakes.Put)
		check := new(fakes.Check)

		put.Returns(example.Version{}, nil, fmt.Errorf("put banana"))

		mux := resource.RunWithCustomization(customization, get.Spy, put.Spy, check.Spy)

		stdin := strings.NewReader(putStdin)
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		err := mux(stdout, stderr, stdin, []string{"/some/absolute-path/out"})

		if exp := "put banana"; err == nil || !strings.Contains(err.Error(), exp) {
			t.Fatalf("expected error containing %q got %s", exp, err)
		}
	})

	t.Run("check fails", func(t *testing.T) {
		customization := resource.Customization{LoggerPrefix: "", LoggerFlags: 0, DisallowUnknownFields: true}

		get := new(fakes.Get)
		put := new(fakes.Put)
		check := new(fakes.Check)

		check.Returns(nil, fmt.Errorf("check banana"))

		mux := resource.RunWithCustomization(customization, get.Spy, put.Spy, check.Spy)

		stdin := strings.NewReader(checkStdin)
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		err := mux(stdout, stderr, stdin, []string{"/some/absolute-path/check"})

		if exp := "check banana"; err == nil || !strings.Contains(err.Error(), exp) {
			t.Fatalf("expected error containing %q got %s", exp, err)
		}
	})
}

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, fmt.Errorf("write banana")
}

func TestRun_Logger(t *testing.T) {
	get := new(fakes.Get)
	put := new(fakes.Put)
	check := new(fakes.Check)

	mux := resource.Run(get.Spy, put.Spy, check.Spy)

	stdin := strings.NewReader(checkStdin)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	_ = mux(stdout, stderr, stdin, []string{"/some/absolute-path/check"})

	_, _, _, logger := check.ArgsForCall(0)

	if logger.Flags() != log.Default().Flags() {
		t.Errorf("expected default flags")
	}

	if logger.Prefix() != log.Default().Prefix() {
		t.Errorf("expected default prefix")
	}
}
