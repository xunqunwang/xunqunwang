package main

import (
	"context"
	"flag"
	"go-online/app/domain/identify/conf"
	"go-online/app/domain/identify/server/gorpc"
	"go-online/lib/ecode/tip"
	"go-online/lib/syscall"
	"os"
	"os/signal"
	"time"

	// "go-online/app/domain/identify/server/grpc"
	"go-online/app/domain/identify/server/http"
	"go-online/app/domain/identify/service"
	"go-online/lib/log"

	// "go-online/lib/net/rpc/warden/resolver/livezk"
	// "go-online/lib/net/trace"

	"google.golang.org/grpc"
)

const (
// discoveryID = "passport.service.identify"
)

var (
	svr        *service.Service
	ws         *grpc.Server
	cancelFunc context.CancelFunc
)

func main() {
	var err error
	flag.Parse()
	// init conf,log,trace,stat,perf.
	if err = conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	// trace.Init(conf.Conf.Tracer)
	// defer trace.Close()
	// service init
	tip.Init(nil)
	svr = service.New(conf.Conf)
	http.Init(conf.Conf, svr)

	// init warden server
	ws, cancelFunc, err = gorpc.New(svr)
	if err != nil {
		panic(err)
	}
	// 先主站内部和chenzhihui测试可用，再对外提供
	// cancel, err := livezk.Register(conf.Conf.LiveZK, conf.Conf.WardenServer.Addr, discoveryID)
	// if err != nil {
	// 	panic(err)
	// }

	// signal handler
	log.Info("domain.identify start")
	signalHandler()
}

func signalHandler() {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		log.Info("identify-service get a signal %s", si.String())
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			// cancel()
			// ws.Shutdown(context.Background())
			ws.Stop()
			time.Sleep(time.Second * 2)
			svr.Close()
			log.Info("identify-service exit")
			cancelFunc()
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
