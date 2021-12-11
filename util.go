package main

import (
	"encoding/base64"
	"fmt"
	"github.com/dop251/goja"
	"math/rand"
	"time"
)

func Print(message string) {
	fmt.Printf("%s\n", message)
}

func Sleep(ms uint) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func StringToByteArray(value string) []byte {
	return []byte(value)
}

func ByteArrayToString(value []byte) string {
	return string(value)
}

func RandomInt(min, max int) int {
	return rand.Int()%(max-min+1) + min
}

func RegisterToGojaVM(vm *goja.Runtime) {
	_ = vm.Set("Print", Print)
	_ = vm.Set("Sleep", Sleep)
	_ = vm.Set("StringToByteArray", StringToByteArray)
	_ = vm.Set("ByteArrayToString", ByteArrayToString)
	_ = vm.Set("RandomInt", RandomInt)
	_ = vm.Set("RandomBase64Bytes", RandomBase64Bytes)
}

func RandomBase64Bytes(n int) string {
	data := make([]byte, n)
	for i := 0; i < n; i++ {
		data[i] = byte(RandomInt(0, 255))
	}
	return base64.StdEncoding.EncodeToString(data)
}
