// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"html/template"
	"net/smtp"
	"strings"
	"time"

	"bytes"
	. "config"
	"logger"
	"model"
	"util"
)

// 发送电子邮件功能
func SendMail(subject, content string, tos []string) error {
	message := `From: Go语言中文网 | Golang中文社区 | Go语言学习园地<` + Config["from_email"] + `>
To: ` + strings.Join(tos, ",") + `
Subject: ` + subject + `
Content-Type: text/html;charset=UTF-8

` + content

	auth := smtp.PlainAuth("", Config["smtp_username"], Config["smtp_password"], Config["smtp_host"])
	err := smtp.SendMail(Config["smtp_addr"], auth, Config["from_email"], tos, []byte(message))
	if err != nil {
		logger.Errorln("Send Mail to", strings.Join(tos, ","), "error:", err)
		return err
	}
	logger.Infoln("Send Mail to", strings.Join(tos, ","), "Successfully")
	return nil
}

// 自定义模板函数
var emailFuncMap = template.FuncMap{
	"time_format": func(s string) string {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
		if err != nil {
			return s
		}

		return t.Format("01-02")
	},
	"substring": util.Substring,
}

var emailTpl = template.Must(template.New("email.html").Funcs(emailFuncMap).ParseFiles(ROOT + "/template/email.html"))

// 订阅邮件通知
func EmailNotice() {

	beginTime := time.Now().Add(-7*24*time.Hour).Format("2006-01-02") + " 00:00:00"

	// 本周晨读（过去 7 天）
	readings, err := model.NewMorningReading().Where("ctime>? AND rtype=0", beginTime).Order("id DESC").FindAll()
	if err != nil {
		logger.Errorln("find morning reading error:", err)
	}

	// 本周精彩文章
	articles, err := model.NewArticle().Where("ctime>? AND status!=2", beginTime).Order("cmtnum DESC, likenum DESC, viewnum DESC").Limit("10").FindAll()
	if err != nil {
		logger.Errorln("find article error:", err)
	}

	// 本周热门主题
	topics, err := model.NewTopic().Where("ctime>? AND flag IN(0,1)", beginTime).Order("tid DESC").Limit("10").FindAll()
	if err != nil {
		logger.Errorln("find topic error:", err)
	}

	data := map[string]interface{}{
		"readings": readings,
		"articles": articles,
		"topics":   topics,
	}

	// 给所有用户发送邮件
	userModel := model.NewUser()

	var (
		lastUid = 0
		users   []*model.User
	)

	for {
		users, err = userModel.Where("uid>?", lastUid).Order("uid ASC").FindAll()
		if err != nil {
			logger.Errorln("find user error:", err)
			continue
		}

		if len(users) == 0 {
			break
		}

		for _, user := range users {
			if user.Unsubscribe == 1 {
				logger.Infoln(user)
				continue
			}

			data["email"] = user.Email
			data["token"] = GenUnsubscribeToken(user.Username, user.Email)

			content, err := genEmailContent(data)
			if err != nil {
				logger.Errorln("from email.html gen email content error:", err)
				continue
			}

			go func(content, email string) {
				SendMail("每周精选", content, []string{"274768166@qq.com"})
			}(content, user.Email)
			break
		}
		break
	}

}

// 生成 退订 邮件的 token
func GenUnsubscribeToken(username, email string) string {
	return util.Md5(username + email + model.UNSUBSCRIBE_TOKEN_KEY)
}

func genEmailContent(data map[string]interface{}) (string, error) {
	buffer := &bytes.Buffer{}
	if err := emailTpl.Execute(buffer, data); err != nil {
		logger.Errorln("execute template error:", err)
		return "", err
	}

	return buffer.String(), nil
}
