package main

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestCaching(t *testing.T) {
	cars := []Car{
		{vinNumber: 1, model: "Mondeo"},
		{vinNumber: 2, model: "Citroen"},
		{vinNumber: 3, model: "Audi"},
		{vinNumber: 4, model: "Jaguar"},
		{vinNumber: 5, model: "Porshe"},
		{vinNumber: 6, model: "Ferrari"},
		{vinNumber: 7, model: "Nissan"},
		{vinNumber: 8, model: "Alfa Romeo"},
		{vinNumber: 9, model: "Volvo"},
		{vinNumber: 10, model: "Volkswagen"},
	}

	t.Run("Add records to cache with concurrent writers and with warmup set up", func(t *testing.T) {
		t.Parallel()

		const maxCap = 10
		m := NewInMemoryCache(&fifo{}, maxCap,
			[]Car{
				{vinNumber: 11, model: "Lamborgini"},
				{vinNumber: 12, model: "Tata"},
				{vinNumber: 13, model: "BMW"}})

		wait := sync.WaitGroup{}
		ready := sync.WaitGroup{}

		wait.Add(len(cars))
		ready.Add(len(cars))
		for _, entity := range cars {
			go func(c Car) {
				defer wait.Done()

				ready.Done()
				ready.Wait()

				m.Update(c)
			}(entity)
		}

		go func() {
			ready.Wait()
			<-time.After(time.Duration(10 * time.Second))
			t.Errorf("Test has exceeded 10s, timeout")
		}()
		wait.Wait()

		for _, entity := range cars {
			if val, err := m.Read(entity.Id()); err != nil {
				t.Errorf("Expected : %v, actual: %v", entity.Id(), err.Error())
			} else {
				t.Logf("Car: %v", val)
			}
		}
	})

	t.Run("Add records to cache with concurrent writers, small capacity and multiple concurrent readers", func(t *testing.T) {
		t.Parallel()

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
			<-time.After(time.Duration(10 * time.Second))
			t.Errorf("Test has exceeded 10s, timeout")
		}()
		wait.Wait()

		var arr []Car
		for _, entity := range cars {
			if car, err := m.Read(entity.Id()); err == nil {
				arr = append(arr, car)
			}
		}

		t.Logf("Records in cache after operation: %v", arr)
		if len(arr) != maxCap {
			t.Errorf("Expected: %d, actual: %d", maxCap, len(arr))
		}
	})

	t.Run("Add records to cache with concurrent writers, where some of the records needs to be only modified", func(t *testing.T) {
		t.Parallel()

		const maxCap = 13
		warmupCars := []Car{
			{vinNumber: 1, model: "MondeoOrg"},
			{vinNumber: 2, model: "CitroenOrg"},
			{vinNumber: 3, model: "AudiOrg"}}

		m := NewInMemoryCache(&fifo{}, maxCap, warmupCars)

		originals := make(map[int]Car)
		for _, car := range warmupCars {
			carFromCache, err := m.Read(car.Id())
			if err != nil {
				t.Errorf("Failed to prepare data for test")
			}
			originals[carFromCache.Id()] = carFromCache
		}

		wait := sync.WaitGroup{}
		ready := sync.WaitGroup{}
		wait.Add(len(cars))
		ready.Add(len(cars))
		for _, entity := range cars {
			go func(c Car) {
				defer wait.Done()

				ready.Done()
				ready.Wait()

				m.Update(c)
			}(entity)
		}
		wait.Wait()

		go func() {
			ready.Wait()
			<-time.After(time.Duration(10 * time.Second))
			t.Errorf("Test has exceeded 10s, timeout")
		}()

		var arr []Car
		for _, car := range cars {
			cachedCar, err := m.Read(car.Id())

			if err == nil {
				arr = append(arr, cachedCar)
			}

			if car, exists := originals[cachedCar.Id()]; exists && car.Model() == cachedCar.Model() {
				t.Errorf("Expected: %v, actual: %v", originals[cachedCar.Id()].Model(), car.Model())
			}
		}

		t.Logf("Records in cache before: %v, after operation: %v\n", originals, arr)
		if len(arr) != maxCap-3 {
			t.Errorf("Expected: %v, actual: %v", maxCap-3, len(arr))
		}
	})

	t.Run("Initialize cache with warmup that exceeds cache capacity", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		NewInMemoryCache(&fifo{}, 2, []Car{
			{vinNumber: 11, model: "Lamborgini"},
			{vinNumber: 12, model: "Tata"},
			{vinNumber: 13, model: "BMW"}})
	})
}

