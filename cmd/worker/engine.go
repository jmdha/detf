package main

import (
	"log"
	"fmt"
	"sync"
	"path"
	"os"
	"os/exec"
	"errors"
	pb "detf/api"
)

type Engine struct {
	repo string
	ref  string
	path string
}

var mu      sync.Mutex
var engines []Engine

func Clone(repo string, ref string, dir string) (string, error) {
	var cmd *exec.Cmd
	cmd = exec.Command("git", "clone", repo)
	cmd.Dir = dir
	if out, err := cmd.Output(); err != nil {
		return "", errors.New(string(out))
	}
	dirs, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	name := dirs[0].Name()
	cmd = exec.Command("git", "checkout", ref)
	cmd.Dir = path.Join(dir, name)
	if _, err := cmd.Output(); err != nil {
		return "", errors.New(fmt.Sprintf("Could not checkout %s - %s in dir %s", repo, ref, path.Join(dir, name)))
	}
	cmd = exec.Command("make")
	cmd.Dir = path.Join(dir, name)
	if _, err := cmd.Output(); err != nil {
		return "", errors.New("Could not build")
	}
	return path.Join(path.Join(dir, name), name), nil
}

func PrepareEngine(repo string, ref string) (Engine, error) {
	log.Printf("Preparing: %s - %s", repo, ref)
	dir, err := os.MkdirTemp("", "DETF")
	if err != nil {
		return Engine {}, err
	}
	log.Printf("Temp dir: %s", dir)
	path, err := Clone(repo, ref, dir); 
	if err != nil {
		return Engine {}, err
	}
	return Engine {
		repo: repo,
		ref:  ref,
		path: path,
	}, nil
}

func GetOrPrepareEngine(repo string, ref string) (Engine, error) {
	// If engine already downloaded return it
	for _, engine := range engines {
		if engine.repo == repo && engine.ref == ref {
			return engine, nil
		}
	}
	// If not, prepare it
	engine, err := PrepareEngine(repo, ref)
	if err != nil {
		return Engine {}, err
	}
	engines = append(engines, engine)
	return engine, nil
}

func GetEngine(pb_engine pb.Engine) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	
	engine, err := GetOrPrepareEngine(
		pb_engine.GetRepo(),
		pb_engine.GetRef(),
	)
	if err != nil {
		return "", err
	} else {
		return engine.path, nil
	}
}
