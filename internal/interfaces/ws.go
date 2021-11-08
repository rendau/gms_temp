package interfaces

type Ws interface {
	Send(channel string, data interface{})
}
