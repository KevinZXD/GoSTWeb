
package alert

import (
"bytes"
"encoding/json"
"fmt"
"net"
"net/http"
"time"
)

var (
	alertCH chan *alertMsgSt
	localIP string
)

const ALERT_URL = "http://aladdin.biz.weibo.com/v2/api/adalert_custom/"
const title = "test title"
const users = "1067892503@qq.com"
const sendType = "all" // wechat and mail

type alertMsgSt struct {
	aggregationKey string
	errMsg         string
	count          int
}

type bodySt struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	User     string `json:"user"`
	SendType string `json:"send_type"`
}

// 调用阿拉丁报警接口
func alert(in *alertMsgSt) {

	body := bodySt{
		Title:    title,
		User:     users,
		SendType: sendType,
		Message:  fmt.Sprintf("\n机器地址:%s\n触发次数:%d\n报警内容:%s\n", localIP, in.count, in.errMsg),
	}

	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", ALERT_URL, bytes.NewBuffer(bodyBytes))
	if err != nil || req == nil {
		fmt.Println("http request alert api failed , err:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("http request alert api failed , err:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(resp)
	}

	return
}

func Alert(aggregationKey, errMsg string) {
	if len(alertCH) == cap(alertCH) {
		fmt.Println("alert chan have full")
		return
	}

	alertCH <- &alertMsgSt{
		aggregationKey: aggregationKey,
		errMsg:         errMsg,
	}

}

func init() {
	localIP, _ = getLocalIp()
	alertCH = make(chan *alertMsgSt, 100)
	alertMap := make(map[string]*alertMsgSt)

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for {
			select {
			case <-ticker.C: // 5 分钟触发一次报警
				for _, v := range alertMap {
					alert(v)
				}
				alertMap = make(map[string]*alertMsgSt) //清空已经聚合的结果

			case v := <-alertCH: // 聚合
				if _, ok := alertMap[v.aggregationKey]; !ok {
					alertMap[v.aggregationKey] = &alertMsgSt{
						aggregationKey: v.aggregationKey,
						errMsg:         v.errMsg,
						count:          0,
					}
				}
				alertMap[v.aggregationKey].count++
			}
		}
	}()

}

// getLocalIp 获取本机出口IPv4的字符串形式
func getLocalIp() (string, error) {
	conn, err := net.DialTimeout("udp", "10.255.255.255:80", 10*time.Millisecond)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.To4().String(), nil
}

