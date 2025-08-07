package main

import (
	"log"
	"net"
	"sync"
	"errors"
	"io"
	"bufio"
	"os"
	"flag"
	"fmt"

	"google.golang.org/grpc"
	pb "detf/api"
)

type test struct {
	active bool
	id     uint64
	wins   uint64
	losses uint64
	draws  uint64
	book   uint64
}

type match struct {
	id  uint64
	pos string
}

type result struct {
	id   uint64
	win  bool
	draw bool
}

type server struct {
	pb.UnimplementedDETFServer
}

var mtx   sync.Mutex
var tests []test
var book  []string

func HandleResult(res result) {
	mtx.Lock()
	defer mtx.Unlock()
	if res.draw {
		tests[res.id].draws += 1
	} else {
		if res.win {
			tests[res.id].wins += 1
		} else {
			tests[res.id].losses += 1
		}
	}
}

func NextMatch() (match, error) {
	mtx.Lock()
	defer mtx.Unlock()
	for _, test := range tests {
		if !test.active {
			continue
		}
		test.book += 1
		return match {
			id:  test.id,
			pos: book[test.book],
		}, nil
	}
	return match {}, 
	       errors.New("no tests currently active")
}

func SendMatch(
	stream pb.DETF_StreamServer,
) error {
	initial_match, err := NextMatch()
	if err != nil {
		return err
	}
	stream.Send(&pb.Match {
		ID: initial_match.id,
	})
	return nil
}

func (s *server) Stream(
	stream pb.DETF_StreamServer,
) error {
	{
		err := SendMatch(stream)
		if err != nil {
			return err
		}
	}
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		HandleResult(
			result {
				id: in.ID,
			},
		)
		{
			err := SendMatch(stream)
			if err != nil {
				return err
			}
		}
	}
}

func LoadBook(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
	        book = append(book, scanner.Text())
	}
}

func main() {
	book_path := flag.String("b", "", "path to opening book")
	port := flag.Int("p", 8080, "port to operate on")
	flag.Parse()
	LoadBook(*book_path)
	tests = append(tests, test {
		active: true,
		id:     0,
		wins:   0,
		losses: 0,
		draws:  0,
		book:   0,
	})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("%v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDETFServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr());
	if err := s.Serve(lis); err != nil {
		log.Fatalf("%v", err)
	}
}
