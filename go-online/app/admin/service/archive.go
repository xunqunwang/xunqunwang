package service

// import (
// 	"context"

// 	lmdl "go-online/app/admin/model"
// 	arcmdl "go-online/app/archive/api"
// 	// "go-online/lib/log"
// )

// // Archives get achives info .
// func (s *Service) Archives(c context.Context, p *lmdl.ArchiveParam) (res map[int64]*arcmdl.Arc, err error) {
// 	var (
// 		arcs *arcmdl.ArcsReply
// 	)
// 	// if arcs, err = s.arcClient.Arcs(c, &arcmdl.ArcsRequest{Aids: p.Aids}); err != nil {
// 	// 	log.Error("s.arcClient.Archives3(%v) error(%v)", p.Aids, err)
// 	// 	return
// 	// }
// 	res = make(map[int64]*arcmdl.Arc, len(p.Aids))
// 	for _, aid := range p.Aids {
// 		if arc, ok := arcs.Arcs[aid]; ok && arc.IsNormal() {
// 			res[aid] = arc
// 		}
// 	}
// 	return
// }
