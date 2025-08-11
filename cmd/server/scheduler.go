package main

import (
	"sync"
	"time"
	"log"
	"errors"
	"strings"
	"os/exec"
	pb "detf/api"
)

type Test struct {
	baseline  pb.Engine
	candidate pb.Engine
	active bool

	book   uint
	wins   uint
	losses uint
	draws  uint
}

var mu    sync.Mutex
var tests []Test

func HandleResult(res pb.Result) {
	mu.Lock()
	defer mu.Unlock()

	if res.Win == res.Draw {
		log.Printf("Received invalid result")
		return
	}

	idx, err := TestIndex(*res.GetBaseline(), *res.GetCandidate())
	if err != nil {
		log.Printf("Received result from unknown test")
		return
	}

	if res.Win {
		tests[idx].wins += 1
	}

	if res.Draw {
		tests[idx].draws += 1
	}
}

func NextMatch() (pb.Match, error) {
	mu.Lock()
	defer mu.Unlock()

	for _, test := range tests {
		if !test.active {
			continue
		}
		pos       := book[test.book - test.book % 2]
		turn      := test.book % 2 == 0
		test.book += 1
		return pb.Match {
			Baseline:  &test.baseline,
			Candidate: &test.candidate,
			Pos:       pos,
			Turn:      turn,
		}, nil
	}
	return pb.Match {}, errors.New("No active tests")
}

func TestIndex(baseline pb.Engine, candidate pb.Engine) (int, error) {
	for i, test := range tests {
		if baseline.GetRepo()  == test.baseline.GetRepo()  &&
		   baseline.GetRef()   == test.baseline.GetRef()   &&
		   candidate.GetRepo() == test.candidate.GetRepo() &&
		   candidate.GetRef()  == test.candidate.GetRef() {
			return i, nil
		}
	}
	return 0, errors.New("Test not found")
}

func TestsContain(baseline pb.Engine, candidate pb.Engine) bool {
	_, err := TestIndex(baseline, candidate)
	return err == nil
}

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

func CheckRefStatus(repo string) {
	master, refs, err := FetchRefs(repo)
	if err != nil {
		log.Printf("Failed to retrieve repo refs: %v", err)
		return
	}
	for _, ref := range refs {
		baseline  := pb.Engine { Repo: repo, Ref: master }
		candidate := pb.Engine { Repo: repo, Ref: ref }
		if TestsContain(baseline, candidate) {
			continue
		}
		log.Printf("Found new repo ref %s - %s", repo, ref)
		mu.Lock()
		tests = append(tests, Test {
			baseline:  baseline,
			candidate: candidate,
			active:    true,
		})
		mu.Unlock()
	}
}

func PrintStatus() {
	for _, test := range tests {
		log.Printf(
			"%s - %s: Total: %d W: %d L: %d D: %d",
			test.candidate.GetRepo(),
			test.candidate.GetRef(),
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
