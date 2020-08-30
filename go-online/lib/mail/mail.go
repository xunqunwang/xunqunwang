package mail

import (
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendMail(mailTo []string, subject string, body string) error {
	//定义邮箱服务器连接信息，如果是阿里邮箱 pass填密码，qq邮箱填授权码
	mailConn := map[string]string{
		"user": "wakai2009dy@126.com",
		"pass": "wk19890703",
		"host": "smtp.126.com", //邮箱服务器地址，如腾讯企业邮箱为smtp.exmail.qq.com
		"port": "25",           //邮箱服务器端口，如腾讯企业邮箱为465
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()
	// m.SetHeader("From", "XD Game"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo...)    //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)    //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err
}

// 阿里邮箱
// SMTP 服务器地址：smtp.mxhichina.com（端口：465/587）

// QQ 邮箱
// POP3 服务器地址：qq.com（端口：995）
// SMTP 服务器地址：smtp.qq.com（端口：465/587）

// 163 邮箱：
// POP3 服务器地址：pop.163.com（端口：110）
// SMTP 服务器地址：smtp.163.com（端口：25）

// 126 邮箱：
// POP3 服务器地址：pop.126.com（端口：110）
// SMTP 服务器地址：smtp.126.com（端口：25）
