package util

import (
	"encoding/base64"
	"fmt"
	"github.com/dop251/goja"
	"math/rand"
	"time"
)

// Print prints any string
func Print(message string) {
	fmt.Printf("%s\n", message)
}

// Sleep delays the execution by ms milliseconds
func Sleep(ms uint) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// StringToByteArray converts a string to a byte slice
func StringToByteArray(value string) []byte {
	return []byte(value)
}

// ByteArrayToString converts a byte slice to string. This is sometimes important for javascript code
func ByteArrayToString(value []byte) string {
	return string(value)
}

// RandomInt generates a random int in the interval I=[min, max]
// This is not to be confused with [min, max[
func RandomInt(min, max int) int {
	return rand.Int()%(max-min+1) + min
}

// RegisterToGojaVM registeres all utility functions that are important to the vm
func RegisterToGojaVM(vm *goja.Runtime) {
	_ = vm.Set("Print", Print)
	_ = vm.Set("Sleep", Sleep)
	_ = vm.Set("StringToByteArray", StringToByteArray)
	_ = vm.Set("ByteArrayToString", ByteArrayToString)
	_ = vm.Set("RandomInt", RandomInt)
	_ = vm.Set("RandomBase64Bytes", RandomBase64Bytes)
}

// RandomBase64Bytes generates a byte slice with random values and transforms it to a base64 encoded string
func RandomBase64Bytes(n int) string {
	data := make([]byte, n)
	for i := 0; i < n; i++ {
		data[i] = byte(RandomInt(0, 255))
	}
	return base64.StdEncoding.EncodeToString(data)
}
