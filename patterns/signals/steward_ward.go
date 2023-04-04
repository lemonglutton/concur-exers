package main

import (
	"log"
	"time"
)

// Sometimes the goroutines are dependent ona resource that we don't have very good control of.
// Maybe goroutine recieves q request to pull data from a web service, or maybe it's monitoring an ephemeral file.
// The point is that it can be very easy for a goroutine to become stuck in a bad state which it cannot recocer without external help.
// Steward's job is to keep checking if hypothetical long-living process is in good condition and if not take apropiate actions

type workerFn func(
	done <-chan struct{},
	pulseInterval time.Duration,
) <-chan struct{}

func RunStewardWardExample() {
	done := make(chan struct{})
	// time after steward rerun ward(worker) if no signal from ward is recieved
	stewardTimeout := time.Duration(10 * time.Second)

	// Frequency of sending Health checks by steward
	stewardPulseRate := time.Duration(10 * time.Second)

	go func() {
		time.Sleep(1 * time.Minute)
		defer close(done)
		log.Printf("Vulcano explosion, terrible accident! Killing program..")
	}()

	for {
		select {
		case <-done:
			return
		case <-RunWithSteward(done, stewardPulseRate, stewardTimeout, worker):
			log.Printf("Steward is alive...")
		}
	}
}

func RunWithSteward(done <-chan struct{}, healthCheckRate time.Duration, timeout time.Duration, worker workerFn) <-chan struct{} {
	stewardHeartbeat := make(chan struct{})
	pulse := time.NewTicker(healthCheckRate)

	go func() {
		defer close(stewardHeartbeat)
		var wardHeartbeat <-chan struct{}
		var wardDone chan struct{}

		startWard := func() {
			wardDone = make(chan struct{})
			wardHeartbeat = worker(orDone(done, wardDone), timeout/2)
		}
		startWard()

		sendStewardPulse := func() {
			select {
			case <-done:
				return
			case stewardHeartbeat <- struct{}{}:
			}
		}

	monitorLoop:
		for {
			timeout := time.After(timeout)
			for {
				select {
				case <-pulse.C:
					sendStewardPulse()
				case <-done:
					return
				case <-wardHeartbeat:
					log.Printf("Ward is healthy..")
					continue monitorLoop

				case <-timeout:
					log.Printf("Ward is unhealthy..- Healing")
					close(wardDone)
					startWard()
					continue monitorLoop
				}
			}
		}
	}()
	return stewardHeartbeat
}

func worker(done <-chan struct{}, pulseInterval time.Duration) <-chan struct{} {
	heartbeat := make(chan struct{})
	pulse := time.NewTicker(pulseInterval)

	sendPulse := func() {
		select {
		case <-done:
			return
		case heartbeat <- struct{}{}:
		}
	}

	go func() {
		defer close(heartbeat)

		for {
			select {
			case <-done:
				return
			case <-pulse.C:
				sendPulse()
			}
		}
	}()

	return heartbeat
}

func orDone(done <-chan struct{}, stream <-chan struct{}) <-chan struct{} {
	orDone := make(chan struct{})
	go func() {
		defer close(orDone)
		for {
			select {
			case <-done:
				return
			case v, ok := <-stream:
				if !ok {
					return
				}

				select {
				case <-done:
					return
				case orDone <- v:
				}
			}
		}
	}()
	return orDone
}
