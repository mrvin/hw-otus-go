package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outs := make([]Out, len(stages))
	out := make(Bi)
	finish := make(chan struct{})

	for i, stage := range stages {
		outs[i] = stage(in)
		in = outs[i]
	}

	go func() {
		select {
		case <-done:
			for _, out := range outs {
				go emptyReading(out)
			}
		case <-finish:
			return
		}
	}()

	go func() {
	loop:
		for {
			select {
			case res, ok := <-outs[len(outs)-1]:
				if !ok {
					close(finish)
					break loop
				}
				out <- res
			case <-done:
				break loop
			}
		}
		close(out)
	}()

	return out
}

func emptyReading(out Out) {
	for range out { //nolint:revive
		// do nothing.
	}
}
