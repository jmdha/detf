package main

import (
	"regexp"
	"strings"
	"os/exec"
)

var (
	abc = regexp.MustCompile(`[0-9]+`)
)

func FetchRefs(repo string, filter string) ([]string, error) {
	var out []string

	eout, err := exec.Command("git", "ls-remote", repo, filter).Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSuffix(string(eout), "\n"), "\n")
	for _, ref := range lines {
		if len(ref) == 0 {
			continue
		}
		out = append(out, ref)
	}

	return out, nil
}

func FetchHead(repo string) (string, error) {
	head, err := FetchRefs(repo, "HEAD")
	if err != nil {
		return "", err
	}
	return strings.Fields(head[0])[0][0:6], nil
}

func FetchPRs(repo string) ([]string, error) {
	var out  []string

	heads, err := FetchRefs(repo, "refs/pull/*/head")
	if err != nil {
		return nil, err
	}
	head_refs, head_ids := SplitStrings(heads)

	merges, err := FetchRefs(repo, "refs/pull/*/merge")
	if err != nil {
		return nil, err
	}
	_, merge_ids := SplitStrings(merges)

	for _, merge_id := range merge_ids {
		merge_num := abc.FindAllString(merge_id, -1)[0]
		for idx, head_id := range head_ids {
			head_num := abc.FindAllString(head_id, -1)[0]
			if merge_num == head_num {
				out = append(out, head_refs[idx][0:6])
			}
		}
	}

	return out, nil
}
