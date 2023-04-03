package main

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

func TestCaching(t *testing.T) {
	entities := []Entity{
		Car{vinNumber: 1, model: "Mondeo"},
		Car{vinNumber: 2, model: "Citroen"},
		Car{vinNumber: 3, model: "Audi"},
		Car{vinNumber: 4, model: "Jaguar"},
		Car{vinNumber: 5, model: "Porshe"},
		Car{vinNumber: 6, model: "Ferrari"},
		Car{vinNumber: 7, model: "Nissan"},
		Car{vinNumber: 8, model: "Alfa Romeo"},
		Car{vinNumber: 9, model: "Volvo"},
		Car{vinNumber: 10, model: "Volkswagen"},
	}

	t.Run("Check if cache evistion works properly. This test should remove the warmup dataset, since its the oldest data stored and complete cache's capacity to its maximum ", func(t *testing.T) {
		t.Parallel()

		const maxCap = 10
		m := NewInMemoryCache(&fifo{}, maxCap,
			[]Entity{
				Car{vinNumber: 11, model: "Lamborgini"},
				Car{vinNumber: 12, model: "Tata"},
				Car{vinNumber: 13, model: "BMW"}})

		wait := sync.WaitGroup{}
		wait.Add(len(entities))

		for _, entity := range entities {
			go func(e Entity) {
				defer wait.Done()
				m.Set(e)

			}(entity)
		}
		wait.Wait()

		for _, entity := range entities {
			if val, err := m.Read(entity.Id()); err != nil {
				t.Errorf("Expected : %v, actual: %v", entity.Id(), err.Error())
			} else {
				t.Logf("Car: %v", val)
			}
		}
	})

	t.Run("Check if cache eviction works properly for small capacities and if after all reads and updates cache max capacity is not crossed", func(t *testing.T) {
		t.Parallel()
		const maxCap = 2
		m := NewInMemoryCache(&fifo{}, maxCap, nil)

		wait := sync.WaitGroup{}
		wait.Add(2 * len(entities))
		for _, entity := range entities {
			go func(e Entity) {
				defer wait.Done()
				m.Set(e)
			}(entity)

			go func(e Entity) {
				defer wait.Done()
				m.Read(e.Id())
			}(entity)

		}
		wait.Wait()

		var arr []Entity
		for _, entity := range entities {
			if cachedEntity, err := m.Read(entity.Id()); err == nil {
				arr = append(arr, cachedEntity)
			}
		}

		t.Logf("Records in cache after operation: %v", arr)
		if len(arr) != maxCap {
			t.Errorf("Expected: %d, actual: %d", maxCap, len(arr))
		}
	})

	t.Run("Check if cache gonna add additional objects alongside warmupset without eviction and if already existed warmup records gonna be modified correctly", func(t *testing.T) {
		t.Parallel()

		const maxCap = 13
		warmupEntities := []Entity{
			Car{vinNumber: 1, model: "MondeoOrg"},
			Car{vinNumber: 2, model: "CitroenOrg"},
			Car{vinNumber: 3, model: "AudiOrg"}}

		m := NewInMemoryCache(&fifo{}, maxCap, warmupEntities)
		originals := make(map[int]Car)
		for _, entity := range warmupEntities {
			cachedEntity, err := m.Read(entity.Id())
			if err != nil {
				t.Errorf("Failed to prepare data for a test")
			}

			cachedCar := cachedEntity.(Car)
			originals[cachedEntity.Id()] = cachedCar
		}

		wait := sync.WaitGroup{}
		wait.Add(len(entities))
		for _, entity := range entities {
			go func(e Entity) {
				defer wait.Done()
				m.Set(e)
			}(entity)
		}
		wait.Wait()

		var cars []Car
		for _, entity := range entities {
			cachedEntity, err := m.Read(entity.Id())
			cachedCar := cachedEntity.(Car)

			if err == nil {
				cars = append(cars, cachedCar)
			}

			if car, exists := originals[cachedEntity.Id()]; exists && car.Model() == cachedCar.Model() {
				t.Errorf("Expected: %v, actual: %v", originals[cachedEntity.Id()].Model(), car.Model())
			}
		}

		t.Logf("Records in cache before: %v, after operation: %v\n", originals, cars)
		if len(cars) != maxCap-3 {
			t.Errorf("Expected: %v, actual: %v", maxCap-3, len(cars))
		}
	})

	t.Run("Initialize cache with warmup that exceeds cache capacity. Cache should react with panic", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		NewInMemoryCache(&fifo{}, 2, []Entity{
			Car{vinNumber: 11, model: "Lamborgini"},
			Car{vinNumber: 12, model: "Tata"},
			Car{vinNumber: 13, model: "BMW"}})
	})
}

func TestPurge(t *testing.T) {
	t.Run("Check if Purge function when invoked doesn't disturb whole process of adding and reading records to/from cache", func(t *testing.T) {
		t.Parallel()

		entities := []Entity{
			Car{vinNumber: 1, model: "Mondeo"},
			Car{vinNumber: 2, model: "Citroen"},
			Car{vinNumber: 3, model: "Audi"},
			Car{vinNumber: 4, model: "Jaguar"},
			Car{vinNumber: 5, model: "Porshe"},
			Car{vinNumber: 6, model: "Ferrari"},
			Car{vinNumber: 7, model: "Nissan"},
			Car{vinNumber: 8, model: "Alfa Romeo"},
			Car{vinNumber: 9, model: "Volvo"},
			Car{vinNumber: 10, model: "Volkswagen"},
		}
		m := NewInMemoryCache(&fifo{}, 10, nil)

		main := sync.WaitGroup{}
		ready := sync.WaitGroup{}

		main.Add(2*len(entities) + 1)
		ready.Add(len(entities))
		for _, entity := range entities {
			go func(e Entity) {
				defer main.Done()
				ready.Done()
				ready.Wait()
				m.Set(e)
			}(entity)

			go func(e Entity) {
				defer main.Done()
				ready.Wait()
				m.Read(e.Id())
			}(entity)

		}

		go func() {
			defer main.Done()
			ready.Wait()
			m.Purge()
		}()
		main.Wait()

		var cars []Car
		for _, car := range entities {
			cachedEntity, _ := m.Read(car.Id())
			cachedCar := cachedEntity.(Car)

			if cachedCar != (Car{}) {
				cars = append(cars, cachedCar)
			}
		}

		t.Logf("Records in cache after Purge: %v", cars)
		if len(cars) == 10 {
			t.Errorf("Expected: <10, actual: %d", len(cars))
		}
	})
}

// //TODO reduce number of allocations in test and rewrite tests to use RunParallel
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
					m.Set(car)
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
						m.Set(c)
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
