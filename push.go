package push

type Pusher interface {
	Setup() error
	Send(tokens []Token) error
}
