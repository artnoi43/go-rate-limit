package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/artnoi43/go-rate-limit/config"
	"github.com/artnoi43/go-rate-limit/lib/utils"
)

var (
	conf *config.Config
	f    flags // CLI flags
)

func init() {
	var err error
	conf, err = config.Load()
	if err != nil {
		panic(err)
	}
	f.parse(conf)
	log.Printf(
		"Starting poller\nURL=%s MaxGuard=%d",
		f.URL,
		f.maxGuard,
	)
}

func main() {
	guard := make(chan struct{}, f.maxGuard)
	quit := make(chan struct{})
	statusChan := make(chan int)
	timeChan := make(chan time.Duration)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
		syscall.SIGTERM, // kill -SIGTERM XXXX
	)

	var mut sync.RWMutex
	var counter int
	var getTimes []time.Duration

	client := new(http.Client)
	req, _ := http.NewRequest(http.MethodGet, conf.URL, nil)

	loop := func(u string) {
		start := time.Now()
		resp, _ := client.Do(req)
		<-guard
		fmt.Println(counter)
		mut.Lock()
		counter++
		mut.Unlock()
		timeChan <- time.Since(start)
		statusChan <- resp.StatusCode
	}
	exit := func() {
		log.Println("Shutting down poller")
		utils.CalcAvgTime(getTimes)
		os.Exit(0)
	}

	for {
		guard <- struct{}{}
		go loop(f.URL)
		go func() {
			select {
			case status := <-statusChan:
				switch status {
				case http.StatusTooManyRequests:
					quit <- struct{}{}
				}
			case t := <-timeChan:
				getTimes = append(getTimes, t)
			case <-sigChan:
				exit()
			case <-quit:
				exit()
			}
		}()
	}
}
