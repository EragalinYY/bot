package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, _ := net.Dial("tcp", "golang.org:80")         //Соединение по TCP
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")         //Отправка строки через
	status, _ := bufio.NewReader(conn).ReadString('\n') //Вывод первой строки ответа
	fmt.Println(status)
}

var a int8
var b int16
var c int32
var d int64
