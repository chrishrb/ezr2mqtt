package api

type Receiver interface {
	Connect(errCh chan error)
}
