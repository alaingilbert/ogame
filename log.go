package ogame

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

// Quiet mode will not show any informative output
func (b *OGame) Quiet(quiet bool) {
	b.quiet = quiet
}

// SetLogger set a custom logger for the bot
func (b *OGame) SetLogger(logger *log.Logger) {
	b.logger = logger
}

// Terminal styling constants
const (
	knrm = "\x1B[0m"
	kred = "\x1B[31m"
	//kgrn = "\x1B[32m"
	kyel = "\x1B[33m"
	//kblu = "\x1B[34m"
	kmag = "\x1B[35m"
	kcyn = "\x1B[36m"
	kwht = "\x1B[37m"
)

func (b *OGame) log(prefix, color string, v ...interface{}) {
	if !b.quiet {
		_, f, l, _ := runtime.Caller(2)
		args := append([]interface{}{fmt.Sprintf(color+"%s"+knrm+" [%s:%d]", prefix, filepath.Base(f), l)}, v...)
		b.logger.Println(args...)
	}
}

func (b *OGame) trace(v ...interface{}) {
	b.log("TRAC", kwht, v...)
}

func (b *OGame) info(v ...interface{}) {
	b.log("INFO", kcyn, v...)
}

func (b *OGame) warn(v ...interface{}) {
	b.log("WARN", kyel, v...)
}

func (b *OGame) error(v ...interface{}) {
	b.log("ERRO", kred, v...)
}

func (b *OGame) critical(v ...interface{}) {
	b.log("CRIT", kred, v...)
}

func (b *OGame) debug(v ...interface{}) {
	b.log("DEBU", kmag, v...)
}

func (b *OGame) println(v ...interface{}) {
	b.log("PRIN", kwht, v...)
}