func TestPurge(t *testing.T) {
	t.Parallel()

	cars := []Car{
		{vinNumber: 1, model: "Mondeo"},
		{vinNumber: 2, model: "Citroen"},
		{vinNumber: 3, model: "Audi"},
		{vinNumber: 4, model: "Jaguar"},
		{vinNumber: 5, model: "Porshe"},
		{vinNumber: 6, model: "Ferrari"},
		{vinNumber: 7, model: "Nissan"},
		{vinNumber: 8, model: "Alfa Romeo"},
		{vinNumber: 9, model: "Volvo"},
		{vinNumber: 10, model: "Volkswagen"},
	}

	m := NewInMemoryCache(&fifo{}, 10, nil)

	main := sync.WaitGroup{}
	ready := sync.WaitGroup{}

	main.Add(2*len(cars) + 1)
	ready.Add(len(cars))
	for _, car := range cars {
		go func(c Car) {
			defer main.Done()

			ready.Done()
			ready.Wait()
			m.Update(c)
		}(car)

		go func(c Car) {
			defer main.Done()

			ready.Wait()
			m.Read(c.Id())
		}(car)

	}

	go func() {
		defer main.Done()
		ready.Wait()
		m.Purge()
	}()
	main.Wait()

	var carCol []Car
	for _, car := range cars {
		cachedCar, _ := m.Read(car.Id())

		if cachedCar != (Car{}) {
			carCol = append(carCol, cachedCar)
		}
	}

	t.Logf("Records in cache after Purge: %v", carCol)
	if len(carCol) == 10 {
		t.Errorf("Expected: <10, actual: %d", len(carCol))
	}
}

//TODO reduce number of allocations in test and rewrite tests to use RunParallel
func BenchmarkReads(b *testing.B) {
	var table = []struct {
		cacheCapacity             int
		numberOfConcurrentReaders int
	}{
		{cacheCapacity: 30, numberOfConcurrentReaders: 100},
		// {cacheCapacity: 300, numberOfConcurrentReaders: 1000},
		// {cacheCapacity: 2400, numberOfConcurrentReaders: 8000},
	}

	for _, v := range table {
		b.Run(fmt.Sprintf("cacheCapacity: %d, numberOfConcurrentReaders: %d", v.cacheCapacity, v.numberOfConcurrentReaders), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				main := sync.WaitGroup{}
				ready := sync.WaitGroup{}
				signal := make(chan struct{})

				main.Add(v.numberOfConcurrentReaders + 1)
				ready.Add(v.numberOfConcurrentReaders)

				m := NewInMemoryCache(&fifo{}, v.cacheCapacity, nil)

				// This for brings more alocations
				for j := 1; j <= v.cacheCapacity; j++ {
					car := Car{vinNumber: i, model: fmt.Sprintf("CarModel_%d", j)}
					m.Update(car)
				}

				for i := 0; i < v.numberOfConcurrentReaders; i++ {
					go func(id int) {
						defer main.Done()
						ready.Done()
						ready.Wait()
						<-signal
						m.Read(rand.Intn(v.cacheCapacity) + 1)
					}(i)
				}

				go func() {
					defer main.Done()
					ready.Wait()
					b.StartTimer()
					close(signal)
				}()
				main.Wait()
			}
		})
	}
}

func BenchmarkWrites(b *testing.B) {
	var table = []struct {
		cacheCapacity            int
		numberOfConcurrentWrites int
	}{
		{cacheCapacity: 10, numberOfConcurrentWrites: 20},
		{cacheCapacity: 300, numberOfConcurrentWrites: 1000},
		// {cacheCapacity: 2400, numberOfConcurrentWrites: 8000},
	}

	for _, v := range table {
		b.Run(fmt.Sprintf("cacheCapacity: %d, numberOfConcurrentWrites: %d", v.cacheCapacity, v.numberOfConcurrentWrites), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				main := sync.WaitGroup{}
				ready := sync.WaitGroup{}
				signal := make(chan struct{})

				ready.Add(v.numberOfConcurrentWrites)
				main.Add(v.numberOfConcurrentWrites + 1)

				m := NewInMemoryCache(&fifo{}, v.cacheCapacity, nil)
				for j := 0; j < v.numberOfConcurrentWrites; j++ {
					car := Car{vinNumber: i, model: fmt.Sprintf("CarModel_%d", j)}
					go func(c Car) {
						defer main.Done()
						ready.Done()
						ready.Wait()

						<-signal
						m.Update(c)
					}(car)
				}

				go func() {
					defer main.Done()
					ready.Wait()
					b.StartTimer()
					close(signal)
				}()
				main.Wait()
			}
		})
	}
}
