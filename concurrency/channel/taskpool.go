package channel

import (
	"context"
	"sync"
)

type Task func()

type TaskPool struct {
	tasks     chan Task
	close     chan struct{}
	closeOnce sync.Once
}

func NewTaskPool(numG int, capacity int) *TaskPool {
	res := &TaskPool{
		tasks: make(chan Task, capacity),
		close: make(chan struct{}),
	}

	for i := 0; i < numG; i++ {
		go func() {
			for {
				select {
				case <-res.close:
					return
				case t := <-res.tasks:
					t()
				}
			}
		}()
	}

	return res
}

func (t *TaskPool) Submit(ctx context.Context, job Task) error {
	// 这里不需要 for 循环，因为值添加一次任务
	select {
	// 可以临时停止这个任务，这里也可以是定时 ctx
	case <-ctx.Done():
		return ctx.Err()
	case t.tasks <- job:
	}
	return nil
}

func (t *TaskPool) Close() error {
	// 这种写法执行让一个 goroutine 接收结束信号
	//t.close <- struct{}{}

	//close(t.close)

	t.closeOnce.Do(func() {
		close(t.close)
	})

	return nil
}
