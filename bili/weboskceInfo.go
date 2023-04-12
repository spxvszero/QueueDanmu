package bili

import (
	"encoding/json"
	"fmt"
	"github.com/monaco-io/request"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type BaseResp struct {
	Code      int64           `json:"code"`
	Message   string          `json:"message"`
	RequestId string          `json:"request_id"`
	Data      json.RawMessage `json:"data"`
}

type WebsocketInfoResp struct {
	Ip       []string `json:"ip"`
	AuthBody string   `json:"auth_body"`
	Host     []string `json:"host"`
	TcpPort  []int64  `json:"tcp_port"`
	WsPort   []int64  `json:"ws_port"`
	WssPort  []int64  `json:"wss_port"`
}

type WebsocketInfoReq struct {
	RoomId int64 `json:"room_id"`
}

func GetWebsocketInfo(roomId int64, akId, akSecret string) (host string, port int64, authBody string, err error) {
	req := &WebsocketInfoReq{
		RoomId: roomId,
	}
	content, _ := json.Marshal(req)
	header := &CommonHeader{
		ContentType:       JsonType,
		ContentAcceptType: JsonType,
		Timestamp:         strconv.FormatInt(time.Now().Unix(), 10),
		SignatureMethod:   HmacSha256,
		SignatureVersion:  BiliVersion,
		Authorization:     "",
		Nonce:             strconv.FormatInt(time.Now().UnixNano(), 10),
		AccessKeyId:       akId,
		ContentMD5:        Md5(string(content)),
	}
	header.Authorization = CreateSignature(header, akSecret)
	cli := request.Client{
		Method: "POST",
		URL:    fmt.Sprintf("%s/v1/common/websocketInfo", TestHttpHost),
		Header: header.ToMap(),
		String: string(content),
	}
	fmt.Println(cli.Header, cli.String)
	result := BaseResp{}
	resp := cli.Send().Scan(&result)
	if !resp.OK() {
		err = fmt.Errorf("[websocketInfo | GetWebsocketInfo] req:%+v resp:%+v", req, resp)
		return
	}
	if result.Code != 0 {
		err = fmt.Errorf("[websocketInfo | GetWebsocketInfo] result.Code req:%+v result:%+v", req, result)
		return
	}
	infoData := &WebsocketInfoResp{}
	fmt.Println(result)
	err = json.Unmarshal(result.Data, infoData)
	if err != nil {
		err = errors.Wrapf(err, "[websocketInfo | GetWebsocketInfo] json.Unmarshal err req:%+v resp:%+v", req, resp)
		return
	}
	// 这里简单写下，实际可以根据需求取
	if len(infoData.WsPort) != 0 && len(infoData.Host) != 0 {
		host, port, authBody = infoData.Host[0], infoData.WsPort[0], infoData.AuthBody
	}
	return
}
