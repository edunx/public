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

var Out *Output

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


func (o *Output) Prefix( info string ) {
    o.prefix = info
}

func (o *Output) Err(format string, v ...interface{}) {
	if o.level >= ERR {
		log.Printf(o.prefix + "[error] "+format+"\n", v...)
	}
}

func (o *Output) Info(format string, v ...interface{}) {
	if o.level >= INFO {
		log.Printf(o.prefix + "[info] "+format+"\n", v...)
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
		format = fmt.Sprintf(o.prefix + "[debug] %s %s:%d %s\n", funcName, filename, line, format)
		log.Printf(format, v...)
	}
}


// 守护日志文件,防止文件被删除,新日志丢失
func (o *Output) DaemonLog() {
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
	dir := filepath.Dir(o.path)
	err = watcher.Add(dir)
	if os.IsNotExist(err) {
		SetOutput(o.path, o.level)
		goto ADD
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Remove == fsnotify.Remove && event.Name == o.path {
				log.Println("log file was removed, will recreate it")
				SetOutput(o.path, o.level)
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

func SetOutput(path string, level int) *Output {
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
