package main

import (
    //"os"
    "fmt"
    "strconv"
)

func main(){
    input_hex := "00000000000000000000000000000000000000000000000000000000000001f4"//os.Args[1]
    hex := "0x" + input_hex
    int, _ := strconv.ParseUint(hex, 0, 64)
    fmt.Println(int)

}

func a(s string)  {
    fmt.Println(s)
}

//1）string 转 byte
//A => 0x41
//
//2）string 转 bytes
//A => 0x41
//fghj => 0x6667686a
//
//3）string 转 bytes32
//qtum => 0x7174756d00000000000000000000000000000000000000000000000000000000