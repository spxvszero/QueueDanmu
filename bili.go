package main

import (
	"fmt"
	"log"
	"queue_danmu/bili"
	"strconv"
)

const (
	//测试IM时候用的直播码，通过界面输入话就不需要了
	Code = "XXXXXXXXXXXXX"
	//bilibili 饭贩上创建项目的ID
	AppID = 1888888888888
)

func RunBiliDanmu(IdentifyCode string) error {

	result,err := bili.SendAppStartReq(IdentifyCode,AppID)
	if err != nil {
		return err
	}
	fmt.Println("[BiliDanmu] Start : ",result)

	var c *bili.BiliWsClient
	for _, v := range result.WebsocketInfo.WssLink {
		c = bili.NewBiliWsClient(&bili.BiliWsClientConfig{
			Path: v,
			AuthBody: result.WebsocketInfo.AuthBody,
		})
		if c != nil {
			fmt.Println("[BiliDanmu] Link ws : ", v, " Success.")
			c.OnReciveDanMuMsg = RecieveDanmu
			c.OnReciveGiftMsg = RecieveLiwu
			break;
		}
		fmt.Println("[BiliDanmu] Link ws : ", v, " Failed. Try Next..")
	}
	if c == nil {
		return fmt.Errorf("[BiliDanmu]Client Init Error")
	}
	go func() {
		c.Run()
	}()

	return nil
}

func RecieveDanmu(danmu bili.Socket_Danmu)  {
	//check if can queue
	if QueueFilterData.QueueKind == QueueTypeDanmu && danmu.Msg == QueueFilterData.QueueContent {
		if QueueFilterData.UserEnable > 0 {
			if danmu.GuardLevel > 0 {
				//1总督 2提督 3舰长
				danmuKind := int64(1) << (danmu.GuardLevel + 1)
				if (int64(QueueFilterData.UserEnable) & danmuKind) > 0{
					levelName := "未知"
					if danmu.GuardLevel == 1 {
						levelName = "总督"
					}else if danmu.GuardLevel == 2 {
						levelName = "提督"
					}else if danmu.GuardLevel == 3{
						levelName = "舰长"
					}
					//enable
					addToQueue(QueueItemData{fmt.Sprintf("[%s]",levelName), danmu.Uname})
				}
			}
			if QueueFilterData.UserEnable & UserTypeFans > 0 && danmu.FansMedalLevel >= QueueFilterData.UserEnableLevel {
				//粉丝满足
				//addToQueue(fmt.Sprintf("[%sLv%d]❁%s",danmu.FansMedalName,danmu.FansMedalLevel, danmu.Uname))
				addToQueue(QueueItemData{fmt.Sprintf("[%sLv%d]",danmu.FansMedalName,danmu.FansMedalLevel), danmu.Uname})
				return
			}
		}else {
			if danmu.GuardLevel > 0 {
				//1总督 2提督 3舰长
				levelName := "未知"
				if danmu.GuardLevel == 1 {
					levelName = "总督"
				}else if danmu.GuardLevel == 2 {
					levelName = "提督"
				}else if danmu.GuardLevel == 3{
					levelName = "舰长"
				}
				//enable
				addToQueue(QueueItemData{fmt.Sprintf("[%s]",levelName), danmu.Uname})
			}else if (danmu.FansMedalLevel >= 1){
				addToQueue(QueueItemData{fmt.Sprintf("[%sLv%d]",danmu.FansMedalName,danmu.FansMedalLevel), danmu.Uname})
			}else {
				addToQueue(QueueItemData{"[路人]",danmu.Uname})
			}

		}
	}
}
func RecieveLiwu(liwu bili.Socket_Liwu)  {
	if QueueFilterData.QueueKind == QueueTypeGift {
		giftTargetVal, err := strconv.Atoi(QueueFilterData.QueueContent)
		if err != nil {
			log.Println("Gift Filter is not correct")
			return
		}
		if liwu.Price >= int64(giftTargetVal) {
			//addToQueue(liwu.Uname)
		}
	}
}
