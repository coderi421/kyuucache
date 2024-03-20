package channel

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Broker struct {
	mutex sync.RWMutex
	//chans []chan Msg
	chans map[string][]chan Msg
}

type Msg struct {
	Topic   string
	Content string
}

func (b *Broker) Subscribe(topic string, capacity int) (<-chan Msg, error) {
	if b.chans == nil {
		return nil, errors.New("queue 尚不可用或已经关闭")
	}
	c := make(chan Msg, capacity)

	b.mutex.Lock()
	defer b.mutex.Unlock()
	res, ok := b.chans[topic]
	if !ok {
		b.chans[topic] = []chan Msg{c}
		return c, nil
	}
	b.chans[topic] = append(res, c)
	return c, nil
}

func (b *Broker) Send(ctx context.Context, msg Msg) error {
	// 可以进一步实现，满队列后的等待
	//_, isDeadline := ctx.Deadline()
	if b.chans == nil {
		return errors.New("queue 尚不可用或已经关闭")
	}
	chans, ok := b.chans[msg.Topic]
	if len(chans) == 3 {
		fmt.Println(111)
	}
	if !ok {
		return errors.New("投递topic消息错误")
	}
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for _, ch := range chans {
		select {
		case ch <- msg:
		default:
			return errors.New("队列已满: " + msg.Topic)
		}
	}

	return nil
}

func (b *Broker) Close() error {
	b.mutex.Lock()
	chans := b.chans
	b.chans = nil
	//b.chans = map[string][]chan Msg{}
	b.mutex.Unlock()

	for _, chs := range chans {
		// 避免了重复 close chan 的问题
		for _, ch := range chs {
			close(ch)
		}
	}
	return nil
}
