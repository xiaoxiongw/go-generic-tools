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

// FilterMap 执行过滤并且转化
// 如果 m 的第二个返回值是 false，那么我们会忽略第一个返回值
// 即便第二个返回值是 false，后续的元素依旧会被遍历
func FilterMap[Src any, Dst any](src []Src, m func(idx int, src Src) (Dst, bool)) []Dst {
	res := make([]Dst, 0, len(src))
	for i, s := range src {
		dst, ok := m(i, s)
		if ok {
			res = append(res, dst)
		}
	}
	return res
}

// Map 创建一个新的切片dst，用于存储经过映射函数 m 处理后的元素
func Map[Src any, Dst any](src []Src, m func(idx int, src Src) Dst) []Dst {
	dst := make([]Dst, len(src))
	for i, s := range src {
		dst[i] = m(i, s)
	}
	return dst
}

// ToMap 将[]Ele映射到map[Key]Ele
// 从Ele中提取Key的函数fn由使用者提供
//
// 注意:
// 如果出现 i < j
// 设：
//
//	key_i := fn(elements[i])
//	key_j := fn(elements[j])
//
// 满足key_i == key_j 的情况，则在返回结果的resultMap中
// resultMap[key_i] = val_j
//
// 即使传入的字符串为nil，也保证返回的map是一个空map而不是nil
func ToMap[Ele any, Key comparable](elements []Ele, fn func(element Ele) Key) map[Key]Ele {
	return ToMapV(
		elements,
		func(element Ele) (Key, Ele) {
			return fn(element), element
		})
}

// ToMapV 将[]Ele映射到map[Key]Val
// 从Ele中提取Key和Val的函数fn由使用者提供
//
// 注意:
// 如果出现 i < j
// 设：
//
//	key_i, val_i := fn(elements[i])
//	key_j, val_j := fn(elements[j])
//
// 满足key_i == key_j 的情况，则在返回结果的resultMap中
// resultMap[key_i] = val_j
//
// 即使传入的字符串为nil，也保证返回的map是一个空map而不是nil
func ToMapV[Ele any, Key comparable, Val any](elements []Ele, fn func(element Ele) (Key, Val)) (resultMap map[Key]Val) {
	resultMap = make(map[Key]Val, len(elements))
	for _, element := range elements {
		k, v := fn(element)
		resultMap[k] = v
	}
	return
}

// 构造map（key是切片元素 value是空结构体）
func toMap[T comparable](src []T) map[T]struct{} {
	var dataMap = make(map[T]struct{}, len(src))
	for _, v := range src {
		// 使用空结构体,减少内存消耗
		dataMap[v] = struct{}{}
	}
	return dataMap
}

// deduplicateFunc 去除输入切片中的重复元素,并返回一个新的切片（equal函数作为参数传入）
func deduplicateFunc[T any](data []T, equal equalFunc[T]) []T {
	var newData = make([]T, 0, len(data))
	// 遍历切片
	for k, v := range data {
		// k为index，v为元素值
		// 判断当前切片元素之后的所有元素中，是否存在等于v的元素，若是不存在就将v添加到newData切片中
		// 用此方法可以去重
		if !ContainsFunc[T](data[k+1:], func(src T) bool {
			return equal(src, v)
		}) {
			newData = append(newData, v)
		}
	}
	return newData
}

// deduplicate 去除输入切片中的重复元素,并返回一个新的切片
func deduplicate[T comparable](data []T) []T {
	// 将切片转成map，这样可以去除切片中重复的元素
	dataMap := toMap[T](data)
	var newData = make([]T, 0, len(dataMap))
	for key := range dataMap {
		newData = append(newData, key)
	}
	return newData
}
