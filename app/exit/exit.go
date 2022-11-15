package exit

import (
	"os"
	"os/signal"
)

func OnExit(onExit ...func()) {
	ctrlC := make(chan os.Signal)
	signal.Notify(ctrlC, os.Interrupt)

	go func() {
		for {
			select {
			case <-ctrlC:

				for _, exit := range onExit {
					exit()
				}

				os.Exit(0)
			}
		}
	}()
}
