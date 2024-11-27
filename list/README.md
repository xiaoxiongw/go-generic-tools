当前实现的列表：
ArrayList：Java中的ArrayList
LinkedList：双向链表LinkedList
ConcurrentList：线程安全的list（加读写锁），在list基础上封装



type List[T any] interface 包含的方法有：
Get：    Get 返回对应下标的元素，下标超出范围的情况下，返回错误
Append： Append 在末尾追加元素
Add：    Add 在特定下标处增加一个新元素；下标不在[0, Len()]范围之内，返回错误；如果index == Len()，作用等同于Append
Set：    Set 重置 index 位置的值；下标超出范围，应该返回错误
Delete： Delete 删除特定下标的元素，并且返回该位置的值；index 超出下标，应该返回错误
Len：    Len 返回长度
Cap：    Cap 返回容量
Range：  Range 遍历 List 的所有元素
AsSlice：AsSlice 将 List 转化为一个切片；在没有元素的情况下，回一个长度和容量都为 0 的切片