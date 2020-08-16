package grpc

import (
	pb "go-online/app/domain/identify/api"
	"go-online/app/domain/identify/dao"
	"go-online/lib/conf/paladin"
	"go-online/lib/net/rpc/warden"
)

// New new a grpc server.
func New(svc pb.IdentifyServer) (ws *warden.Server, cf func(), err error) {
	var (
		cfg warden.ServerConfig
		ct  paladin.TOML
	)
	if err = paladin.Get("grpc.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	ws = warden.NewServer(&cfg)
	pb.RegisterIdentifyServer(ws.Server(), svc)
	if ws, err = ws.Start(); err != nil {
		return
	}
	cf, err = dao.RegisterToConsul("grpc", cfg.Addr)
	return
}
