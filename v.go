package gitlab

import (
	"os/exec"
	"strings"
)

// LibVersion is a string returned by `git describe --tags`
func LibVersion(gitdir string) string {
	cmd := exec.Command("git", "describe", "--tags")
	cmd.Dir = gitdir

	tag, err := cmd.Output()
	if err != nil {
		// FIXME: allow no tags
		version := ""
		return version
	}

	version := strings.TrimSpace(string(tag))

	return version
}
