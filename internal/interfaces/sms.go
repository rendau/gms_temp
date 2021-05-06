package interfaces

// Sms - is interface of sms
type Sms interface {
	Send(phones string, msg string) bool
	GetBalance() (bool, float64)
}
