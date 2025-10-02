package main

import (
	"time"
	pb "detf/api"

	"github.com/corentings/chess/v2"
	"github.com/corentings/chess/v2/uci"
)

func Sim(match pb.Match) (pb.Result, error) {
	baseline, err := InitEngine(*match.GetBaseline())
	if err != nil {
		return pb.Result {}, err
	}
	defer baseline.Close()

	candidate, err := InitEngine(*match.GetCandidate())
	if err != nil {
		return pb.Result {}, err
	}
	defer candidate.Close()

	fen, err := chess.FEN(match.GetPos())
	if err != nil {
		return pb.Result {}, err
	}
	game := chess.NewGame(fen)

	init := game.Position().Turn()
	for game.Outcome() == chess.NoOutcome {
		cPos := uci.CmdPosition { Position: game.Position() }
			cGo  := uci.CmdGo { MoveTime: time.Second }

		var move chess.Move
		if match.GetTurn() && game.Position().Turn() == init {
			if err := baseline.Run(cPos, cGo); err != nil {
				return pb.Result {}, err
			}
			move = *baseline.SearchResults().BestMove
		} else {
			if err := candidate.Run(cPos, cGo); err != nil {
				return pb.Result {}, err
			}
			move = *candidate.SearchResults().BestMove
		}	

		if err := game.Move(&move, nil); err != nil {
			return pb.Result {}, err
		}
	}
	
	draw := game.Outcome() == chess.Draw
	iwin := (init == chess.White && game.Outcome() == chess.WhiteWon) ||
	        (init == chess.Black && game.Outcome() == chess.BlackWon)
	win  := (match.GetTurn() && iwin) || (!match.GetTurn() && !iwin)

	return pb.Result {
		Baseline:  match.GetBaseline(),
		Candidate: match.GetCandidate(),
		Draw:      draw,
		Win:       win,
	}, nil
}

func InitEngine(baseline pb.Engine) (*uci.Engine, error) {
	str, err := GetEngine(baseline)
	if err != nil {
		return nil, err
	}
	return uci.New(str)
}
