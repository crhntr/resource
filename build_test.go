package resource

import "testing"

func TestBuild(t *testing.T) {
	var b Build

	t.Setenv("BUILD_ID", "id")
	t.Setenv("BUILD_NAME", "name")
	t.Setenv("BUILD_JOB_NAME", "jobname")
	t.Setenv("BUILD_PIPELINE_NAME", "pipelinename")
	t.Setenv("BUILD_TEAM_NAME", "teamname")
	t.Setenv("BUILD_CREATED_BY", "createdby")
	t.Setenv("ATC_EXTERNAL_URL", "atcexternalurl")

	if val, found := b.ID(); !found {
		t.Errorf("expected value to be found")
	} else if exp := "id"; val != exp {
		t.Errorf("expected ID to return %q but got %q", exp, val)
	}
	if val, found := b.Name(); !found {
		t.Errorf("expected value to be found")
	} else if exp := "name"; val != exp {
		t.Errorf("expected Name to return %q but got %q", exp, val)
	}
	if val, found := b.JobName(); !found {
		t.Errorf("expected value to be found")
	} else if exp := "jobname"; val != exp {
		t.Errorf("expected JobName to return %q but got %q", exp, val)
	}
	if val, found := b.PipelineName(); !found {
		t.Errorf("expected value to be found")
	} else if exp := "pipelinename"; val != exp {
		t.Errorf("expected PipelineName to return %q but got %q", exp, val)
	}
	if val, found := b.TeamName(); !found {
		t.Errorf("expected value to be found")
	} else if exp := "teamname"; val != exp {
		t.Errorf("expected TeamName to return %q but got %q", exp, val)
	}
	if val, found := b.CreatedBy(); !found {
		t.Errorf("expected value to be found")
	} else if exp := "createdby"; val != exp {
		t.Errorf("expected CreatedBy to return %q but got %q", exp, val)
	}
	if val, found := b.ATCExternalURL(); !found {
		t.Errorf("expected value to be found")
	} else if exp := "atcexternalurl"; val != exp {
		t.Errorf("expected ATCExternalURL to return %q but got %q", exp, val)
	}
}
