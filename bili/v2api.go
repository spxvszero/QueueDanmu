package bili

import (
	"encoding/json"
	"fmt"
	"github.com/monaco-io/request"
	"strconv"
	"time"
	_ "unsafe"
)

//用户信息
type UserInfo struct {
	Uid		int64	`json:"uid"`
	Uname	string	`json:"uname"`
	Uface	string	`json:"uface"`
}

type Bili_BaseResp struct {
	Code 	int `json:"code"`
	Message string `json:"message"`
	Data 	json.RawMessage `json:"data"`
}

//  场次信息
type GameInfo struct {
	//  场次id,心跳key(心跳保持20s-60s)调用一次,超过60s无心跳自动关闭,长连停止推送消息
	GameId 	string `json:"game_id"`
}

//  长连信息
type WebsocketInfo struct {
	//  长连使用的请求json体 第三方无需关注内容,建立长连时使用即可
	AuthBody	string		`json:"auth_body"`
	//  wss 长连地址
	WssLink		[]string	`json:"wss_link"`
}

//  主播信息
type AnchorInfo struct {
	// 主播房间号
	RoomId 	int 	`json:"room_id"`
	// 主播昵称
	Uname 	string `json:"uname"`
	// 主播头像
	Uface 	string `json:"uface"`
	// 主播uid
	Uid 	int		`json:"uid"`
}

type AppStartReq struct {
	Code 	string	`json:"code"`
	AppId 	int64	`json:"app_id"`
}

type AppStartResp struct {
	GameInfo 		GameInfo 		`json:"game_info"`
	WebsocketInfo 	WebsocketInfo 	`json:"websocket_info"`
	AnchorInfo		AnchorInfo 		`json:"anchor_info"`
}

type AppGameHeartBreakReq struct {
	GameId 	string	`json:"game_id"`
}
type AppEndReq struct {
	AppId 	int64	`json:"app_id"`
	GameId 	string	`json:"game_id"`
}

/*
项目开启
接口描述：开启项目第一步，平台会根据入参进行鉴权校验。鉴权通过后，返回长连信息、场次信息和主播信息。开发者拿到长连和心跳信息后，需要参照长连说明和项目心跳，与平台保持健康的 心跳机制。
接口地址：/v2/app/start
方法：POST
请求参数：
参数名	必选	类型	描述
code	是	string	主播身份码
app_id	是	integer(13位长度的数值，注意不要用普通int，会溢出的)	项目ID
*/
func SendAppStartReq(code string, app_id int64) (*AppStartResp, error) {
	req := &AppStartReq{
		Code:  code,
		AppId: app_id,
	}
	content, _ := json.Marshal(req)
	header := generateHeaderWithContent(content)
	cli := request.Client{
		Method: "POST",
		URL:    fmt.Sprintf("%s/v2/app/start", TestHttpHost),
		Header: header.ToMap(),
		String: string(content),
	}
	fmt.Println(cli.Header, cli.String)
	result := Bili_BaseResp{}
	resp := cli.Send().Scan(&result)
	if !resp.OK() {
		return nil, fmt.Errorf("[BILI][GameStart] req:%+v resp:%+v\n", req, resp)
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("[BILI][GameStart] result.Code req:%+v result:%+v\n", req, result)
	}
	infoData := &AppStartResp{}
	fmt.Println(result)
	err := json.Unmarshal(result.Data, infoData)
	if err != nil {
		return nil, fmt.Errorf("[BILI][GameStart] json.Unmarshal err req:%+v resp:%+v\n", req, resp)
	}

	return infoData, nil
}


/*
项目关闭
接口描述：项目关闭时需要主动调用此接口，使用对应项目Id及项目开启时返回的game_id作为唯一标识，调用后会同步下线互动道具等内容，项目关闭后才能进行下一场次互动。
接口地址：/v2/app/end
方法：POST
请求参数：
参数名	必选	类型	描述
app_id	是	integer(13位长度的数值，注意不要用普通int，会溢出的)	项目id
game_id	是	string	场次id
*/
func SendAppEndReq(app_id int64, game_id string) error {
	req := &AppEndReq{
		AppId: app_id,
		GameId: game_id,
	}
	content, _ := json.Marshal(req)
	header := generateHeaderWithContent(content)
	cli := request.Client{
		Method: "POST",
		URL:    fmt.Sprintf("%s/v2/app/end", TestHttpHost),
		Header: header.ToMap(),
		String: string(content),
	}
	fmt.Println(cli.Header, cli.String)
	result := Bili_BaseResp{}
	resp := cli.Send().Scan(&result)
	if !resp.OK() {
		return fmt.Errorf("[BILI][GameEnd] req:%+v resp:%+v\n", req, resp)
	}
	if result.Code != 0 {
		return fmt.Errorf("[BILI][GameEnd] result.Code req:%+v result:%+v\n", req, result)
	}
	return nil
}


/*
项目心跳
接口描述：项目开启后，需要持续间隔20秒调用一次该接口。平台超过60s未收到项目心跳，会自动关闭当前场次（game_id），同时将道具相关功能下线，以确保下一场次项目正常运行。
接口地址：/v2/app/heartbeat
方法：POST
请求参数：
参数名	必选	类型	描述
game_id	是	string	场次id
*/
func SendAppGameHeartBreak(game_id string) error {
	req := &AppGameHeartBreakReq{
		GameId: game_id,
	}
	content, _ := json.Marshal(req)
	header := generateHeaderWithContent(content)
	cli := request.Client{
		Method: "POST",
		URL:    fmt.Sprintf("%s/v2/app/heartbeat", TestHttpHost),
		Header: header.ToMap(),
		String: string(content),
	}
	fmt.Println(cli.Header, cli.String)
	result := Bili_BaseResp{}
	resp := cli.Send().Scan(&result)
	if !resp.OK() {
		return fmt.Errorf("[BILI][GameHeartBreak] req:%+v resp:%+v\n", req, resp)
	}
	if result.Code != 0 {
		return fmt.Errorf("[BILI][GameHeartBreak] result.Code req:%+v result:%+v\n", req, result)
	}
	return nil
}

func generateHeaderWithContent(content []byte) *CommonHeader {
	header := &CommonHeader{
		ContentType:       JsonType,
		ContentAcceptType: JsonType,
		Timestamp:         strconv.FormatInt(time.Now().Unix(), 10),
		SignatureMethod:   HmacSha256,
		SignatureVersion:  BiliVersion,
		Authorization:     "",
		Nonce:             strconv.FormatInt(time.Now().UnixNano(), 10),
		AccessKeyId:       AccessKey,
		ContentMD5:        Md5(string(content)),
	}
	header.Authorization = CreateSignature(header, AccessSecret)
	return header
}