package main

import "fmt"

func huiwen(x string) bool {
	lens := len(x)     //由于x只能是整数，所以使用len()得到整数x的长度
	for i := range x { //用range循环判断第i位和第lens-i-1位是否一样(判断是否为回文的条件)
		if x[i] != x[lens-i-1] {
			return false //不一样说明不是回文，提前中断函数返回false
		}
	}
	return true //循环结束说明是回文，返回true
}

func main() {
	var x string      //把整数x视为字符串
	fmt.Scan(&x)      //输入整数x
	flag := huiwen(x) //调用判断是否为回文的函数
	if flag {
		fmt.Println("true") //为true的话说明是回文，打印true
	} else {
		fmt.Println("false") //为false的话说明不是回文，打印false
	}
}
