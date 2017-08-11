package apns

import (
	"crypto/tls"
	"fmt"

	"testing"
	"time"
)

func TestClient(t *testing.T) {
	cert, err := tls.LoadX509KeyPair("/pem/iphone-cert-dev.pem", "/pem/iphone-key-dev.pem")
	if err != nil {
		panic(err)
	}
	apnsClient :=
		//NewClientWithProxy(cert, 0, 0, "127.0.0.1:1080").Development()
		NewClient(cert, 0).Development()
	aps := Aps{
		Alert: Alert{
			Title: "消息通知",
			Body:  "this is message",
		},
		Badge: 1,
	}
	payload := &Payload{Aps: aps, URL: "http://example.com", TaskID: "201708011"}

	apnsClient.BoundID = "this.app.package.name"
	resp, err := apnsClient.Push("239435909zbdnfioeirotufhdjfdkl", payload, time.Now().Unix())
	if err != nil {
		panic(err)
	}
	fmt.Println("StatusCode:", resp.StatusCode, "ApnsID:", resp.ApnsID, "Reason:", resp.Reason)
}
