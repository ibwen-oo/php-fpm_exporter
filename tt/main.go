package main

import "fmt"

func main() {
	b := []byte{70, 105, 108, 101, 32, 110, 111, 116, 32, 102, 111, 117, 110, 100, 46, 10}
	//b := bytes.Buffer{}[70 105 108 101 32 110 111 116 32 102 111 117 110 100 46 10]
	fmt.Printf("%s", b)
}
