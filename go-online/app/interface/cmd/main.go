package main

import (
	"flag"
	"go-online/app/interface/conf"
	"go-online/app/interface/http"
	"go-online/app/interface/service"
	"go-online/lib/ecode/tip"
	"go-online/lib/log"
	"go-online/lib/net/trace"
	"go-online/lib/os/signal"
	"go-online/lib/syscall"
	"os"
	"time"
)

var (
	s *service.Service
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	tip.Init(nil)
	s = service.New(conf.Conf)
	http.Init(conf.Conf, s)
	log.Info("admin start")
	signalHandler()
}

func signalHandler() {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			time.Sleep(time.Second * 2)
			log.Info("get a signal %s, stop the admin process", si.String())
			s.Close()
			s.Wait()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
