package db

type Connector interface {
	IsConnected() bool
	Connect() error
	Close()
}
