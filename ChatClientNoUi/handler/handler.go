package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
	"zChatRoom/proto"

	"github.com/pzqf/zUtil/zTime"

	"github.com/pzqf/zEngine/zNet"
)

var contentList = []string{
	"1.I’m an office worker. 我是上班族",
	"2.I work for the government. 我在政府机关做事",
	"3.I’m happy to meet you. 很高兴见到你",
	"4.I like your sense of hum123 我喜欢你的幽默感",
	"5.I’m glad to see you again. 很高兴再次见到你",
	"6.I’ll call you. 我会打电话给你",
	"7.I feel like sleeping/ taking a walk. 我想睡/散步",
	"8.I want something to eat. 我想吃点东西",
	"9.I need your help. 我需要你的帮助",
	"10.I would like to talk to you for a minute. 我想和你谈一下",
	"11.I have a lot of problems. 我有很多问题",
	"12.I hope our dreams come true. 我希望我们的梦想成真",
	"13.I’m looking forward to seeing you. 我期望见到你",
	"14.I’m supposed to go on a diet / get a raise. 我应该节食/涨工资",
	"15.I heard that you’re getting married. Congratulations.听说你要结婚了，恭喜!",
	"16.I see what your mean. 我了解你的意思",
	"17.I can’t do this. 我不能这么做",
	"18.Let me explain why I was late. 让我解释迟到的理由",
	"19.Let’s have a beer or something. 咱们喝点啤酒什么的",
	"20.Where is your office? 你们的办公室在哪?",
	"21.What is your plan? 你的计划是什么?",
	"22.When is the store closing? 这家店什么时候结束营业?",
	"23.Are you sure you can come by at nine? 你肯定你九点能来吗?",
	"24.Am I allowed to stay out past 10? 我可以十点过后再回家吗?",
	"25.The meeting was scheduled for two hours, but it is now over yet. 会议原定了两个小时，不过现在还没有结束",
	"26.Tom’s birthday is this week. 汤姆的生日就在这个星期",
	"27.Would you care to see it/ sit down for a while? 你要不要看/坐一会呢?",
	"28.Can you cover for me on Friday/help me/ tell me how to get there? 星期五能不能请你替我个班/你能帮我吗/你能告诉我到那里怎么走吗?",
	"29.Could you do me a big favor? 能否请你帮我个忙?",
	"30.He is crazy about Crazy English. 他对疯狂英语很着迷",
	"31.Can you imagine how much he paid for that car?你能想象他买那车花了多少钱吗?",
	"32.Can you believe that I bought a TV for $25?你能相信我用25美元买了一台电视机吗?",
	"33.Did you know he was having an affair/cheating on his wife? 你知道他有外遇了吗?/欺骗他的妻子吗?",
	"34.Did you hear about the new project? 你知道那个新项目吗?",
	"35.Do you realize that all of these shirts are half off? 你知道这些衬衫都卖半价了吗?",
	"36.Are you mind if I take tomorrow off? 你介意我明天请假吗?",
	"37.I enjoy working with you very much. 我很喜欢和你一起工作",
	"38.Did you know that Stone ended up marrying his secretary? 你知道吗?斯通最终和他的秘书结婚了",
	"39.Let’s get together for lunch. 让我们一起吃顿午餐吧",
	"40.How did you do on your test?　你这次考试的结果如何?",
	"41.Do you think you can come? 你认为你能来吗?",
	"42.How was your weekend ? 你周末过得怎么样?",
	"43.Here is my card. 这是我的名片",
	"44.He is used to eating out all the time. 他已经习惯在外面吃饭了",
	"45.I’m getting a new computer for birthday present. 我得到一台电脑作生日礼物",
	"46.Have you ever driven a BMW? 你有没有开过“宝马”?",
	"47.How about if we go tomorrow instead? 我们改成明天去怎么样?",
	"48.How do you like Hong Kong? 你喜欢香港吗?",
	"49.How do you want your steak? 你的牛排要几分熟?",
	"50.How did the game turn out? 球赛结果如何?",
	"51.How did Mary make all of her money? 玛丽所有的钱是怎么赚到的?",
	"52.How was your date? 你的约会怎么样?",
	"53.How are you doing with your new boss? 你跟你的新上司处得如何?",
	"54.How should I tell him the bad news? 我该如何告诉他这个坏消息?",
	"55.How much money did you make? 你赚了多少钱?",
	"56.How much does it cost to go abroad? 出国要多少钱?",
	"57.How long will it take to get to your house? 到你家要多久?",
	"58.How long have you been here? 你在这里多久了?",
	"59.Hownice/pretty/cold/funny/stupid/boring/interesting.",
	"60.How about going out for dinner? 出去吃晚餐如何?",
	"61.I’m sorry that you didn’t get the job. 很遗憾，你没有得到那份工作",
	"62.I’m afraid that it’s not going to work out. 我恐怕这事不会成的",
	"63.I guess I could come over. 我想我能来",
	"64.it okay to smoke in the office? 在办公室里抽烟可以吗?",
	"65.It was kind of exciting. 有点剌激",
	"66.I know what you want. 我知道你想要什么",
	"67.that why you don’t want to go home? 这就是你不想回家的原因吗?",
	"68.I’m sure we can get you a great / good deal. 我很肯定我们可以帮你做成一笔好交易",
	"69.Would you help me with the report? 你愿意帮我写报告吗?",
	"70.I didn’t know he was the richest person in the world.我不知道他是世界上最有钱的人",
	"71.I’ll have to ask my boss/wife first.我必须先问一下我的老板/老婆",
	"72.I take it you don’t agree. 这么说来，我认为你是不同意",
	"73.I tried losing weight, but nothing worked. 我曾试着减肥，但是毫无效果",
	"74.It doesn’t make any sense to get up so early.那么早起来没有任何意义",
	"75.It took years of hard work to speak good English. 讲一口流利的英语需要多年的刻苦操练",
	"76.It feels like spring/ I’ve been here before. 感觉好象春天到了/我以前来过这里",
	"77.I wonder if they can make it. 我在想他们是不是能办得到",
	"78.It’s not as cold / hot as it was yesterday. 今天不想昨天那么冷/热",
	"79.It’s not his work that bothers me; it’s his attitude. 困扰我的不是他的工作，而是他的态度",
	"80.It sounds like you enjoyed it. 听起来你好象蛮喜欢的",
	"81.It seems to me that be would like to go back home. 我觉得他好象想要回家",
	"82.It looks very nice. 看起来很漂亮",
	"83.everything under control? 一切都在掌握之中吗?",
	"84.I thought you could do a better job. 我以为你的表现会更好",
	"85.It’s time for us to say “No” to America. 是我们对美国说不的时候了",
	"86.The show is supposed to be good. 这场表演应当是相当好的",
	"87.It really depends on who is in charge. 那纯粹要看谁负责了",
	"88.It involves a lot of hard work. 那需要很多的辛勤工作",
	"89.That might be in your favor. 那可能对你有利",
	"90.I didn’t realize how much this meant to you. 我不知道这个对你的意义有这么大",
	"91.I didn’t mean to offend you. 我不是故意冒犯你",
	"92.I was wondering if you were doing anything this weekend. 我想知道这个周末你有什么要做",
	"93.May I have your attention., please? 请大家注意一下",
	"94.This is great golfing / swimming/ picnic weather. 这是个打高尔夫球/游泳/野餐的好天气",
	"95.Thanks for taking me the movie. 谢谢你带我去看电影",
	"96.I am too tired to speak. 我累得说不出活来",
	"97.Would you tell me your phone number? 你能告诉我你的电话号码吗?",
	"98.Where did you learn to speak English? 你从哪里学会说英语的呢?",
	"99.There is a TV show about AIDS on right now. 电视正在播放一个关于爱滋病的节目",
	"100.What do you think of his new job/ this magazine? 你对他的新工作/这本杂志看法如何?",
}

