package mipush

import (
	"encoding/json"
	"fmt"
	"imserver/lib/log"
	"net/url"
)

// NOTE  define reference struct http://dev.xiaomi.com/doc/?p=533
type XMMessage struct {
	Payload               string //	消息的内容。
	RestrictedPackageName string //    App的包名。备注：V2版本支持一个包名，V3版本支持多包名（中间用逗号分割）。
	PassThrough           int8   //	pass_through的值可以为： 0 表示通知栏消息1 表示透传消息
	Title                 string //	通知栏展示的通知的标题。
	Description           string //	通知栏展示的通知的描述。

	// notify_type的值可以是DEFAULT_ALL或者以下其他几种的OR组合：DEFAULT_ALL = -1;
	// DEFAULT_SOUND  = 1; 使用默认提示音提示；DEFAULT_VIBRATE = 2; 使用默认震动提示
	// DEFAULT_LIGHTS = 4;     使用默认led灯光提示；
	NotifyType int8
	TaskId     string     // 上报数据使用
	xmuv       url.Values // 含有本条消息所有属性的数组
}

func (xm *XMMessage) buildXMPostParam() {
	xmuv := url.Values{}
	xmuv.Set("payload", xm.Payload)
	xmuv.Set("restricted_package_name", xm.RestrictedPackageName)
	xmuv.Set("pass_through", fmt.Sprintf("%d", xm.PassThrough))
	xmuv.Set("title", xm.Title)
	xmuv.Set("description", xm.Description)
	xmuv.Set("notify_type", fmt.Sprintf("%d", xm.NotifyType))
	xmuv.Set("extra.task_id", xm.TaskId)
	xm.xmuv = xmuv
	return
}

// 可选项。默认情况下，通知栏只显示一条推送消息。如果通知栏要显示多条推送消息，需要针对不同的消息设置不同的notify_id（相同notify_id的通知栏消息会覆盖之前的）。
// 可选项 notify_id 0-4同一个notifyId在通知栏只会保留一条
func (xm *XMMessage) SetNotifyId(notifyId string) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	xm.xmuv.Set("notify_id", notifyId)
}

// 可选项。如果用户离线，设置消息在服务器保存的时间，单位：ms。服务器默认最长保留两周。
// time_to_live 可选项，当用户离线是，消息保留时间，默认两周，单位ms
func (xm *XMMessage) SetTimeToLive(timeToLive int64) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	xm.xmuv.Set("time_to_live", fmt.Sprintf("%d", timeToLive))
}

// 可选项。定时发送消息。用自1970年1月1日以来00:00:00.0 UTC时间表示（以毫秒为单位的时间）。注：仅支持七天内的定时消息。
func (xm *XMMessage) SetTimeToSend(timeToSend int64) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	xm.xmuv.Set("time_to_send", fmt.Sprintf("%d", timeToSend))
}

// 根据user_account，发送消息给设置了该user_account的所有设备。可以提供多个user_account，user_account之间用“,”分割。参数仅适用于“/message/user_account”HTTP API。
func (xm *XMMessage) SetUserAccount(UserAccount string) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	xm.xmuv.Set("user_account", UserAccount)
}

// 针对不同的userAccount推送不同的消息
// 根据user_accounts，发送消息给设置了该user_account的所有设备。可以提供多个user_account，user_account之间用“,”分割。
func (xm *XMMessage) SetUserAccounts(UserAccount string) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	xm.xmuv.Set("user_accounts", UserAccount)
}

//	根据registration_id，发送消息到指定设备上。可以提供多个registration_id，发送给一组设备，不同的registration_id之间用“,”分割。
func (xm *XMMessage) SetRegId(deviceToken string) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	xm.xmuv.Set("registration_id", deviceToken)
}

// 根据topic，发送消息给订阅了该topic的所有设备。参数仅适用于“/message/topic”HTTP API。
func (xm *XMMessage) SetTopic(UserAccount string) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	xm.xmuv.Set("topic", UserAccount)
}

// push return result
type Response struct {
	Result      string `json:"result,omitempty"`      //“result”: string，”ok” 表示成功, “error” 表示失败。
	Reason      string `json:"reason,omitempty"`      //reason: string，如果失败，reason失败原因详情。
	Code        int    `json:"code,omitempty"`        //“code”: integer，0表示成功，非0表示失败。
	Data        Data   `json:"data,omitempty"`        //“data”: string，本身就是一个json字符串（其中id字段的值就是消息的Id）。
	Description string `json:"description,omitempty"` //“description”: string， 对发送消息失败原因的解释。
	Info        string `json:"info,omitempty"`        //“info”: string，详细信息。
}

func (resp Response) String() string {
	return fmt.Sprintf("Result:%s,Reason:%s,Code:%d,Date:%s,Description:%s,Info:%s", resp.Result, resp.Reason, resp.Code, resp.Data, resp.Description, resp.Info)
}

type Data struct {
	Id string `json:"id,omitempty"`
}

// Unmarshal push return result body to Response
func (resp *Response) Unmarshal(body []byte) (err error) {
	err = json.Unmarshal(body, resp)
	if err != nil {
		log.Error("json.Unmarshal()body:%s,error:%v", string(body), err)
		return err
	}
	return nil
}
