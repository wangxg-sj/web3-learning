package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	num := 10
	plusTen(&num)
	fmt.Println("num加10后的结果:", num)
	slice := []int{1, 2, 3, 4, 5}
	doubleSliceItem(&slice)
	fmt.Println("slice每个元素乘以2后的结果:", slice)
	goRoutine()

	taskSchedule([]func(int){
		func(i int) {
			fmt.Printf("任务%d执行\n", i)
			time.Sleep(time.Duration(i) * time.Second)
		},
		func(i int) {
			fmt.Printf("任务%d执行\n", i)
			time.Sleep(time.Duration(i) * time.Second)
		},
	})

	rect := Rectangle{width: 2, height: 3}
	fmt.Println("矩形的面积:", rect.Area())
	fmt.Println("矩形的周长:", rect.Perimeter())

	circle := Circle{radius: 2}
	fmt.Println("圆的面积:", circle.Area())
	fmt.Println("圆的周长:", circle.Perimeter())

	emp := Employee{
		Person: Person{
			name: "张三",
			age:  30,
		},
		employeeID: 1001,
	}
	emp.PrintInfo()

	goRoutineComm(make(chan int, 10))
	// 锁机制
	safeCount()
	atomicCount()
}

func plusTen(x *int) {
	*x += 10
}

func doubleSliceItem(s *[]int) {
	for i, v := range *s {
		(*s)[i] = v * 2
	}
}

func goRoutine() {
	num := 10
	//开启一个带参数的goroutine打印从0到num的奇数

	go func(n int) {
		fmt.Print("打印从0到num的奇数:")
		for i := 0; i <= n; i++ {
			if i%2 != 0 {
				fmt.Print(i)
			}
		}
		fmt.Println()
	}(num)
	// 打印偶数
	go func(n int) {
		fmt.Print("打印从0到num的偶数:")
		for i := 0; i <= n; i++ {
			if i%2 == 0 {
				fmt.Print(i)
			}
		}
		fmt.Println()
	}(num)

	time.Sleep(10 * time.Second)
}

func taskSchedule(tasks []func(int)) {
	var wg sync.WaitGroup
	for i, f := range tasks {
		j := i + 1
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			//开始时间
			start := time.Now()
			f(j)
			fmt.Printf("任务%d执行时间: %v\n", j, time.Since(start))
		}(i)
	}
	wg.Wait()
	fmt.Println("所有任务执行完成")
}

// Shape 面向对象
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	width, height float64
}

type Circle struct {
	radius float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}
func (c Rectangle) Perimeter() float64 {
	return 2 * (c.width + c.height)
}

func (c Circle) Perimeter() float64 {
	return math.Pi * c.radius * c.radius
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

type Person struct {
	name string
	age  int
}

type Employee struct {
	Person
	employeeID int
}

func (e Employee) PrintInfo() {
	fmt.Printf("姓名: %s, 年龄: %d, 员工ID: %d\n", e.name, e.age, e.employeeID)
}

// Channel
func goRoutineComm(c chan int) {
	var wg sync.WaitGroup

	wg.Add(1)
	// 只读channel
	go func(reader <-chan int) {
		defer wg.Done()
		for {
			select {
			case data, ok := <-reader:
				if ok {
					fmt.Printf("只读channel接收: %d\n", data)
				} else {
					return
				}
			}
		}
	}(c)
	wg.Add(1)
	// 只写channel
	go func(writer chan<- int) {
		defer wg.Done()
		func(writer chan<- int) {
			defer close(writer)
			for i := 0; i < 10; i++ {
				writer <- i
				fmt.Printf("只写channel发送: %d\n", i)
				time.Sleep(1 * time.Second)
			}
		}(writer)

	}(c)
	wg.Wait()
}

// 锁机制
func safeCount() {
	var mutex sync.Mutex
	var wg sync.WaitGroup
	count := 0
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			mutex.Lock()
			defer mutex.Unlock()
			for j := 0; j < 1000; j++ {
				count++
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("最终计数: %d\n", count)
}

func atomicCount() {
	var wg sync.WaitGroup
	// atomic 原子操作
	var count int32 = 0
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				atomic.AddInt32(&count, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("最终计数: %d\n", count)
}
