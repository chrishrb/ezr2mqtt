package api

type Message struct {
	Room string
	Type string
	Data any
}

type Meta struct {
	Name string
}
