package myutils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/robfig/cron"
	"io/ioutil"
	"net/http"
)

/**
获取Token
*/
type Token struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
}

/**
报警内容
*/
type Content struct {
	Content string `json:"content"`
}

/**
企业微信发送文本消息
*/
type TextMessage struct {
	Touser  string  `json:"touser"`
	Toparty string  `json:"toparty"`
	Totag   string  `json:"totag"`
	Msgtype string  `json:"msgtype"`
	Agentid int     `json:"agentid"`
	Text    Content `json:"text"`
	Safe    int     `json:"safe"`
}

// 通过定时任务刷新token保存在这里
var t = &Token{}

// 企业ID
var corpid = beego.AppConfig.String("corpid")

// 应用的凭证密钥
var secret = beego.AppConfig.String("corpsecret")

// 企业应用的id，整型
var agentid, _ = beego.AppConfig.Int("agentid")

func init() {
	// 获取token的url
	tokenUrl := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + corpid + "&corpsecret=" + secret
	req := httplib.Get(tokenUrl)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	c := cron.New()
	// 每小时刷新一次token
	spec := "0 0 * * * ?"
	_ = c.AddFunc(spec, func() {
		err := req.ToJSON(t)
		if err != nil {
			println(err)
		}
		println(t.AccessToken)
	})
	// 服务启动后，手动刷新一次
	c.Start()
	err := req.ToJSON(t)
	if err != nil {
		println(err)
	}
	println(t.AccessToken)
}

func GetToken() string {
	return t.AccessToken
}

func SendMessageToWeChatApp(users string, content string) string {
	// 发送应用消息的url
	url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + GetToken()
	msg := &TextMessage{}
	msg.Touser = users
	msg.Agentid = agentid
	msg.Msgtype = "text"
	msg.Safe = 0
	msg.Text = Content{Content: content}
	byteArr, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return "send failed!"
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(byteArr))
	if err != nil {
		fmt.Println(err)
		return "send failed!"
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "send failed!"
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "send failed!"
	}
	result := string(body)
	fmt.Println("response Body:", result)
	return result
}
