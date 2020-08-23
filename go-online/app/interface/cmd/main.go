package main

import (
	"flag"
	"go-online/app/interface/di"
	"go-online/lib/conf/paladin"
	"go-online/lib/ecode/tip"
	"go-online/lib/log"
	"go-online/lib/net/trace/zipkin"
	"go-online/lib/os/signal"
	"go-online/lib/syscall"
	"os"
	"time"
)

var (
	closeFunc func()
)

func main() {
	var err error
	flag.Parse()
	log.Init(nil)
	defer log.Close()
	log.Info("admin start")
	zipkin.Init(nil)
	paladin.Init()
	tip.Init(nil)
	_, closeFunc, err = di.InitApp()
	if err != nil {
		panic(err)
	}
	signalHandler()
}

func signalHandler() {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		log.Info("get a signal %s", si.String())
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			closeFunc()
			log.Info("interface exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
