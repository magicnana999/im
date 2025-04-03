package heap

import "testing"

// 测试用的比较函数
func lessFunc(i, j int) bool {
	return i < j
}

// go test -v -run=TestMinHeapAdd
func TestMinHeapAdd(t *testing.T) {
	// 功能测试
	h := NewMinHeap(lessFunc, 5)
	values := []int{3, 1, 4, 1, 5}
	for _, v := range values {
		err := h.Add(v)
		if err != nil {
			t.Errorf("Add(%d) failed: %v", v, err)
		}
	}
	if h.Count() != 5 {
		t.Errorf("Expected count 5, got %d", h.Count())
	}
	top, err := h.Top()
	if err != nil || top != 1 {
		t.Errorf("Expected top 1, got %d, err: %v", top, err)
	}
	err = h.Add(6) // 超出最大容量
	if err != FullHeap {
		t.Errorf("Expected FullHeap error, got %v", err)
	}
}

// go test -bench=BenchmarkMinHeapAddSingle -benchtime=5s -v -run=^$
func BenchmarkMinHeapAddSingle(b *testing.B) {
	h := NewMinHeap(lessFunc, 0) // 无容量限制
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Add(i)
	}
}

// go test -bench=BenchmarkMinHeapAddConcurrent -benchtime=5s -v -run=^$
// go test -bench=BenchmarkMinHeapAddConcurrent -benchtime=5s -v -run=^$ -race
func BenchmarkMinHeapAddConcurrent(b *testing.B) {
	h := NewMinHeap(lessFunc, 0)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Add(100)
		}
	})
}

// go test -v -run=TestMinHeapRemove
func TestMinHeapRemove(t *testing.T) {
	// 功能测试
	h := NewMinHeap(lessFunc, 0)
	values := []int{3, 1, 4, 1, 5}
	for _, v := range values {
		h.Add(v)
	}
	expected := []int{1, 1, 3, 4, 5} // 最小堆排序后
	for i, exp := range expected {
		item, err := h.Remove()
		if err != nil {
			t.Errorf("Remove failed: %v", err)
		}
		if item != exp {
			t.Errorf("Expected %d at index %d, got %d", exp, i, item)
		}
	}
	_, err := h.Remove() // 空堆
	if err != EmptyHeap {
		t.Errorf("Expected EmptyHeap error, got %v", err)
	}
}

// go test -bench=BenchmarkMinHeapRemoveSingle -benchtime=5s -v -run=^$
func BenchmarkMinHeapRemoveSingle(b *testing.B) {
	h := NewMinHeap(lessFunc, 0)
	for i := 0; i < b.N; i++ {
		h.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Remove()
	}
}

// go test -bench=BenchmarkMinHeapRemoveConcurrent -benchtime=5s -v -run=^$
func BenchmarkMinHeapRemoveConcurrent(b *testing.B) {
	h := NewMinHeap(lessFunc, 0)
	for i := 0; i < 10000; i++ {
		h.Add(i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Remove()
		}
	})
}

// go test -v -run=TestMinHeapLen
func TestMinHeapLen(t *testing.T) {
	// 功能测试
	h := NewMinHeap(lessFunc, 0)
	if h.Count() != 0 {
		t.Errorf("Expected count 0, got %d", h.Count())
	}
	h.Add(1)
	h.Add(2)
	if h.Count() != 2 {
		t.Errorf("Expected count 2, got %d", h.Count())
	}
	h.Remove()
	if h.Count() != 1 {
		t.Errorf("Expected count 1, got %d", h.Count())
	}
}

// go test -v -run=TestMinHeapIterate
func TestMinHeapIterate(t *testing.T) {
	// 功能测试
	h := NewMinHeap(lessFunc, 0)
	values := []int{3, 1, 4, 1, 5}
	for _, v := range values {
		h.Add(v)
	}
	result := h.Iterate()
	expected := []int{1, 1, 3, 4, 5} // 最小堆排序后
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
		return
	}
	for i, exp := range expected {
		if result[i] != exp {
			t.Errorf("Expected %d at index %d, got %d", exp, i, result[i])
		}
	}
	// 验证原堆未受影响
	if h.Count() != 5 {
		t.Errorf("Iterate modified heap, count expected 5, got %d", h.Count())
	}
}
