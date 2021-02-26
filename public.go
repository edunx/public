package public

import (
	"github.com/edunx/lua"
	"unsafe"
)

var (
    VM  *lua.LState
    Out *Output
)

const (
	ERR   = 14
	INFO  = 16
	DEBUG = 18
)

type Output struct {
    prefix string
	path   string
	level  int
}

//自定义对象
type UserData interface { //用户再次读取修改配置的时候 需要转化成userdata
	ToUserData(*lua.LState) *lua.LUserData
}

type Logger interface {
	Err(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type Any struct {
	V        interface{}
	Handler  unsafe.Pointer
}
