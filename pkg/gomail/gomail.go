package gomail

type Gomail interface {
	SendEmail(request *SendEmail) error
}
