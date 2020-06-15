package gameloop

import (
	"fmt"
	"time"
)

// StartLoop is physics game loop based on https://gafferongames.com/post/fix_your_timestep/
func StartLoop(handler GameHandler, updateRate time.Duration, startingState interface{}) chan<- struct{} {

	var t float64
	deltaTime := updateRate.Seconds()

	currentTime := time.Now()
	var accumulator float64

	var (
		previous interface{}
		current  interface{}
	)
	current = startingState

	stopChann := make(chan struct{}, 0)
	maxFrameTime := time.Duration(time.Second / 4).Seconds()

	go func() {
		for {
			select {
			case <-stopChann:
				fmt.Println("stopping loop")
				return
			default:
				newTime := time.Now()
				frameTime := newTime.Sub(currentTime).Seconds()
				if frameTime > maxFrameTime {
					frameTime = maxFrameTime
				}

				currentTime = newTime
				accumulator += frameTime

				for accumulator >= deltaTime {
					previous = current
					current = handler.Integrate(current, t, deltaTime)
					t += deltaTime
					accumulator -= deltaTime
				}

				alpha := accumulator / deltaTime

				handler.Render(current, t, 1.0-alpha)
			}
		}
	}()
	return stopChann
}
