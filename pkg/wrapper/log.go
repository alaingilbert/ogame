package wrapper

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

func (b *OGame) log(prefix, color string, v ...any) {
	if !b.quiet {
		_, f, l, _ := runtime.Caller(2)
		args := append([]any{fmt.Sprintf(color+"%s"+knrm+" [%s:%d]", prefix, filepath.Base(f), l)}, v...)
		b.logger.Println(args...)
	}
}

func (b *OGame) trace(v ...any) {
	b.log("TRAC", kwht, v...)
}

func (b *OGame) info(v ...any) {
	b.log("INFO", kcyn, v...)
}

func (b *OGame) warn(v ...any) {
	b.log("WARN", kyel, v...)
}

func (b *OGame) error(v ...any) {
	b.log("ERRO", kred, v...)
}

func (b *OGame) critical(v ...any) {
	b.log("CRIT", kred, v...)
}

func (b *OGame) debug(v ...any) {
	b.log("DEBU", kmag, v...)
}

func (b *OGame) println(v ...any) {
	b.log("PRIN", kwht, v...)
}
