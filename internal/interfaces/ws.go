package interfaces

type Ws interface {
	Send2User(usrId int64, data map[string]string) bool
	Send2Users(usrIds []int64, data map[string]string) bool
}
