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
	"sync/atomic"
	"unsafe"

	"github.com/go-generic/internal/queue"
)

// ConcurrentLinkedQueue 并发安全的无界队列
type ConcurrentLinkedQueue[T any] struct {
	// *node[T]
	head unsafe.Pointer
	// *node[T]
	tail unsafe.Pointer
}

// NewConcurrentLinkedQueue 创建一个新的并发安全的无界队列
func NewConcurrentLinkedQueue[T any]() *ConcurrentLinkedQueue[T] {
	// 创建一个空node，头指针 尾指针都指向这个地址
	head := &node[T]{}
	ptr := unsafe.Pointer(head)
	return &ConcurrentLinkedQueue[T]{
		head: ptr,
		tail: ptr,
	}
}

// Enqueue 并发安全无界队列入队
func (c *ConcurrentLinkedQueue[T]) Enqueue(t T) error {
	// 创建入队节点，并获取节点指针ptr
	newNode := &node[T]{val: t}
	newPtr := unsafe.Pointer(newNode)
	for {
		// 获取尾节点
		tailPtr := atomic.LoadPointer(&c.tail)
		tail := (*node[T])(tailPtr)
		// 获取尾节点的next指针
		tailNext := atomic.LoadPointer(&tail.next)
		// 尾节点的next指针不为空，说明尾节点在循环过程中被修改，continue 重新开始
		if tailNext != nil {
			// 已经被人修改了，我们不需要修复，因为预期中修改的那个人会把 c.tail 指过去
			continue
		}
		// 尝试将新节点接到尾节点后面
		// 比较 &tail.next 的当前值是否等于 tailNext,如果相等,则将 tail.next 的值更新为 newPtr
		if atomic.CompareAndSwapPointer(&tail.next, tailNext, newPtr) {
			// 如果失败也不用担心，说明有人抢先一步了
			// 添加成功，更新队列的tail指针
			atomic.CompareAndSwapPointer(&c.tail, tailPtr, newPtr)
			return nil
		}
	}
}

// Dequeue 并发安全无界队列出队
func (c *ConcurrentLinkedQueue[T]) Dequeue() (T, error) {
	for {
		// 获取队列头节点
		headPtr := atomic.LoadPointer(&c.head)
		head := (*node[T])(headPtr)
		// 获取队列尾节点
		tailPtr := atomic.LoadPointer(&c.tail)
		tail := (*node[T])(tailPtr)
		// 判断队列是否为空
		if head == tail {
			// 不需要做更多检测，在当下这一刻，我们就认为没有元素，即便这时候正好有人入队
			// 但是并不妨碍我们在它彻底入队完成——即所有的指针都调整好——之前，
			// 认为其实还是没有元素
			var t T
			return t, queue.ErrEmptyQueue
		}
		// 更改head指针为当前头节点的next指向的节点
		headNextPtr := atomic.LoadPointer(&head.next)
		if atomic.CompareAndSwapPointer(&c.head, headPtr, headNextPtr) {
			// 返回队首节点（head 指针指向队列的头节点,但实际上队列的第一个元素是 head.next 指向的节点）
			headNext := (*node[T])(headNextPtr)
			return headNext.val, nil
		}
	}
}

type node[T any] struct {
	val T
	// *node[T]
	next unsafe.Pointer
}
