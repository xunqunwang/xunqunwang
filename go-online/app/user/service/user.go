package service

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"go-online/app/user/model"
	"go-online/lib/ecode"
	"go-online/lib/log"
	"go-online/lib/mail"
)

func (s *Service) SendVerificationCode(email string) (err error) {
	// check email
	user := new(model.User)
	if err = s.DB.Where("email = ?", email).First(user).Error; err != nil {
		log.Error("SendVerificationCode(%s) error(%v)", email, err)
		err = ecode.UserNotExist
		return
	}

	// generate verification code, which is valid for 10 minites.
	code := GenValidateCode(6)
	sT := time.Now()
	eT := sT.Add(10 * time.Minute)
	vcode := &model.VCode{
		Code:      code,
		UserId:    user.Id,
		CreatedAt: sT,
		ExpiredAt: eT,
	}
	if err = s.DB.Save(vcode).Error; err != nil {
		log.Error("SendVerificationCode(%s) error(%v)", email, err)
		return
	}

	// send verification code to user's email
	mailTo := []string{email}
	subject := "[xunqunwang] Please verify your device"
	body := fmt.Sprintf(`Hey %s!

A sign in attempt requires further verification because we did not recognize your device. To complete the sign in, enter the verification code on the unrecognized device.

Device: Chrome on Windows
Verification code: %s

Thanks,
The ThreePeople Team`, user.LoginName, code)

	// body := `Hey wkKaidy!

	// A sign in attempt requires further verification because we did not recognize your device. To complete the sign in, enter the verification code on the unrecognized device.

	// Device: Chrome on Windows
	// Verification code: 436291

	// Thanks,
	// The ThreePeople Team`
	if err = mail.SendMail(mailTo, subject, body); err != nil {
		log.Error("SendVerificationCode(%s) error(%v)", email, err)
		return
	}
	return
}

func (s *Service) ResetPassword(code, password string) (err error) {
	// check verification code
	vcode := new(model.VCode)
	if err = s.DB.Where("code = ?", code).First(vcode).Error; err != nil {
		log.Error("ResetPassword(%s,%s) error(%v)", code, password, err)
		err = ecode.CaptchaErr
		return
	}

	// modify user password
	if err = s.DB.Model(model.User{}).Where("id = ?", vcode.UserId).
		Update("password", password).Error; err != nil {
		log.Error("ResetPassword(%s,%s) error(%v)", code, password, err)
		return
	}
	return
}

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}
