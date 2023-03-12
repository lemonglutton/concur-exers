package main

import (
	"log"
	"math/rand"
	"time"
)

const pulseFreq = 3 * time.Second
const jobSchedulerFreq = 10 * time.Second

func RunHeartBeatExample() {
	const numberOfJobs = 20

	done := make(chan struct{})
	heartbeat, jobWorker := doExpensiveWork(done, numberOfJobs)

	defer log.Printf("Leaving expensiveWorkConsumer...\n")
	go func() {
		defer close(done)
		time.Sleep(30 * time.Second)
		log.Printf("Accident happened, killing process..")

	}()

	for {
		select {
		case <-done:
			return
		case <-heartbeat:
			log.Printf("Jobs process is healthy..")
		case jobId := <-jobWorker:
			log.Printf("Job#%d has been finished..", jobId)
		}
	}

}

func doExpensiveWork(done <-chan struct{}, ids int) (<-chan time.Time, <-chan interface{}) {
	jobs := make(chan interface{})
	heartbeat := make(chan time.Time)

	pulse := time.NewTicker(pulseFreq)
	jobScheduler := time.NewTicker(jobSchedulerFreq)

	go func() {
		defer close(heartbeat)
		defer close(jobs)

		sendPulse := func(t time.Time) {
			select {
			case heartbeat <- t:
			default:
			}
		}

		doTheJob := func(id int) {
			time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
			select {
			case <-done:
				return
			case jobs <- id:
			}
		}

		var cnt int
		for {
			select {
			case <-done:
				return
			case t := <-pulse.C:
				sendPulse(t)
			case <-jobScheduler.C:
				go doTheJob(cnt)
				cnt++
			}
		}
	}()
	return heartbeat, jobs
}
