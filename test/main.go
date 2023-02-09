/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2022-03-23 10:19:50
 * @LastEditTime: 2023-02-09 16:45:04
 */
package main

import (
	"log"
	"time"
)

func main() {
	forever := make(chan bool)
	log.Println("执行开始")
	test()
	log.Println("执行完成")
	<-forever
}

func test() {
	// 构建一个通道
	ch := make(chan int)
	// 开启一个并发匿名函数
	go func() {
		// 从3循环到0
		for i := 3; i >= 0; i-- {
			// 发送3到0之间的数值
			ch <- i
			// 每次发送完时等待
			time.Sleep(time.Second)
		}
	}()
	// 遍历接收通道数据
	// for data := range ch {
	// 	// 打印通道数据
	// 	log.Println(data)
	// 	// 当遇到数据0时, 退出接收循环
	// 	//   if data == 0 {
	// 	// 		  break
	// 	//   }
	// }

	for {
		select {
		case data := <-ch:
			log.Println(data)
		default:
			// time.Sleep(time.Second)
			// log.Println("无数据")

		}
	}

}
