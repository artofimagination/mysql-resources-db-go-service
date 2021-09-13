package initialization

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/proemergotech/log/v3"
)

type startFn func(chan<- error)
type stopFn func(time.Duration) error

type runner struct {
	stoppers []stopper
	errorChs map[string]chan error
}

type stopper struct {
	name string
	stop stopFn
}

func newRunner() *runner {
	return &runner{
		stoppers: make([]stopper, 0),
		errorChs: make(map[string]chan error),
	}
}

func (r *runner) start(name string, start startFn, stop stopFn) {
	errorCh := make(chan error)
	r.errorChs[name] = errorCh

	r.stoppers = append(r.stoppers, stopper{
		name: name,
		stop: stop,
	})

	start(errorCh)
	fmt.Println("-------------------------------------------------------------")
	log.Info(context.Background(), name+" started")
	fmt.Println("-------------------------------------------------------------")
}

func (r *runner) stop() {
	for i := len(r.stoppers) - 1; i >= 0; i-- {
		stopper := r.stoppers[i]

		if err := stopper.stop(5 * time.Second); err != nil {
			err = errors.Wrap(err, stopper.name+" graceful shutdown failed")
			log.Error(context.Background(), err.Error(), "error", err)
		}
		log.Info(context.Background(), stopper.name+" shutdown complete")
	}
}

func (r *runner) errors() <-chan error {
	errCollector := make(chan error, len(r.errorChs))
	for name, errorCh := range r.errorChs {
		name := name
		errorCh := errorCh
		go func() {
			if err := <-errorCh; err != nil {
				errCollector <- errors.Wrap(err, name+" fatal error")
			}
		}()
	}
	return errCollector
}