func Init() {
	if err := zNet.RegisterHandler(proto.PlayerLogin, PlayerLoginRes); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
	if err := zNet.RegisterHandler(proto.PlayerLogout, PlayerLogoutRes); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
	if err := zNet.RegisterHandler(proto.PlayerEnterRoom, PlayerEnterRoomRes); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
	if err := zNet.RegisterHandler(proto.PlayerLeaveRoom, PlayerLeaveRoomRes); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
	if err := zNet.RegisterHandler(proto.PlayerSpeak, PlayerSpeakRes); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
	if err := zNet.RegisterHandler(proto.SpeakBroadcast, SpeakBroadcast); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
	if err := zNet.RegisterHandler(proto.RoomList, RoomListRes); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
	if err := zNet.RegisterHandler(proto.RoomPlayerList, RoomPlayerListRes); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
}

func PlayerLoginRes(session zNet.Session, protoId int32, d []byte) {
	var data proto.PlayerLoginRes
	//err := packet.DecodeData(&data)
	err := json.Unmarshal(d, &data)
	if err != nil {
		return
	}

	if data.Code != proto.ErrNil {
		return
	}
	fmt.Println("登录成功")
}

func PlayerLogoutRes(session zNet.Session, protoId int32, d []byte) {}

func PlayerEnterRoomRes(session zNet.Session, protoId int32, d []byte) {
	var data proto.PlayerEnterRoomRes
	//err := packet.DecodeData(&data)
	err := json.Unmarshal(d, &data)
	if err != nil {
		return
	}

	if data.Code != proto.ErrNil {
		return
	}

	//开始自动聊天
	for {
		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(len(contentList) - 1)
		req, _ := json.Marshal(proto.PlayerSpeakReq{
			Content: contentList[index],
		})
		err := session.Send(proto.PlayerSpeak, req)
		if err != nil {
			fmt.Println(err)
			return
		}

		time.Sleep(time.Second * 30)
	}
}

