package pass_by_value

import "fmt"

func init() {
	fmt.Println("一切都是值传递")
}

// 基本类型（int, float, string 等）
// 行为：传递的是值的副本。
// 复制：传递时会复制整个值，而不是地址。
func basicInt() {
	f1 := func(a int) int {
		a = a + 100
		fmt.Println("in func1:", a)
		return a
	}

	v := 10
	fmt.Println("before func1 v:", v)
	f1(v)
	fmt.Println("after func1:", v)
}

// 基本类型（int, float, string 等）
// 行为：传递的是值的副本。
// 复制：传递时会复制整个值，而不是地址。
func basicString() {
	f1 := func(a string) string {
		a = fmt.Sprintf("%s %s", a, "ahaah")
		fmt.Println("in func1:", a)
		return a
	}

	v := "10"
	fmt.Println("before func1 v:", v)
	f1(v)
	fmt.Println("after func1:", v)
}

//1. 基本类型（int, float, string 等）
//行为：传递的是值的副本。
//复制内容：传递时会复制整个值，而不是地址。

//2. 指针类型
//行为：传递的是指针的值（即地址）的副本。
//复制内容：复制的是地址，而不是地址指向的内容。
//与 Java 的对比：你提到“和 Java 中的对象传递一样”，这基本正确。Java 中的对象引用传递也是复制引用（地址），而不是对象本身。Go 中的指针传递也是如此。

//3. 集合类型（map, slice, chan）
//行为：传递的是值的副本，但这些类型内部包含指针。
//复制内容：复制的是这些类型的“头部”结构（包含指向底层数据的指针），而不是底层数据本身。
//细节：
//slice：复制的是 {data指针, len, cap} 结构。
//map：复制的是指向映射表的指针。
//chan：复制的是指向通道结构的指针。

//4. 结构体
//行为：传递的是整个结构体的副本。
//复制内容：复制的是结构体中所有字段的值（包括嵌套的字段）。

//5. 结构体的指针
//行为：传递的是指针的副本。
//复制内容：复制的是地址，而不是结构体本身。

//6，结构体中引用其他结构体
//func processHeap(h *MinHeap[int]) {
//    h.Push(3)
//}
//在这种情况下， 其内部的impl 是 *minHeapImpl[T] 还是 minHeapImpl[T] 都无所谓，因为在processHeap调用时，复制了地址，地址指向的内存区域不变，所以即使impl 是 minHeapImpl[T]也不会复制一份副本，对不对？
//
//func processHeap(h MinHeap[int]) {
//    h.Push(3)
//}
//在这种情况下，h本身会被复制一份副本，如果impl是指针，那么将复制地址，不会复制整个结构体的内容，相反则会复制 minHeapImpl 结构题的所有内容，对不对？

//这岂不是间接的论证了，在一个结构体的receiver是指针还是值的问题上，为了避免修改了副本而没有改变原始的内容这一错误的结果，receiver应该是指针

//6，结构体内部引用其他结构体指针还是值？
//在外城结构体receiver是指针的情况下，结构体内部引用其他结构体，可以不是指针，是结构体的值。这样是高效还是低效？
//1. 值类型嵌入（当前设计：impl minHeapImpl[T]）
//内存布局
//MinHeap[T] 的内存布局是连续的：
//impl（包含 data 切片和 less 函数指针）直接存储在 MinHeap[T] 的内存块中。
//整个结构体的大小是固定的（impl 的字段大小 + mu + maxSize）。
//访问 h.impl.data 时，只需一次偏移计算，无需额外的指针解引用。
//性能
//优点：
//高效访问：访问 impl 的字段（如 data 或 less）是直接的内存偏移，不涉及指针解引用，减少了 CPU 的间接寻址开销。
//内存局部性：impl 和 MinHeap[T] 的其他字段（如 mu）在同一块内存中，缓存命中率更高。
//	反之则不连续，访问 h.impl.data 需要两次解引用：先解 h.impl 得到 minHeapImpl[T] 的地址，再偏移到 data。
//无额外分配：创建 MinHeap[T] 时，impl 不需要单独分配内存，与外层结构体一起分配。

//在接收者是指针的情况下（*MinHeap[T]），结构体内部引用其他结构体使用值类型（impl minHeapImpl[T]）是高效的，原因如下：
//
//避免不必要开销：值类型嵌入减少了内存分配和解引用的成本。
//语义清晰：impl 是 MinHeap[T] 的核心部分，嵌入值类型符合逻辑。
//性能优化：内存连续性和直接访问提高了效率。
//无拷贝问题：指针接收者确保了外层结构体不被复制，嵌套结构体的值类型不会引发副本问题。
//如果改为指针类型（*minHeapImpl[T]），会引入低效因素（额外分配、间接访问），而收益有限，除非有特定需求（如共享或动态替换 impl）。

//7 内部结构体的receiver
//内部依赖的minHeapImpl本身也需要自己的receiver来暴露方法给MinHeap内部使用（必须的，比如必须实现container/interface接口），此时 minHeapImpl的receiver是指针还是结构体有什么影响？ 怎样才好？

//推荐使用 *minHeapImpl[T] 作为接收者，理由是正确性、一致性和性能的综合优势。
//保持 MinHeap[T] 的 impl 为值类型嵌入，无需改为指针类型，因为外层已经是 *MinHeap[T]，拷贝问题已解决。
