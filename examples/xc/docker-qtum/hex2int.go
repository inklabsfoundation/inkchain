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
