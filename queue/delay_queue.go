// Copyright 2021 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package queue

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-generic/internal/queue"
)

// DelayQueue 延时队列
// 每次出队的元素必然都是已经到期的元素，即 Delay() 返回的值小于等于 0
// 延时队列本身对时间的精确度并不是很高，其时间精确度主要取决于 time.Timer
// 所以如果你需要极度精确的延时队列，那么这个结构并不太适合你。
// 但是如果你能够容忍至多在毫秒级的误差，那么这个结构还是可以使用的
type DelayQueue[T Delayable] struct {
	q             queue.PriorityQueue[T] // 基于小顶堆的优先队列
	mutex         *sync.Mutex
	dequeueSignal *cond // 出队时发出信号
	enqueueSignal *cond // 入队时发出信号
}

func NewDelayQueue[T Delayable](c int) *DelayQueue[T] {
	m := &sync.Mutex{}
	res := &DelayQueue[T]{
		// 根据延时时间
		q: *queue.NewPriorityQueue[T](c, func(src T, dst T) int {
			// src 来源  dst 目标
			srcDelay := src.Delay()
			dstDelay := dst.Delay()
			// 来源delay>目标delay return>0
			if srcDelay > dstDelay {
				return 1
			}
			if srcDelay == dstDelay {
				return 0
			}
			// 来源delay<目标delay return<0
			return -1
		}),
		mutex:         m,
		dequeueSignal: newCond(m),
		enqueueSignal: newCond(m),
	}
	return res
}

func (d *DelayQueue[T]) Enqueue(ctx context.Context, t T) error {
	for {
		select {
		// 先检测 ctx 有没有过期
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		// ctx 没有过期
		d.mutex.Lock()
		// 对小顶堆的优先队列进行入队操作
		err := d.q.Enqueue(t)
		switch err {
		// 入队未发生错误
		case nil:
			d.enqueueSignal.broadcast()
			return nil
		// 队列已满
		case queue.ErrOutOfCapacity:
			// 获取 dequeueSignal 信号通道
			signal := d.dequeueSignal.signalCh()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-signal: // 在此处阻塞
			}
		default:
			d.mutex.Unlock()
			return fmt.Errorf("ekit: 延时队列入队的时候遇到未知错误 %w，请上报", err)
		}
	}
}

func (d *DelayQueue[T]) Dequeue(ctx context.Context) (T, error) {
	var timer *time.Timer
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()
	for {
		select {
		// 先检测 ctx 有没有过期
		case <-ctx.Done():
			var t T
			return t, ctx.Err()
		default:
		}
		d.mutex.Lock()
		val, err := d.q.Peek()
		switch err {
		case nil:
			delay := val.Delay()
			if delay <= 0 {
				val, err = d.q.Dequeue()
				d.dequeueSignal.broadcast()
				// 理论上来说这里 err 不可能不为 nil
				return val, err
			}
			signal := d.enqueueSignal.signalCh()
			if timer == nil {
				timer = time.NewTimer(delay)
			} else {
				timer.Reset(delay)
			}
			select {
			case <-ctx.Done():
				var t T
				return t, ctx.Err()
			case <-timer.C:
				// 到了时间
				d.mutex.Lock()
				// 原队头可能已经被其他协程先出队，故再次检查队头
				val, err := d.q.Peek()
				if err != nil || val.Delay() > 0 {
					d.mutex.Unlock()
					continue
				}
				// 验证元素过期后将其出队
				val, err = d.q.Dequeue()
				d.dequeueSignal.broadcast()
				return val, err
			case <-signal:
				// 进入下一个循环。这里可能是有新的元素入队，也可能是到期了
			}
		case queue.ErrEmptyQueue:
			signal := d.enqueueSignal.signalCh()
			select {
			case <-ctx.Done():
				var t T
				return t, ctx.Err()
			case <-signal:
			}
		default:
			d.mutex.Unlock()
			var t T
			return t, fmt.Errorf("ekit: 延时队列出队的时候遇到未知错误 %w，请上报", err)
		}
	}
}

type Delayable interface {
	Delay() time.Duration
}

type cond struct {
	signal chan struct{}
	l      sync.Locker
}

func newCond(l sync.Locker) *cond {
	return &cond{
		signal: make(chan struct{}),
		l:      l,
	}
}

// broadcast 唤醒等待者
// 如果没有人等待，那么什么也不会发生
// 必须加锁之后才能调用这个方法
// 广播之后锁会被释放，这也是为了确保用户必然是在锁范围内调用的
func (c *cond) broadcast() {
	signal := make(chan struct{})
	old := c.signal
	c.signal = signal
	c.l.Unlock()
	close(old)
}

// signalCh 返回一个 channel，用于监听广播信号
// 必须在锁范围内使用
// 调用后，锁会被释放，这也是为了确保用户必然是在锁范围内调用的
func (c *cond) signalCh() <-chan struct{} {
	res := c.signal
	c.l.Unlock()
	return res
}
