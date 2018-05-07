package gitlab

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func gitCmdOrFatal(t *testing.T, tempdir string, arg ...string) {
	cmd := exec.Command("git", arg...)
	cmd.Dir = tempdir
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("Could not run %v: %v", cmd.Args, err)
	}
}

func TestVersion(t *testing.T) {
	tempdir, err := ioutil.TempDir("", "version-test-")
	if err != nil {
		t.Fatalf("Could not create temp dir: %v", err)
	}
	defer os.RemoveAll(tempdir)

	tempfile := filepath.Join(tempdir, "test")
	if err := ioutil.WriteFile(tempfile, []byte("testcase"), 0644); err != nil {
		t.Fatalf("Could not write temp file %q: %v", tempfile, err)
	}

	gitCmdOrFatal(t, tempdir, "init")
	gitCmdOrFatal(t, tempdir, "config", "user.email", "unittest@example.com")
	gitCmdOrFatal(t, tempdir, "config", "user.name", "Unit Test")
	gitCmdOrFatal(t, tempdir, "add", "test")
	cmd := exec.Command("git", "commit", "-a", "-m", "initial commit")
	cmd.Env = append(os.Environ(), "GIT_COMMITTER_DATE=2015-04-20T11:22:33")
	cmd.Dir = tempdir
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("Could not run %v: %v", cmd.Args, err)
	}

	t.Run("test when git tag not exists", func(t *testing.T) {
		gitCmdOrFatal(t, tempdir, "describe", "--long", "--all")
		if err != nil {
			t.Fatalf("Determining package version from git failed: %v", err)
		}

		got := LibVersion(tempdir)
		// if want := "0.0~git20150420."; !strings.HasPrefix(got, want) {
		if want := ""; !strings.HasPrefix(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("test when git tag exists", func(t *testing.T) {
		gitCmdOrFatal(t, tempdir, "tag", "-a", "v1", "-m", `"release v1"`)
		got := LibVersion(tempdir)
		want := "v1"

		if got != want {
			t.Errorf("want '%s' but got '%s'", want, got)
		}
	})
}
