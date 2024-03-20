package channel

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBroker_Send(t *testing.T) {
	// test in concurrency
	t.Parallel()
	b := &Broker{
		chans: map[string][]chan Msg{},
	}

	//	发送者
	go func() {
		time.Sleep(time.Second * 5)
		for {
			err := b.Send(context.Background(), Msg{
				Topic:   "tian",
				Content: time.Now().String(),
			})
			if err != nil {
				t.Log(err)
				return
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 生产者
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("消费者 %d", i)
		go func() {
			fmt.Println(name)
			defer wg.Done()
			msgs, err := b.Subscribe("tian", 100)
			if err != nil {
				t.Log(err)
				return
			}

			for msg := range msgs {
				fmt.Println(name, msg)
			}
		}()
	}
	wg.Wait()
	time.Sleep(10 * time.Second)
}