func PlayerLeaveRoomRes(session zNet.Session, protoId int32, d []byte) {

	log.Println("离开房间成功")
}

func PlayerSpeakRes(session zNet.Session, protoId int32, d []byte) {
	var data proto.PlayerSpeakRes
	//_ = packet.DecodeData(&data)
	_ = json.Unmarshal(d, &data)
	if data.Code != proto.ErrNil {
		return
	}
}

func SpeakBroadcast(session zNet.Session, protoId int32, d []byte) {
	var data proto.ChatMessage
	err := json.Unmarshal(d, &data)
	if err != nil {
		return
	}

	//fmt.Println(formatSpeakContent(data))
}

func RoomListRes(session zNet.Session, protoId int32, d []byte) {
	var data proto.RoomListRes
	//err := packet.DecodeData(&data)
	err := json.Unmarshal(d, &data)
	if err != nil {
		return
	}

	if data.Code != proto.ErrNil {
		//cui.ShowDialog(data.Message, cui.DialogTypeError)
		return
	}

	var list []string
	list = append(list, "\t\tid\t\t\tname")
	for _, v := range data.RoomList {
		list = append(list, fmt.Sprintf("\t\t%d\t\t\t%s", v.Id, v.Name))
	}
	//cui.ShowRoomUi(list)

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(data.RoomList) - 1)

	req, _ := json.Marshal(proto.PlayerEnterRoomReq{
		RoomId: data.RoomList[index].Id,
	})

	err = session.Send(proto.PlayerEnterRoom, req)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func RoomPlayerListRes(session zNet.Session, protoId int32, d []byte) {
	var data proto.RoomPlayerListRes
	//err := packet.DecodeData(&data)
	err := json.Unmarshal(d, &data)
	if err != nil {
		return
	}

	if data.Code != proto.ErrNil {
		//cui.ShowDialog(data.Message, cui.DialogTypeError)
		return
	}

	var list []string
	for _, v := range data.RoomPlayerList {
		list = append(list, v.Name)
	}
}

func formatSpeakContent(data proto.ChatMessage) string {
	str := ""
	if data.Name != "" {
		str = fmt.Sprintf("%s [ %s ] say: %s", zTime.Seconds2String(data.Time), data.Name, data.Content)
	} else {
		str = fmt.Sprintf("%s   %s", zTime.Seconds2String(data.Time), data.Content)
	}

	return str
}
