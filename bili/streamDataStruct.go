package bili

import "encoding/json"

const (
	//弹幕
	CMD_Danmu = "LIVE_OPEN_PLATFORM_DM"
	//礼物
	CMD_Liwu = "LIVE_OPEN_PLATFORM_SEND_GIFT"
	//付费留言
	CMD_Liuyan = "LIVE_OPEN_PLATFORM_SUPER_CHAT"
	//付费留言下线
	CMD_LiuyanOff = "LIVE_OPEN_PLATFORM_SUPER_CHAT_DEL"
	//付费大航海
	CMD_Dahanghai = "LIVE_OPEN_PLATFORM_GUARD"
)

type SocketBase struct {
	Cmd		string	`json:"cmd"`
	Data	json.RawMessage	`json:"data"`
}

/*
字段名	类型	描述
uname	string	用户昵称
uid	int64	用户UID
uface	string	用户头像
timestamp	int64	弹幕发送时间秒级时间戳
room_id	int64	弹幕接收的直播间
msg	string	弹幕内容
msg_id	string	消息唯一id
guard_level	int64	对应房间大航海等级 1总督 2提督 3舰长
fans_medal_wearing_status	bool	该房间粉丝勋章佩戴情况
fans_medal_name	string	粉丝勋章名
fans_medal_level	int64	对应房间勋章信息
emoji_img_url	string	表情包图片地址
dm_type	int64	弹幕类型 0：普通弹幕 1：表情包弹幕
*/
type Socket_Danmu struct {
	Uname 	string	`json:"uname"`
	Uid 	int64	`json:"uid"`
	Uface	string	`json:"uface"`
	Timestamp int64	`json:"timestamp"`
	RoomId 	int64	`json:"room_id"`
	Msg		string	`json:"msg"`
	GuardLevel	int64	`json:"guard_level"`
	FansMedalWearingStatus	bool	`json:"fans_medal_wearing_status"`
	FansMedalName	string	`json:"fans_medal_name"`
	FansMedalLevel	int64	`json:"fans_medal_level"`
	EmojiImgUrl	string	`json:"emoji_img_url"`
	DmType		int64	`json:"dm_type"`
}

/*
字段名	类型	描述
room_id	int64	房间号
uid	int64	送礼用户UID
uname	string	送礼用户昵称
uface	string	送礼用户头像
gift_id	int64	道具id(盲盒:爆出道具id)
gift_name	string	道具名(盲盒:爆出道具名)
gift_num	int64	赠送道具数量
price	int64	支付金额(1000 = 1元 = 10电池),盲盒:爆出道具的价值
paid	bool	是否是付费道具
fans_medal_level	int64	实际送礼人的勋章信息
fans_medal_name	string	粉丝勋章名
fans_medal_wearing_status	bool	该房间粉丝勋章佩戴情况
guard_level	int64	大航海等级
timestamp	int64	收礼时间秒级时间戳
anchor_info	anchor_info结构体	主播信息
msg_id	string	消息唯一id
gift_icon	string	道具icon
combo_gift	bool	是否是combo道具
combo_info	combo_info结构体	连击信息
*/
type Socket_Liwu struct {
	RoomId	int64	`json:"room_id"`
	Uid		int64	`json:"uid"`
	Uname	string	`json:"uname"`
	Uface	string	`json:"uface"`
	GiftId	int64	`json:"gift_id"`
	GiftName string	`json:"gift_name"`
	GiftNum	int64	`json:"gift_num"`
	Price	int64	`json:"price"`
	Paid	bool	`json:"paid"`
	FansMedalLevel	int64	`json:"fans_medal_level"`
	FansMedalName	string	`json:"fans_medal_name"`
	FansMedalWearingStatus	bool	`json:"fans_medal_wearing_status"`
	GuardLevel	int64	`json:"guard_level"`
	Timestamp	int64	`json:"timestamp"`
	AnchorInfo	AnchorInfo	`json:"anchor_info"`
	MsgId		string	`json:"msg_id"`
	GiftIcon	string	`json:"gift_icon"`
	ComboGift	bool	`json:"combo_gift"`
	ComboInfo	ComboInfo	`json:"combo_info"`
}

/*
combo_info	类型	描述
combo_base_num	int64	每次连击赠送的道具数量
combo_count	int64	连击次数
combo_id	string	连击id
combo_timeout	int64	连击有效期秒
*/
type ComboInfo struct {
	ComboBaseNum	int64		`json:"combo_base_num"`
	ComboCount		int64		`json:"combo_count"`
	ComboId			string		`json:"combo_id"`
	ComboTimeout	int64		`json:"combo_timeout"`
}


/*
字段名	类型	描述
room_id	int64	直播间id
uid	int64	购买用户UID
uname	string	购买的用户昵称
uface	string	购买用户头像
message_id	int64	留言id(风控场景下撤回留言需要)
message	string	留言内容
rmb	int64	支付金额(元)
timestamp	int64	赠送时间秒级
start_time	int64	生效开始时间
end_time	int64	生效结束时间
guard_level	int64	对应房间大航海等级
fans_medal_level	int64	对应房间勋章信息
fans_medal_name	string	对应房间勋章名字
fans_medal_wearing_status	bool	该房间粉丝勋章佩戴情况
msg_id	string	消息唯一id
*/
type Socket_Liuyan struct {
	RoomId		int64		`json:"room_id"`
	Uid			int64		`json:"uid"`
	Uname		string		`json:"uname"`
	Uface		string		`json:"uface"`
	MessageId	int64		`json:"message_id"`
	Message		string		`json:"message"`
	Rmb			int64		`json:"rmb"`
	Timestamp	int64		`json:"timestamp"`
	StartTime	int64		`json:"start_time"`
	EndTime		int64		`json:"end_time"`
	GuardLevel	int64		`json:"guard_level"`
	FansMedalLevel	int64			`json:"fans_medal_level"`
	FansMedalName	string			`json:"fans_medal_name"`
	FansMedalWearingStatus	bool		`json:"fans_medal_wearing_status"`
	MsgId		string		`json:"msg_id"`
}

/*
字段名	类型	描述
room_id	int64	直播间id
message_ids	[]int64	留言id
msg_id	string	消息唯一id
*/
type Socket_Liuyan_Off struct {
	RoomId		int64		`json:"room_id"`
	MessageIds	[]int64		`json:"message_ids"`
	MsgId		string		`json:"msg_id"`
}

/*
字段名	类型	描述
user_info	user_info结构体	用户信息
guard_level	int64	大航海等级
guard_num	int64	大航海数量
guard_unit	string	大航海单位
fans_medal_level	int64	粉丝勋章等级
fans_medal_name	string	粉丝勋章名
fans_medal_wearing_status	bool	该房间粉丝勋章佩戴情况
room_id	int64	房间号
msg_id	string	消息唯一id
timestamp	int64	上舰时间秒级时间戳
*/
type Socket_Dahanghai struct {
	UserInfo	UserInfo	`json:"user_info"`
	GuardLevel	int64		`json:"guard_level"`
	GuardNum	int64		`json:"guard_num"`
	GuardUnit	string		`json:"guard_unit"`
	FansMedalLevel	int64			`json:"fans_medal_level"`
	FansMedalName	string		`json:"fans_medal_name"`
	FansMedalWearingStatus	bool	`json:"fans_medal_wearing_status"`
	RoomId		int64		`json:"room_id"`
	MsgId		string		`json:"msg_id"`
	Timestamp	int64		`json:"timestamp"`
}