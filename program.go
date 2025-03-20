package prog

import (
	"context"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// program implements svc.Service
type program struct {
	ctx       context.Context
	wg        errgroup.Group
	quit      chan struct{}
	initFuncs []func() error
	mainFunc  func() error
	stopFunc  func()
}

type Option func(*program)

// NewProgram returns a new program
func NewProgram(ctx context.Context, mainFunc func() error, options ...Option) *program {
	p := &program{ctx: ctx, mainFunc: mainFunc}

	for _, option := range options {
		option(p)
	}
	return p
}

func WithInit(f func() error) Option {
	return func(p *program) {
		p.initFuncs = append(p.initFuncs, f)
	}
}

func WithStopFunc(f func()) Option {
	return func(p *program) {
		p.stopFunc = f
	}
}

func WithPprof() Option {
	return func(p *program) {
		p.initFuncs = append(p.initFuncs, func() error {
			var w sync.WaitGroup
			w.Add(1)
			go func() {
				server := &http.Server{Addr: ":0", Handler: nil}
				ln, err := net.Listen("tcp", ":0")
				if err != nil {
					w.Done()
					log.Println(err)
					return
				}
				port := ln.Addr().(*net.TCPAddr).Port
				log.Println("pprof listen on port:", port)
				w.Done()
				//
				server.Serve(ln)
			}()
			w.Wait()
			return nil
		})
	}
}

func (p *program) start() error {
	p.quit = make(chan struct{})

	p.wg.Go(func() error {
		p.mainFunc()
		<-p.quit
		return nil
	})

	return nil
}

func (p *program) run() error {
	for _, f := range p.initFuncs {
		if err := f(); err != nil {
			return err
		}
	}

	if err := p.start(); err != nil {
		return err
	}

	select {
	case <-signalChan():
	case <-p.ctx.Done():
	}

	close(p.quit)
	p.wg.Wait()
	if p.stopFunc != nil {
		p.stopFunc()
	}
	return nil
}

func signalChan() <-chan os.Signal {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	return signalChan
}

func (p *program) Run() {
	if err := p.run(); err != nil {
		log.Fatal(err)
	}
}
