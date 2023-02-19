package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"
)

type JobServer struct {
	queue chan job
}

type job struct {
	request    string
	responseCh chan any
}

type waitable interface {
	Wait() error
}

func NewJobServer(workers int) *JobServer {
	jobServer := &JobServer{
		queue: make(chan job, 1000),
	}
	args := []string{
		"/usr/bin/env",
		"python3.11",
		"-m",
		"distribute_challenge.serve",
	}
	for i := 0; i < workers; i++ {
		go func() {
			var process waitable
			for {
				if process != nil {
					log.Printf("worker process exited, err=%v", process.Wait())
					// Do not respawn too fast
					time.Sleep(200 * time.Millisecond)
				}
				var err error
				process, err = newWorker(args, &jobServer.queue)
				if err != nil {
					log.Fatalf("newWorker() failed: %v", err)
				}
			}
		}()
	}

	return jobServer
}

func (jobServer *JobServer) Execute(request []byte) ([]byte, error) {
	j := job{
		request:    base64.StdEncoding.EncodeToString(request),
		responseCh: make(chan any, 1),
	}
	jobServer.queue <- j
	switch response := (<-j.responseCh).(type) {
	case string:
		decoded, err := base64.StdEncoding.DecodeString(response)
		if err != nil {
			return nil, fmt.Errorf("Can't decode base64 response from process: %v", err)
		}
		return decoded, nil
	default:
		return nil, response.(error)
	}
}
