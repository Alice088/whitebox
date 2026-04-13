package status

import (
	"fmt"
	"os"
	"time"
)

type AnimationController struct {
	statusGen   *StatusEngine
	isAnimating bool
	stopChan    chan struct{}
	doneChan    chan struct{}
}

func NewAnimationController(statusGen *StatusEngine) *AnimationController {
	return &AnimationController{
		statusGen: statusGen,
		stopChan:  make(chan struct{}),
		doneChan:  make(chan struct{}),
	}
}

func (ac *AnimationController) Start() {
	if ac.isAnimating {
		return
	}

	ac.isAnimating = true

	go func() {
		defer func() {
			ac.isAnimating = false
			close(ac.doneChan)
		}()

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				animatedStatus := ac.statusGen.NextAnimated()
				_, _ = fmt.Fprintf(os.Stdout, "\r%s", animatedStatus)
			case <-ac.stopChan:
				_, _ = fmt.Fprint(os.Stdout, "\r\033[K")
				return
			}
		}
	}()
}

func (ac *AnimationController) Stop() {
	if !ac.isAnimating {
		return
	}

	close(ac.stopChan)
	<-ac.doneChan
}
