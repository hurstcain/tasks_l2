package main

import (
	"github.com/hurstcain/tasks_l2/develop/dev11/internal/service"
	"log"
	"os"
	"os/signal"
)

func main() {
	s := service.NewService()

	chExit := make(chan os.Signal, 1)
	signal.Notify(chExit, os.Interrupt)

	go func() {
		select {
		case <-chExit:
			log.Println("Closing server...")
			s.Stop()
		}
	}()

	s.Run()
}
