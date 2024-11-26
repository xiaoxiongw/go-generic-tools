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

package slice

/**
切片缩容
*/

func calCapacity(c, l int) (int, bool) {
	// 容量cal<64 ，没必要缩容
	if c <= 64 {
		return c, false
	}
	// 容量>2048，并且已使用的长度len没有达到cap的一半，缩容到原来的5/8
	if c > 2048 && (c/l >= 2) {
		factor := 0.625
		return int(float32(c) * float32(factor)), true
	}
	// 64<容量<=2048，并且已使用的长度len没有达到cap的1/4，缩容到原来的1/2
	if c <= 2048 && (c/l >= 4) {
		return c / 2, true
	}
	return c, false
}

func Shrink[T any](src []T) []T {
	// 获取长度len和容量cap
	c, l := cap(src), len(src)
	n, changed := calCapacity(c, l)
	if !changed {
		return src
	}
	// 重新创建一个切片，将原来的切片数据拷贝到新的切片中
	s := make([]T, 0, n)
	s = append(s, src...)
	return s
}
