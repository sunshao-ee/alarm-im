package controllers

import (
	"alarm-im/myutils"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
)

type ImController struct {
	beego.Controller
}

func (c *ImController) Post() {
	// 获取im用户名，如: test01,test02
	tos := c.GetString("tos")
	// 微信api要求多个用户之间使用‘|’分隔
	tos = strings.Join(strings.Split(tos, ","), "|")
	// 获取报警内容，如：
	// [P0]		[PROBLEM]	[home-host][]	[1999端口挂了 	all(#3) net.port.listen 	port=1999 	0==0]	[O3 2019-01-16 	16:49:00]
	// Priority	Status		Endpoint		[Note			Func	Metric				SortedTags]			[CurrentStep	FormattedTime]
	content := c.GetString("content")
	trim := strings.Trim(strings.Trim(content, "["), "]")
	fmt.Printf("\ntrim:%s\n", trim)
	split := strings.Split(trim, "][")
	fmt.Printf("\nsplit %s %d\n", strings.Join(split, "|"), len(split))

	if len(split) == 6 {
		level := split[0]
		status := split[1]
		if status != "OK" {
			status = "故障消息"
		} else {
			status = "警告消息"
		}
		endpoint := split[2]
		noteSplit := strings.Split(split[4], " ")
		note := noteSplit[0]
		num := strings.Replace(strings.Split(split[5], " ")[0], "O", "", 1)
		time := strings.Split(split[5], " ")[1] + " " + strings.Split(split[5], " ")[2]
		content = "级别：\t" + level + "\n" + "类型：\t" + status + "\n" + "节点：" + "\t" + endpoint + "\n" + "描述：" + "\t" + note + "\n" + "次数：" + "\t" + num + "\n" + "时间：" + "\t" + time
	}
	fmt.Println(content)
	result := myutils.SendMessageToWeChatApp(tos, content)
	_, err := c.Ctx.ResponseWriter.Write([]byte(result))
	if err != nil {
		fmt.Println(err)
	}
}
