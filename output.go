package public 

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	ERR   = 14
	INFO  = 16
	DEBUG = 18
)

type Output struct {
	path  string
	level int
}

func (o *Output) Err(format string, v ...interface{}) {
	if o.level >= ERR {
		log.Printf(C.ID+" "+C.Node+" [error] "+format+"\n", v...)
	}
}

func (o *Output) Info(format string, v ...interface{}) {
	if o.level >= INFO {
		log.Printf(C.ID+" "+C.Node+" [info] "+format+"\n", v...)
	}
}

func (o *Output) Debug(format string, v ...interface{}) {
	if o.level >= DEBUG {
		filename, line, funcName := "???", 0, "???"
		pc, filename, line, ok := runtime.Caller(1)
		if ok {
			funcName = runtime.FuncForPC(pc).Name()
			filename = filepath.Base(filename)
		}
		format = fmt.Sprintf(C.ID+" "+C.Node+" [debug] %s %s:%d %s\n", funcName, filename, line, format)
		log.Printf(format, v...)
	}
}

func SetOutput(path string, level int) Logger {
	out := &Output{path: path, level: level}

	logfile, err := os.OpenFile(out.path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("set logger output file failed: %v", err)
		return nil
	}
	log.SetOutput(io.MultiWriter(logfile, os.Stderr))
	Out = out
	return out
}

// 守护日志文件,防止文件被删除,新日志丢失
func DaemonLog(path string, level int) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("new watcher error: %v\n", err)
		return
	}

	defer func() {
		if err := watcher.Close(); err != nil {
			log.Printf("close watcher error: %s\n", err)
		}
	}()

ADD:
	dir := filepath.Dir(path)
	err = watcher.Add(dir)
	if os.IsNotExist(err) {
		SetOutput(path, level)
		goto ADD
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Remove == fsnotify.Remove && event.Name == path {
				log.Println("log file was removed, will recreate it")
				SetOutput(path, level)
				goto ADD
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}
