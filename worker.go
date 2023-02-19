package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

type workerStatus int

const (
	workerStarting workerStatus = iota
	workerReady
	workerBusy
	workerDead
)

type worker struct {
	cmd     *exec.Cmd
	stdin   chan string
	stdout  chan string
	stderr  chan string
	status  workerStatus
	jobs_in *chan job
}

func pipe_reader(r io.Reader) chan string {
	ch := make(chan string, 100)

	go func(r io.Reader, ch chan string) {
		defer close(ch)
		br := bufio.NewReader(r)
		for {
			line, err := br.ReadString('\n')
			if err != nil {
				log.Printf("br.ReadString() failed: %v", err)
				break
			}
			ch <- strings.TrimRight(line, "\n")
		}
	}(r, ch)

	return ch
}

func pipe_writer(w io.Writer) chan string {
	ch := make(chan string, 100)

	go func(w io.Writer, ch chan string) {
		bw := bufio.NewWriter(w)
		for {
			line, ok := <-ch
			if !ok {
				log.Printf("pipe_writer(): error reading from channel, exiting")
				break
			}
			_, err := bw.WriteString(line + "\n")
			if err != nil {
				log.Printf("bw.WriteString() failed: %v", err)
				break
			}
			bw.Flush()
		}
	}(w, ch)

	return ch
}

func newWorker(args []string, jobs_in *chan job) (waitable, error) {
	w := &worker{
		cmd:     exec.Command(args[0], args[1:]...),
		jobs_in: jobs_in,
	}

	stdin, err := w.cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("cmd.StdinPipe() failed: %v", err)
	}

	stdout, err := w.cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("cmd.StdoutPipe() failed: %v", err)
	}

	stderr, err := w.cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("cmd.StderrPipe() failed: %v", err)
	}

	if err := w.cmd.Start(); err != nil {
		return nil, fmt.Errorf("cmd.Start() failed: %v", err)
	}

	w.stdin = pipe_writer(stdin)
	w.stdout = pipe_reader(stdout)
	w.stderr = pipe_reader(stderr)
	w.status = workerStarting
	go w.serve()

	return w.cmd, nil
}

func (w *worker) getStatus() workerStatus {
	return w.status
}

func (w *worker) serve() {
	log.Printf("worker serve() has started")
	defer log.Printf("worker serve() has finished")
	defer close(w.stdin)
	defer func() { w.status = workerDead }()

	waitReady := func() bool {
		for {
			select {
			case line, ok := <-w.stdout:
				if !ok {
					log.Printf("worker starting: failed to read from stdout pipe")
					return false
				}
				log.Printf("worker stdout: %s", line)
				if line == "ready" {
					return true
				}
			case line, ok := <-w.stderr:
				if !ok {
					log.Printf("worker starting: failed to read from stderr pipe")
					return false
				}
				log.Printf("worker stderr: %s", line)
			}
		}
	}

	if !waitReady() {
		return
	}

	for {
		w.status = workerReady

		job := <-*w.jobs_in
		w.stdin <- job.request

		w.status = workerBusy

		response_line, err := func() (string, error) {
			for {
				select {
				case line, ok := <-w.stderr:
					if !ok {
						return "", fmt.Errorf("failed to read from stderr, process terminated unexpectedly")
					}
					log.Printf("running stderr: %s", line)
				case response_line, ok := <-w.stdout:
					if !ok {
						return "", fmt.Errorf("failed to read from stdout, process terminated unexpectedly")
					}
					return response_line, nil
				}
			}
		}()

		if err != nil {
			job.responseCh <- err
			return
		} else {
			job.responseCh <- response_line
		}
	}
}
