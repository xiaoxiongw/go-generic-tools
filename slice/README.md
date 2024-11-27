Slice工具类函数说明文件：

Add：    在切片的index出添加元素

Max：    获取切片最大值 (Number类型的切片)
Min：    获取切片最小值 (Number类型的切片)
Sum：    求和 (Number类型的切片)
// 上述三个函数在使用 float32 或者 float64 的时候要小心精度问题

Contains：判断Slice切片中是否包含某个元素,
ContainsFunc： 同上，应该优先使用Contains方法
ContainsAny： 判断切片中是否存在子切片中的任何一个元素
ContainsAnyFunc： 同上，应该优先使用ContainsAny
ContainsAll： 判断切片中是否存在子切片中的所有元素
ContainsAllFunc： 同上，应该优先使用ContainsAll
Delete： 删除 index 处的元素
FilterDelete： 删除符合条件的元素（考虑到性能问题，所有操作都会在原切片上进行）

DiffSet： 找出src和dst之间的差集（src 中存在但在 dst 中不存在的元素），已去重，并且返回顺序不固定
DiffSetFunc： 同上，应该优先使用DiffSet

IntersectSet： 取两个切片的交集（只支持comparable类型），已去重
IntersectSetFunc: 支持任意类型，优先使用IntersectSet

Find： 在Slice中查找元素，找到则返回；需要传入查找函数。
FindAll： 在Slice中查找所有符合条件的元素
Index： 在Slice中查询某个元素，找到则返回下标；未找到则返回-1
IndexFunc： 同上，应该优先使用Index

FilterMap： 对切片进行过滤，传入映射函数m，返回满足条件的元素组成的新切片
Map： 返回经映射函数m处理后的切片元素，返回的是一个新数组

ToMap： 将[]Ele映射到map[Key]Ele，从Ele中提取Key的函数fn由使用者提供
ToMapV： 将[]Ele映射到map[Key]Val，从Ele中提取Key和Val的函数fn由使用者提供

Reverse： 将切片反转（返回的是一个新的切片）
ReverseSelf： 将切片反转（在原来的基础上修改）

SymmetricDiffSet： 求两个切片的 对称差集（属于一个切片，但不属于两个切片的交集）
SymmetricDiffSetFunc： 优先使用 SymmetricDiffSet，已去重

UnionSet： 求两个切片的并集，只支持 comparable
UnionSetFunc： 求两个切片的并集，支持任意类型，优先使用 UnionSet，已去重；求并集函数作为参数传入














