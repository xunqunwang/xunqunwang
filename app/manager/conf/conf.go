package conf

import (
	"errors"
	"flag"

	"go-online/lib/cache/memcache"
	"go-online/lib/conf"
	"go-online/lib/database/orm"
	"go-online/lib/log"
	bm "go-online/lib/net/http/blademaster"
	"go-online/lib/net/http/blademaster/middleware/permit"
	"go-online/lib/net/rpc/warden"
	"go-online/lib/net/trace"
	xtime "go-online/lib/time"

	"github.com/BurntSushi/toml"
)

// Config .
type Config struct {
	Cfg          *cfg
	App          *bm.App
	ORM          *orm.Config
	Log          *log.Config
	Tracer       *trace.Config
	Memcache     *memcache.Config
	HTTPServer   *bm.ServerConfig
	HTTPClient   *bm.ClientConfig
	DsbClient    *bm.ClientConfig
	UnameTicker  xtime.Duration
	WardenServer *warden.ServerConfig
	Permit       *permit.Config2
}

type cfg struct {
	RankGroupMaxPs int
}

var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

// Init .
func Init() (err error) {
	if confPath != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

func remote() (err error) {
	if client, err = conf.New(); err != nil {
		return
	}
	if err = load(); err != nil {
		return
	}
	go func() {
		for range client.Event() {
			log.Info("config reload")
			if load() != nil {
				log.Error("config reload err")
			}
		}
	}()
	return
}

func load() (err error) {
	var (
		s       string
		ok      bool
		tmpConf *Config
	)
	if s, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
