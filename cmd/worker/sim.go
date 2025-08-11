package main

import (
	"log"
	pb "detf/api"
)

func Sim(match pb.Match) (pb.Result, error) {
	baseline, err  := GetEngine(*match.GetBaseline())
	candidate, err := GetEngine(*match.GetCandidate())
	if err != nil {
		log.Fatalf("Failed to get engine: %v", err)
	}
	log.Printf("%s %s", baseline, candidate)
	return pb.Result {
		Baseline: match.GetBaseline(),
		Candidate:  match.GetCandidate(),
		Win:  true,
		Draw: false,
	}, nil
}
