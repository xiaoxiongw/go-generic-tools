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
	"errors"
	"github.com/go-generic"

	"github.com/go-generic/internal/slice"
)

var (
	ErrOutOfCapacity = errors.New("queue: 超出最大容量限制")
	ErrEmptyQueue    = errors.New("queue: 队列为空")
)

// PriorityQueue 是一个基于小顶堆的优先队列
// 当capacity <= 0时，为无界队列，切片容量会动态扩缩容
// 当capacity > 0 时，为有界队列，初始化后就固定容量，不会扩缩容
type PriorityQueue[T any] struct {
	// 用于比较前一个元素是否小于后一个元素
	compare generic.Comparator[T]
	// 队列容量
	capacity int
	// 队列中的元素，为便于计算父子节点的index，0位置留空，根节点从1开始
	data []T
}

// Len 优先队列长度
func (p *PriorityQueue[T]) Len() int {
	return len(p.data) - 1
}

// Cap 无界队列返回0，有界队列返回创建队列时设置的值
func (p *PriorityQueue[T]) Cap() int {
	return p.capacity
}

// IsBoundless 判断是否是无界队列，true无界队列 false有界队列
func (p *PriorityQueue[T]) IsBoundless() bool {
	return p.capacity <= 0
}

// 判断优先队列是否满（有界）
func (p *PriorityQueue[T]) isFull() bool {
	return p.capacity > 0 && len(p.data)-1 == p.capacity
}

// 判断优先队列是否为空
func (p *PriorityQueue[T]) isEmpty() bool {
	return len(p.data) < 2
}

// 回优先队列中的最小元素,而不将其从队列中移除
func (p *PriorityQueue[T]) Peek() (T, error) {
	if p.isEmpty() {
		var t T
		return t, ErrEmptyQueue
	}
	return p.data[1], nil
}

// Enqueue 新元素入队
func (p *PriorityQueue[T]) Enqueue(t T) error {
	// 判断是否满
	if p.isFull() {
		return ErrOutOfCapacity
	}

	p.data = append(p.data, t)
	//进行上浮操作,将新元素上移到合适的位置,以满足小顶堆的性质
	node, parent := len(p.data)-1, (len(p.data)-1)/2
	//从新元素的位置开始,与其父节点进行比较。
	//如果新元素小于父节点,则交换它们的位置。
	//重复这个过程,直到新元素的位置满足小顶堆的性质或到达根节点。
	for parent > 0 && p.compare(p.data[node], p.data[parent]) < 0 {
		p.data[parent], p.data[node] = p.data[node], p.data[parent]
		node = parent
		parent = parent / 2
	}

	return nil
}

// 优先队列出库
func (p *PriorityQueue[T]) Dequeue() (T, error) {
	if p.isEmpty() {
		var t T
		return t, ErrEmptyQueue
	}

	pop := p.data[1]
	// 将最后一个元素移动到堆顶,并缩小切片 data 的长度
	p.data[1] = p.data[len(p.data)-1]
	p.data = p.data[:len(p.data)-1]
	// 如果是无界队列，则对data切片缩容
	p.shrinkIfNecessary()
	// 从data[1]开始往后，构造成一个堆序列
	p.heapify(p.data, len(p.data)-1, 1)
	return pop, nil
}

// 对无界队列进行缩容
func (p *PriorityQueue[T]) shrinkIfNecessary() {
	if p.IsBoundless() {
		p.data = slice.Shrink[T](p.data)
	}
}

// 将一个无序的数组或线性数据结构转换为一个满足堆性质的数据结构
func (p *PriorityQueue[T]) heapify(data []T, n, i int) {
	minPos := i
	for {
		if left := i * 2; left <= n && p.compare(data[left], data[minPos]) < 0 {
			minPos = left
		}
		if right := i*2 + 1; right <= n && p.compare(data[right], data[minPos]) < 0 {
			minPos = right
		}
		if minPos == i {
			break
		}
		//
		data[i], data[minPos] = data[minPos], data[i]
		i = minPos
	}
}

// NewPriorityQueue 创建优先队列 capacity <= 0 时，为无界队列，否则有有界队列
func NewPriorityQueue[T any](capacity int, compare generic.Comparator[T]) *PriorityQueue[T] {
	sliceCap := capacity + 1 // 切片长度 = 容量+1
	if capacity < 1 {
		capacity = 0
		sliceCap = 64
	}
	return &PriorityQueue[T]{
		capacity: capacity,
		data:     make([]T, 1, sliceCap),
		compare:  compare,
	}
}
