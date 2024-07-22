package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	log.Println("执行开始")

	log.Println(add(1, 1.2))
	log.Println()

	log.Println(toSlice(1, 2, 3))
	log.Println()

	s := Stack[int]{}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	log.Println(s.data)
	log.Println(s.Pop())
	log.Println(s.data)
	log.Println()

	// numbers := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	// Sort(numbers)
	// log.Println(numbers)
	// log.Println()

	log.Println("执行完成")

}

type Number interface {
	int8 | int16 | int | int32 | int64 | uint8 | uint16 | uint | uint32 | uint64 | float32 | float64
}

func add[T Number](a, b T) T {
	return a + b
}

// 多个不同类型的参数类型
func mutilGenericsParams[T string, X Number](a T, b X) {

}

func toSlice[T any](args ...T) []T {
	return args
}

// 泛型类型
// 泛型栈类型
type Stack[T any] struct {
	data []T
}

func (s *Stack[T]) Push(a T) {
	s.data = append(s.data, a)
}

func (s *Stack[T]) Pop() T {
	l := len(s.data)
	r := s.data[l-1]
	s.data = s.data[:l-1]
	return r
}

// 类型约束可以让泛型函数或类型只接受特定类型的参数。
// 在 Go 中，类型约束可以使用 interface{} 类型和类型断言来实现。
// 例如，下面是一个泛型函数，它可以接受实现了 fmt.Stringer 接口的类型
func Print[T fmt.Stringer](a T) {
	log.Println(a.String())
}

// T 和 U 分别表示实现了 fmt.Stringer 和 io.Reader 接口的任意类型，函数接受一个类型为 T 的参数和一个类型为 U 的参数，
// 并调用它们的方法输出其字符串表示和读取数据
func Print2[T fmt.Stringer, U io.Reader](a T, b U) {
	log.Println(a.String())
	io.Copy(os.Stdout, b)
}

// func Sort[T comparable](s []T) {
// 	sort.Slice(s, func(i, j int) bool {
// 		return s[i] < s[j]
// 	})
// }
