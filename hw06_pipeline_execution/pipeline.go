package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = chainChannel(in, done)

		in = stage(in)
	}

	return in
}

func chainChannel(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for v := range in {
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}()

	return out
}
