package main

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestCaching(t *testing.T) {
	cars := []Car{
		{vinNumber: "1", model: "Mondeo"},
		{vinNumber: "2", model: "Citroen"},
		{vinNumber: "3", model: "Audi"},
		{vinNumber: "4", model: "Jaguar"},
		{vinNumber: "5", model: "Porshe"},
		{vinNumber: "6", model: "Ferrari"},
		{vinNumber: "7", model: "Nissan"},
		{vinNumber: "8", model: "Alfa Romeo"},
		{vinNumber: "9", model: "Volvo"},
		{vinNumber: "10", model: "Volkswagen"},
	}

	t.Run("Add records to cache with warmup", func(t *testing.T) {
		const maxCap = 10
		m := NewInMemoryCache(&fifo{}, maxCap,
			[]Car{
				Car{vinNumber: "11", model: "Lamborgini"},
				Car{vinNumber: "12", model: "Tata"},
				Car{vinNumber: "13", model: "BMW"}})

		wait := sync.WaitGroup{}
		ready := sync.WaitGroup{}

		wait.Add(len(cars))
		ready.Add(len(cars))
		for _, entity := range cars {
			go func(c Car) {
				defer wait.Done()

				ready.Done()
				ready.Wait()

				time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
				m.Update(c)
			}(entity)
		}

		go func() {
			ready.Wait()
			<-time.After(time.Duration(40 * time.Second))
			t.Errorf("Test has exceeded 40s, timeout")
		}()
		wait.Wait()

		for _, entity := range cars {
			if _, err := m.Read(entity.Id()); err != nil {
				t.Errorf("Expected : %v, actual: %v", entity.Id(), err.Error())
			}
		}
	})

	t.Run("Add records to cache with small capacity and multiple concurrent readers", func(t *testing.T) {
		const maxCap = 2
		m := NewInMemoryCache(&fifo{}, maxCap, nil)

		wait := sync.WaitGroup{}
		ready := sync.WaitGroup{}
		wait.Add(2 * len(cars))
		ready.Add(len(cars))
		for _, entity := range cars {
			go func(c Car) {
				defer wait.Done()

				ready.Done()
				ready.Wait()

				time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
				m.Update(c)
			}(entity)

			go func(c Car) {
				defer wait.Done()

				ready.Wait()
				m.Read(c.Id())
			}(entity)

		}

		go func() {
			ready.Wait()
			<-time.After(time.Duration(45 * time.Second))
			t.Errorf("Test has exceeded 45s, timeout")
		}()
		wait.Wait()

		var arr []Car
		for _, entity := range cars {
			if car, err := m.Read(entity.Id()); err == nil {
				arr = append(arr, car)
			}
		}

		if len(arr) != maxCap {
			t.Errorf("Expected: %d, actual: %d", maxCap, len(arr))
		}
	})

	t.Run("Modify records if they are already in cache", func(t *testing.T) {
		const maxCap = 13
		warmupCars := []Car{
			{vinNumber: "1", model: "MondeoOrg"},
			{vinNumber: "2", model: "CitroenOrg"},
			{vinNumber: "3", model: "AudiOrg"}}

		m := NewInMemoryCache(&fifo{}, maxCap, warmupCars)

		originals := make(map[string]Car)
		for _, car := range warmupCars {
			carFromCache, err := m.Read(car.Id())
			if err != nil {
				t.Errorf("Failed to prepare data for test")
			}
			originals[carFromCache.Id()] = carFromCache
		}

		wait := sync.WaitGroup{}
		ready := sync.WaitGroup{}
		for _, entity := range cars {
			go func(c Car) {
				defer wait.Done()

				ready.Done()
				ready.Wait()

				time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
				m.Update(c)
			}(entity)
		}
		wait.Wait()

		var counter int
		for _, car := range cars {
			cachedCar, err := m.Read(car.Id())

			if err == nil {
				counter++
			}

			if car, exists := originals[cachedCar.Id()]; exists && car.Model() == cachedCar.Model() {
				t.Errorf("Expected: %v, actual: %v", originals[cachedCar.Id()].Model(), car.Model())
			}
		}

		if counter != maxCap-3 {
			t.Errorf("Expected: %v, actual: %v", maxCap-3, counter)
		}
	})

	t.Run("Too high warump set", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		NewInMemoryCache(&fifo{}, 2, []Car{
			{vinNumber: "11", model: "Lamborgini"},
			{vinNumber: "12", model: "Tata"},
			{vinNumber: "13", model: "BMW"}})
	})

}

func TestPurge(t *testing.T) {}

func TestPrint(t *testing.T) {}

func BenchmarkFifo(b *testing.B) {}

func BenchmarkLru(b *testing.B) {}

func BenchmarkLfu(b *testing.B) {

}
