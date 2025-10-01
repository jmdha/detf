package main

import (
	"sync"
	"time"
	"log"
	"errors"
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

	idx, err := TestIndex(*res.GetBaseline(), *res.GetCandidate())
	if err != nil {
		log.Printf("Received result from unknown test")
		return
	}

	if res.Draw {
		tests[idx].draws += 1
	} else if res.Win {
		tests[idx].wins   += 1
	} else {
		tests[idx].losses += 1
	}
}

func NextMatch() (pb.Match, error) {
	mu.Lock()
	defer mu.Unlock()

	for i := 0; i < len(tests); i++ {
		if !tests[i].active {
			continue
		}
		pos       := book[tests[i].book - tests[i].book % 2]
		turn      := tests[i].book % 2 == 0
		tests[i].book += 1
		return pb.Match {
			Baseline:  &tests[i].baseline,
			Candidate: &tests[i].candidate,
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
		CheckRefStatus(repo)
		PrintStatus()
		// Sleep to avoid sending too many requests
		time.Sleep(10 * time.Second)
	}
}
