package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {

	Reading()
}

func Reading() {
	m := NewInMemoryCache(&fifo{}, 10,
		[]Entity{
			Car{vinNumber: "1", model: "Mondeo"},
			Car{vinNumber: "2", model: "Citroen"},
			// Car{vinNumber: "3", model: "Audi"},
			// Car{vinNumber: "4", model: "Jaguar"},
		})

	cars := []Entity{
		Car{vinNumber: "5", model: "Ferrari"},
		Car{vinNumber: "6", model: "Porshe"},
		Car{vinNumber: "7", model: "Nissan"}}

	wait := sync.WaitGroup{}
	ready := sync.WaitGroup{}
	mu := sync.Mutex{}

	wait.Add(2 * len(cars))
	ready.Add(len(cars))

	var res []Entity
	for _, car := range cars {
		time.Sleep(1 * time.Second)
		go func(e Entity) {
			defer ready.Done()
			defer wait.Done()

			err := m.Update(e)
			if err != nil {
				fmt.Println(err)
			}
		}(car)

		go func(id string) {
			ready.Wait()
			time.Sleep(time.Duration(rand.Intn(5)+5) * time.Second)
			defer wait.Done()

			resp, err := m.Read(id)
			if err != nil {
				fmt.Println(err)
			}

			mu.Lock()
			res = append(res, resp)
			mu.Unlock()
		}(car.Id())
	}
	wait.Wait()
	m.Print()

	fmt.Println("Finish, results:", res)
}

func ConcurrentlyAdding() {
	m := NewInMemoryCache(&fifo{}, 3, nil)

	cars := []Car{
		{vinNumber: "1", model: "Mondeo"},
		{vinNumber: "2", model: "Citroen"},
		{vinNumber: "3", model: "Audi"},
		// {vinNumber: "4", model: "Jaguar"},
		// {vinNumber: "5", model: "Porshe"},
		// {vinNumber: "6", model: "Ferrari"},
		// {vinNumber: "7", model: "Nissan"},
		// {vinNumber: "8", model: "Alfa Romeo"},
		// {vinNumber: "9", model: "Volvo"},
		// {vinNumber: "10", model: "Volkswagen"},
		// {vinNumber: "11", model: "Lamborgini"},
		// {vinNumber: "12", model: "BMW"},
	}
	rand.Seed(time.Now().UnixNano())

	wait := sync.WaitGroup{}
	ready := sync.WaitGroup{}
	wait.Add(len(cars))
	ready.Add(len(cars))

	for _, car := range cars {
		go func(c Car) {
			defer wait.Done()
			ready.Done()
			ready.Wait()

			time.Sleep(time.Duration(rand.Intn(15)) * time.Second)
			m.Update(c)
		}(car)
	}

	wait.Wait()
	m.Print()

	m.Update(Car{vinNumber: "11", model: "Lamborgini"})
	m.Print()

}
