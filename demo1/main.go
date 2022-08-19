package main

import (
	"errors"
	"fmt"
	"log"
	"plugin"
	"reflect"
)

type Builder struct {
	result CPUSet
	done   bool
}

// NewBuilder returns a mutable CPUSet builder.
func NewBuilder() *Builder {
	return &Builder{
		result: CPUSet{
			elems: map[int]struct{}{},
		},
	}
}

// Add adds the supplied elements to the result. Calling Add after calling
// Result has no effect.
func (b *Builder) Add(elems ...int) {
	if b.done {
		return
	}
	for _, elem := range elems {
		b.result.elems[elem] = struct{}{}
	}
}

// Result returns the result CPUSet containing all elements that were
// previously added to this builder. Subsequent calls to Add have no effect.
func (b *Builder) Result() CPUSet {
	b.done = true
	return b.result
}

// CPUSet is a thread-safe, immutable set-like data structure for CPU IDs.
type CPUSet struct {
	elems map[int]struct{}
}

// NewCPUSet returns a new CPUSet containing the supplied elements.
func NewCPUSet(cpus ...int) CPUSet {
	b := NewBuilder()
	for _, c := range cpus {
		b.Add(c)
	}
	return b.Result()
}

// NewCPUSetInt64 returns a new CPUSet containing the supplied elements, as slice of int64.
func NewCPUSetInt64(cpus ...int64) CPUSet {
	b := NewBuilder()
	for _, c := range cpus {
		b.Add(int(c))
	}
	return b.Result()
}

// Size returns the number of elements in this set.
func (s CPUSet) Size() int {
	return len(s.elems)
}

// IsEmpty returns true if there are zero elements in this set.
func (s CPUSet) IsEmpty() bool {
	return s.Size() == 0
}

// Contains returns true if the supplied element is present in this set.
func (s CPUSet) Contains(cpu int) bool {
	_, found := s.elems[cpu]
	return found
}

// Equals returns true if the supplied set contains exactly the same elements
// as this set (s IsSubsetOf s2 and s2 IsSubsetOf s).
func (s CPUSet) Equals(s2 CPUSet) bool {
	return reflect.DeepEqual(s.elems, s2.elems)
}

func init() {
	log.Println("pkg1 init")
}

type MyInterface interface {
	M1()
}

func LoadAndInvokeSomethingFromPlugin(pluginPath string) error {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return err
	}

	// // 导出整型变量
	v, err := p.Lookup("V")
	if err != nil {
		return err
	}
	*v.(*CPUSet) = NewCPUSetInt64(1, 2, 3)

	// // 导出函数变量
	// f, err := p.Lookup("F_print_V")
	// if err != nil {
	// 	return err
	// }
	// f.(func())()

	// mySet := NewCPUSet(1, 2, 3)

	// // 导出自定义类型变量
	f1, err := p.Lookup("Foo")
	if err != nil {
		return err
	}
	i, ok := f1.(MyInterface)
	if !ok {
		return errors.New("f1 does not implement MyInterface")
	}
	i.M1()

	// 导出函数变量
	fM, err := p.Lookup("M")
	if err != nil {
		return err
	}
	fM.(func())()

	return nil
}

func main() {
	// fmt.Println("try to LoadAndInvokeSomethingFromPlugin...")
	err := LoadAndInvokeSomethingFromPlugin("../demo1-plugins/plugin1.so")
	if err != nil {
		fmt.Println("LoadAndInvokeSomethingFromPlugin error:", err)
		return
	}

}

/* de scris un plugin care afiseaza "hello world" pt of functie myfunc, si sa apelez din kubelet myfunc pt a afisa "hello world" */
