package main

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestCaching(t *testing.T) {
	cars := []Entity{
		Car{vinNumber: "1", model: "Mondeo"},
		Car{vinNumber: "2", model: "Citroen"},
		Car{vinNumber: "3", model: "Audi"},
		Car{vinNumber: "4", model: "Jaguar"},
		Car{vinNumber: "5", model: "Porshe"},
		Car{vinNumber: "6", model: "Ferrari"},
		Car{vinNumber: "7", model: "Nissan"},
		Car{vinNumber: "8", model: "Alfa Romeo"},
		Car{vinNumber: "9", model: "Volvo"},
		Car{vinNumber: "10", model: "Volkswagen"},
	}

	t.Run("Add records to cache with warmup", func(t *testing.T) {
		const maxCap = 10
		m := NewInMemoryCache(&fifo{}, maxCap,
			[]Entity{
				Car{vinNumber: "11", model: "Lamborgini"},
				Car{vinNumber: "12", model: "Tata"},
				Car{vinNumber: "13", model: "BMW"}})

		wait := sync.WaitGroup{}
		ready := sync.WaitGroup{}

		wait.Add(len(cars))
		ready.Add(len(cars))
		for _, entity := range cars {
			go func(e Entity) {
				defer wait.Done()

				ready.Done()
				ready.Wait()

				time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
				m.Update(e)
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
			go func(e Entity) {
				defer wait.Done()

				ready.Done()
				ready.Wait()

				time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
				m.Update(e)
			}(entity)

			go func(e Entity) {
				defer wait.Done()

				ready.Wait()
				m.Read(e.Id())
			}(entity)

		}

		go func() {
			ready.Wait()
			<-time.After(time.Duration(45 * time.Second))
			t.Errorf("Test has exceeded 45s, timeout")
		}()
		wait.Wait()

		var arr []Entity
		for _, entity := range cars {
			if car, _ := m.Read(entity.Id()); car != nil {
				arr = append(arr, car)
			}
		}

		if len(arr) != maxCap {
			t.Errorf("Expected: %d, actual: %d", maxCap, len(arr))
		}
	})

	t.Run("Modify records if they are already in cache", func(t *testing.T) {
		const maxCap = 13
		warmupCars := []Entity{
			Car{vinNumber: "1", model: "MondeoOrg"},
			Car{vinNumber: "2", model: "CitroenOrg"},
			Car{vinNumber: "3", model: "AudiOrg"}}

		m := NewInMemoryCache(&fifo{}, maxCap, warmupCars)

		var originals []Entity
		for _, car := range warmupCars {
			entity, err := m.Read(car.Id())
			if err != nil {
				t.Errorf("Failed to prepare data for test")
			}
			originals = append(originals, entity)
		}

		wait := sync.WaitGroup{}
		ready := sync.WaitGroup{}
		for _, entity := range cars {
			go func(e Entity) {
				defer wait.Done()

				ready.Done()
				ready.Wait()

				time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
				m.Update(e)
			}(entity)
		}
		wait.Wait()

	})

	t.Run("Too high warump set", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		NewInMemoryCache(&fifo{}, 2, []Entity{
			Car{vinNumber: "11", model: "Lamborgini"},
			Car{vinNumber: "12", model: "Tata"},
			Car{vinNumber: "13", model: "BMW"}})
	})

}

func TestPurge(t *testing.T) {}

func TestPrint(t *testing.T) {}

func BenchmarkFifo(b *testing.B) {}

func BenchmarkLru(b *testing.B) {}

func BenchmarkLfu(b *testing.B) {

}
