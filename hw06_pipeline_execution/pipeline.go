package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		out = chainChannel(out, done)

		out = stage(out)
	}

	return out
}

func chainChannel(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for v := range in {
			select {
			case <-done:
				return
			default:
			}

			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}()

	return out
}
