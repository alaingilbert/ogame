package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

type logging struct {
	quiet  bool
	logger *log.Logger
}

// Quiet mode will not show any informative output
func (b *logging) Quiet(quiet bool) {
	b.quiet = quiet
}

// SetLogger set a custom logger for the bot
func (b *logging) SetLogger(logger *log.Logger) {
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

func (b *logging) log(prefix, color string, v ...interface{}) {
	if !b.quiet {
		_, f, l, _ := runtime.Caller(2)
		args := append([]interface{}{fmt.Sprintf(color+"%s"+knrm+" [%s:%d]", prefix, filepath.Base(f), l)}, v...)
		b.logger.Println(args...)
	}
}

func (b *logging) trace(v ...interface{}) {
	b.log("TRAC", kwht, v...)
}

func (b *logging) info(v ...interface{}) {
	b.log("INFO", kcyn, v...)
}

func (b *logging) warn(v ...interface{}) {
	b.log("WARN", kyel, v...)
}

func (b *logging) error(v ...interface{}) {
	b.log("ERRO", kred, v...)
}

func (b *logging) critical(v ...interface{}) {
	b.log("CRIT", kred, v...)
}

func (b *logging) debug(v ...interface{}) {
	b.log("DEBU", kmag, v...)
}

func (b *logging) println(v ...interface{}) {
	b.log("PRIN", kwht, v...)
}

func (b *logging) Fatal(v ...interface{}) {
	b.log("FATAL", kred, v...)
}
