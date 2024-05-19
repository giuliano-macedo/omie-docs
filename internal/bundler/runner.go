package bundler

import "sync"

type bundlersWorker struct {
	args       Args
	waitGroup  *sync.WaitGroup
	errChannel chan error
	bundler    Bundler
}

type Runner struct {
	Args Args
}

func (worker *bundlersWorker) run() {
	worker.errChannel <- worker.bundler.Bundle(worker.args)
	worker.waitGroup.Done()
}

func (runner *Runner) Run(bundlers ...Bundler) error {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(bundlers))
	errChannel := make(chan error, len(bundlers))

	for _, bundler := range bundlers {
		worker := bundlersWorker{
			args:       runner.Args,
			waitGroup:  &waitGroup,
			errChannel: errChannel,
			bundler:    bundler,
		}
		go worker.run()
	}
	waitGroup.Wait()

	for i := 0; i < len(bundlers); i++ {
		bundlerErr := <-errChannel
		if bundlerErr != nil {
			return bundlerErr
		}
	}
	return nil
}
