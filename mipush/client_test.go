package mipush

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestClient(t *testing.T) {
	xmm := &XMMessage{
		Payload:               "http://example.com",
		RestrictedPackageName: "push.app.package.name",
		PassThrough:           1, // 0 表示通知栏消息1 表示透传消息
		Title:                 "消息通知",
		Description:           "this is message",
		NotifyType:            1,
		TaskId:                "11",
	}

	// 设置是否被覆盖，不同的数字，可显示多行
	s := "20170811"
	xmm.SetNotifyId(s[0:10])
	xmm.SetUserAccount("139320")
	client := NewClient("this is secret")
	client.SetProductionUrl(User_Account_Url)
	//
	//	xmm.SetRegId(".........................fGU5js=")
	//	client := NewClient("this is secret")
	//	client.SetProductionUrl(Reg_url)
	resp, err := client.Push(xmm)
	if err != nil {
		fmt.Println(err)
	}
	if resp.Result == RESULT_OK {
		tt := strings.Split(resp.Info, " ")
		if len(tt) == 6 {
			m, _ := strconv.Atoi(tt[4])
			fmt.Println(m + 1)
		}
	}
	fmt.Printf("%+v", resp)
}
