package training

import (
	"fmt"
	"github.com/magicnana999/im/util/id"
	"github.com/timandy/routine"
	"strings"
	"sync"
	"time"
)

func CompareAndSwap() {
	var m sync.Map

	key := "key"
	val := "{\"key\":100}"

	f := func() {
		time.Sleep(time.Second)
		swapped := m.CompareAndSwap(key, "", val)
		fmt.Println(routine.Goid(), swapped)
	}

	go f()
	go f()
	go f()

	time.Sleep(time.Second * 5)
	fmt.Println("OK")
}

func LoadAndStore() {
	var m sync.Map

	key := "key"

	f := func() {
		val := strings.ToLower(id.GenerateXId())
		v, e := m.LoadOrStore(key, val)
		fmt.Println(routine.Goid(), "当前routine:", val, "是否已存在", e, "map中的值", v)
	}

	go f()
	go f()
	go f()

	time.Sleep(time.Second * 5)
	fmt.Println("OK")
}

func rangeFunc() {

	type foo struct {
		name string
		val  string
	}

	f1 := &foo{"foo1", "foo1value"}
	f2 := &foo{"foo2", "foo2value"}
	f3 := &foo{"foo3", "foo2value"}

	var m sync.Map

	m.LoadOrStore(f1.name, f1)
	m.LoadOrStore(f2.name, f2)
	m.LoadOrStore(f3.name, f3)

	m.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})

}
