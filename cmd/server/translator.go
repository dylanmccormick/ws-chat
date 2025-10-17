package server

import "context"

// This is the translator. Also known as the byte wizard. It translates []byte into Message{} and takes in some user context to be able to ascribe a Message to a User

type Translator struct {
}

func (t *Translator) BytesToMessage(ctx context.Context, data []byte) (Message, error) {

	return Message{}, nil
}
