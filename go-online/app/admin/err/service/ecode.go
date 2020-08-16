package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"go-online/app/admin/err/model"
	"go-online/lib/ecode"
	"go-online/lib/log"
)

var (
	CurrentVer int64  = 1
	CurrentMd5 string = ""
)

// func (s *Service) GetEcodeList() (list []*model.Ecode, err error) {
// 	if err = s.dao.DB.Find(&list).Error; err != nil {
// 		log.Error("GetEcodeList error(%v)", err)
// 		return
// 	}
// 	return
// }

func (s *Service) GetEcodeList(ver int64) (data *model.EcodeData, err error) {
	var (
		list  []*model.Ecode
		bytes []byte
	)
	if ver == CurrentVer {
		err = ecode.NotModified
		return
	}
	if err = s.dao.DB.Find(&list).Error; err != nil {
		log.Error("GetEcodeList error(%v)", err)
		return
	}
	data = &model.EcodeData{Ver: CurrentVer}
	data.Code = make(map[int]string)
	for _, v := range list {
		data.Code[v.Code] = v.Message
	}
	if bytes, err = json.Marshal(data.Code); err != nil {
		return
	}
	mb := md5.Sum(bytes)
	md5 := hex.EncodeToString(mb[:])
	if CurrentMd5 != md5 {
		CurrentVer += 1
		CurrentMd5 = md5
	}
	data.MD5 = md5
	return
}
