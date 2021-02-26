package public

import (
	"fmt"
	"github.com/edunx/lua"
	"net"
	"math"
	"os"
	"reflect"
	"unsafe"
)

func CheckUserData(L *lua.LState, idx int) UserData {
	v := L.ToUserData(idx)

	switch ud := v.Value.(type) {
	case UserData:
		return ud
	default:
		L.RaiseError("must common.Userdata , got %v", ud)
		return nil
	}
}

func Round(val float64, precision int64) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(precision))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= 0.5 {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func GetLocalIP() string {
    addresses, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Printf("get local ip error: %v\n", err)
		os.Exit(1)
	}

	for _, address := range addresses {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}
// B2S converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// S2B converts string to a byte slice without memory allocation.
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func S2B(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return
}