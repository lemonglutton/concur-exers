package main

func main() {
	rooms := []*Room{NewRoom("#general"), NewRoom("#random")}
	s := NewServer(rooms)
	s.Listen()
}
