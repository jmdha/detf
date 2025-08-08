package main

import (
	"sync"
	"time"
	"log"
	"errors"
	"strings"
	"os/exec"
)

type Test struct {
	repo   string
	ref    string
	active bool

	book   uint
	wins   uint
	losses uint
	draws  uint
}

type Match struct {
	repo string
	ref  string
	pos  string
	turn bool
}

type Result struct {
	repo string
	ref  string
	win  bool
	draw bool
}

var mu    sync.Mutex
var tests []Test

func HandleResult(res Result) {
	mu.Lock()
	defer mu.Unlock()

	if res.win == res.draw {
		log.Printf("Received invalid result")
		return
	}

	idx, err := TestIndex(res.repo, res.ref)
	if err != nil {
		log.Printf("Received result from unknown test: %s - %s", res.repo, res.ref)
		return
	}

	if res.win {
		tests[idx].wins += 1
	}

	if res.draw {
		tests[idx].draws += 1
	}
}

func NextMatch() (Match, error) {
	mu.Lock()
	defer mu.Unlock()

	for _, test := range tests {
		if !test.active {
			continue
		}
		pos       := book[test.book - test.book % 2]
		turn      := test.book % 2 == 0
		test.book += 1
		return Match {
			repo: test.repo,
			ref:  test.ref,
			pos:  pos,
			turn: turn,
		}, nil
	}
	return Match {}, errors.New("No active tests")
}

func TestIndex(repo string, ref string) (int, error) {
	for i, test := range tests {
		if repo == test.repo && ref == test.ref {
			return i, nil
		}
	}
	return 0, errors.New("Test not found")
}

func TestsContain(repo string, ref string) bool {
	_, err := TestIndex(repo, ref)
	return err == nil
}

func FetchRefs(repo string) ([]string, error) {
	cmd := exec.Command("git", "ls-remote", repo, "pull/*/head")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	str   := string(out)
	lines := strings.Split(strings.TrimSuffix(str, "\n"), "\n")
	refs  := make([]string, len(lines))
	for i, line := range lines {
		refs[i] = strings.Fields(line)[1]
	}
	return refs, nil
}

func CheckRefStatus(repo string) {
	refs, err := FetchRefs(repo)
	if err != nil {
		log.Printf("Failed to retrieve repo refs: %v", err)
		return
	}
	for _, ref := range refs {
		if TestsContain(repo, ref) {
			continue
		}
		log.Printf("Found new repo ref %s - %s", repo, ref)
		mu.Lock()
		tests = append(tests, Test {
			repo:   repo,
			ref:    ref,
			active: true,
		})
		mu.Unlock()
	}
}

func PrintStatus() {
	for _, test := range tests {
		log.Printf(
			"%s - %s: Total: %d W: %d L: %d D: %d",
			test.repo,
			test.ref,
			test.wins + test.losses + test.draws,
			test.wins,
			test.losses,
			test.draws,
		)
	}
}

func SchedulerStart(repo string) {
	for {
		// Sleep to avoid sending too many requests
		time.Sleep(10 * time.Second)
		CheckRefStatus(repo)
		PrintStatus()
	}
}
