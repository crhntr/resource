package resource

import "os"

// Build is a convenience type to help you load build metadata from the environment in Get and Put functions.
// It must only be called in Get and Put. It will always return ("", false) in Check. Concourse does not set
// build variables in check calls.
type Build struct{}

func (Build) ID() (string, bool)             { return os.LookupEnv("BUILD_ID") }
func (Build) Name() (string, bool)           { return os.LookupEnv("BUILD_NAME") }
func (Build) JobName() (string, bool)        { return os.LookupEnv("BUILD_JOB_NAME") }
func (Build) PipelineName() (string, bool)   { return os.LookupEnv("BUILD_PIPELINE_NAME") }
func (Build) TeamName() (string, bool)       { return os.LookupEnv("BUILD_TEAM_NAME") }
func (Build) CreatedBy() (string, bool)      { return os.LookupEnv("BUILD_CREATED_BY") }
func (Build) ATCExternalURL() (string, bool) { return os.LookupEnv("ATC_EXTERNAL_URL") }
