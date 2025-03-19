package prog

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// program implements svc.Service
type program struct {
	ctx      context.Context
	wg       errgroup.Group
	quit     chan struct{}
	initFunc func() error
	mainFunc func() error
	stopFunc func()
}

func NewProgram(ctx context.Context, f func() error) *program {
	return &program{mainFunc: f}
}

func (p *program) Init(f func() error) {
	p.initFunc = f
}

func (p *program) Stop(f func()) {
	p.stopFunc = f
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
	if p.initFunc != nil {
		if err := p.initFunc(); err != nil {
			return err
		}
	}

	if err := p.start(); err != nil {
		return err
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-signalChan:
	case <-p.ctx.Done():
	}

	close(p.quit)
	p.wg.Wait()
	return nil
}

func (p *program) Run() {
	if err := p.run(); err != nil {
		log.Fatal(err)
	}
}
