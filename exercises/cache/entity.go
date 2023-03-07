package main

type Entity interface {
	Id() string
}

type Car struct {
	vinNumber string
	model     string
	serialNum string
}

func (c Car) Id() string {
	return c.vinNumber
}

func (c Car) Model() string {
	return c.model
}

func (c Car) SerialNumber() string {
	return c.serialNum
}
