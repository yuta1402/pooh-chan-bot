package poohchanbot

import (
	"errors"
	"math/rand"

	"github.com/line/line-bot-sdk-go/linebot"
)

type sendingMessageQueue struct {
	queue []linebot.SendingMessage
}

func newSendingMessageQueue() *sendingMessageQueue {
	return &sendingMessageQueue{
		queue: []linebot.SendingMessage{},
	}
}

func (q *sendingMessageQueue) enque(sendingMessage linebot.SendingMessage) {
	q.queue = append(q.queue, sendingMessage)
}

func (q *sendingMessageQueue) deque(sendingMessage linebot.SendingMessage) (linebot.SendingMessage, error) {
	if len(q.queue) <= 0 {
		return nil, errors.New("sending message queue is empty.")
	}

	top := q.queue[0]
	q.queue = q.queue[1:]

	return top, nil
}

func generateRandomTextMessageQueue(messages [][]string) *sendingMessageQueue {
	r := rand.Intn(len(messages))
	q := newSendingMessageQueue()

	for _, text := range messages[r] {
		q.enque(linebot.NewTextMessage(text))
	}

	return q
}
