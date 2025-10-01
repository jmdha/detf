package main

import (
	"strings"
	"os/exec"
)

func FetchRefs(repo string) (string, []string, error) {
	// Find master ref
	out_m, err := exec.Command("git", "ls-remote", repo, "HEAD").Output()
	if err != nil {
		return "", nil, err
	}
	master := strings.Fields(string(out_m))[0][0:6]

	// Find PR refs
	out_r, err := exec.Command("git", "ls-remote", repo, "pull/*/head").Output()
	if err != nil {
		return "", nil, err
	}
	lines := strings.Split(strings.TrimSuffix(string(out_r), "\n"), "\n")
	refs  := make([]string, len(lines))
	for i, line := range lines {
		refs[i] = strings.Fields(line)[0][0:6]
	}
	return master, refs, nil
}
