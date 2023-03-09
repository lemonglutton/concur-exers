package main

type Car struct {
	vinNumber int
	model     string
	serialNum string
}

func (c Car) Id() int {
	return c.vinNumber
}

func (c Car) Model() string {
	return c.model
}

func (c Car) SerialNumber() string {
	return c.serialNum
}
