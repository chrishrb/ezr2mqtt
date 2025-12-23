package api

type Message struct {
	Room int
	Type string
	Data any
}

type Meta struct {
	Name string
}
