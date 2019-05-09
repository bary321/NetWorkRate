package NetWorkRate

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"time"
)

type CustomLogger struct {
	lumberjack.Logger
	Switch bool
}

func ExampleLogger(filename string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    1,
		MaxBackups: 0,
		MaxAge:     10,
	}
}

func NewCustomLogger(filename string, sw bool) *CustomLogger {
	return &CustomLogger{*ExampleLogger(filename), sw}
}

func (c *CustomLogger) Println(i *IORates) {
	if c.Switch {
		t := time.Now().UTC()
		i.Time = t.Format("2006-01-02T15:04:05.000")
		_, err := fmt.Fprintf(&c.Logger, "%v\n", i)
		if err != nil {
			fmt.Println(err)
		}
	}
}
