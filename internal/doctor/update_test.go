package doctor

import (
	"os/exec"
	"strings"
	"testing"
)

func TestRunnerUpdatesList(t *testing.T) {
	got := RunnerUpdates()
	if len(got) != 3 {
		t.Fatalf("want 3 runners, got %d", len(got))
	}
	for _, u := range got {
		if !strings.HasSuffix(u.NpmPkg, "@latest") {
			t.Fatalf("%s missing @latest: %s", u.Key, u.NpmPkg)
		}
	}
}

func TestRunnerPill(t *testing.T) {
	if got := (RunnerUpdateStatus{Updating: true}).Pill(); got != "[UPDATING RUNNERS...]" {
		t.Fatalf("pill=%q", got)
	}
	if got := (RunnerUpdateStatus{Done: true}).Pill(); got != "[RUNNERS UP-TO-DATE]" {
		t.Fatalf("pill=%q", got)
	}
}

func TestUpdateRunnerFakeExec(t *testing.T) {
	old := execCommand
	defer func() { execCommand = old }()
	execCommand = func(name string, args ...string) *exec.Cmd {
		// echo always succeeds; simulates npm without network.
		return exec.Command("go", "version")
	}
	if err := UpdateRunner(RunnerUpdate{Key: "test", NpmPkg: "x@latest"}, func(string) {}); err != nil {
		t.Fatal(err)
	}
}
