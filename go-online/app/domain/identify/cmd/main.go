package main

import (
	"flag"
	"go-online/lib/ecode/tip"
	"go-online/lib/syscall"
	"os"
	"os/signal"
	"time"

	"go-online/lib/conf/paladin"
	"go-online/lib/log"

	"go-online/app/domain/identify/di"
	"go-online/lib/net/trace/zipkin"
)

var (
	closeFunc func()
)

func main() {
	var err error
	flag.Parse()
	// init conf,log,trace,stat,perf.
	log.Init(nil)
	defer log.Close()
	log.Info("domain.identify start")
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
		log.Info("domain.identify get a signal %s", si.String())
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			closeFunc()
			time.Sleep(time.Second * 2)
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
