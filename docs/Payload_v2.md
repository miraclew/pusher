## 基本结构 ##

客户端收到的消息体是一个JSON字符串:

    {
        id: string
        type: int,
        sub_type: int,
        chat_id: int,
        sender_id: int,
        ttl: int, // 消息有效时间
        sent_at: int, // 发送时间(push服务收到该消息时间)
        sent_at_ms: int, // 发送时间(毫秒)
        body: {

        },
        extra: { // 额外数据
        }
    }

Body是透传的，内容为API中的body

## 消息类型定义 ##
Type & Subtype:

    TYPE_USER_MSG = 1;  // 用户发送的消息
    ST_UM_CHAT = 1001;  // 聊天
    ST_UM_TYPING = 1002; // 输入中

    TYPE_CHAT_EVENT = 2 // 对话事件
    ST_CE_CREATED = 2001
    ST_CE_MEMBER_JOIN = 2002
    ST_CE_MEMBER_LEAVE = 2003
    ST_CE_SYS_MSG = 2004 // 系统发送的文字消息
    ST_CE_CHAT_UPDATE = 2005 // Chat 信息更新

    TYPE_SYSTEM_EVENT = 3 // 系统事件
    ST_SE_ROLL_MSG = 3001 //滚屏消息


    TYPE_GENERAL_EVENT        = 4
    ST_GE_NEW_LIKE            = 4001 //点赞
    ST_GE_NEW_COMMENT         = 4002 //评论
    ST_GE_NEW_POST            = 4003 //帖子
    ST_GE_USER_NB_UPDATE      = 4100 // 用户数值更新

    TYPE_LOVER_EVENT          = 7
    ST_LE_NEW_LIKE            = 7001 //点赞
    ST_LE_NEW_COMMENT         = 7002 //评论
    ST_LE_NEW_POST            = 7003 //帖子
    ST_LE_ROOM_INFO           = 7005 //房间信息
    ST_LE_ROOM_INFO           = 7005 ; //房间信息
    ST_LE_GIFT                = 7006 ; //礼物
    ST_LE_HEART               = 7007 ; //真心话


    // 推送 (点击打开应用)
    const TYPE_NOTIFICATION_EVENT   = 5;
    const ST_NE_ALERT               = 5001; // 普通提醒，只需打开应用
    const ST_NE_NEW_LETTER          = 5002; // 新情书
    const ST_NE_LETTER_ARRIVE       = 5003; // 情书送达
    const ST_NE_TREE_UPGRADE        = 5004; // 爱情树升级
    const ST_NE_NEW_FANS            = 5005; // 新粉丝
    const ST_NE_NEW_WATERING        = 5006; // 新浇水

## Body 定义 ##

### 聊天消息 TYPE_USER_MSG/ST_UM_CHAT ###

    {
        mime: string
        content: {...}
		is_system: int
		is_bottle: int
		attach: int// 贴着 1: above, 2: below
    }

不同类型的mime对应content定义:

text

    {
        text: string // 聊天内容(包括emoji转义字符)
    }


image

    {
        url: string // 图片地址
    }


audio

    {
        url: string // 音频地址
    }

video

    {
        url: string // 视频地址
    }
    
link
 
	{
		type: int
		icon: string
		title: string
		desc: string
		url: string
	}
	
    const LINK_TYPE_URL             = 1; // 网址
    const LINK_TYPE_URL_BROWSER     = 2; // 网址(系统浏览器打开)
    const LINK_TYPE_POST            = 3; // 帖子详情
    const LINK_TYPE_LOVE_ROOM       = 4; // 恋人空间
    const LINK_TYPE_LOVE_ROOM_POST  = 5; // 恋人空间帖子详情
    const LINK_TYPE_LOVER           = 6; // 恋人专属

	
### 对话中的系统消息 ST_CE_SYS_MSG ###

    {
        text: string
    }

### 滚屏消息 ST_SE_ROLL_MSG ###

    {
        text: string
        times: int // 滚屏次数
    }


### ST_GE_USER_NB_UPDATE ###

    {
        followers: int,
        new_followers: int,
        visitors: int,
        new_visitors: int,
        new_letters: int,
        new_water: int,
        tree_level: int,
        new_lover_applies: int,
    }

## TYPE_LOVER_EVENT

该类型下所有子类型的消息格式相同

    {
        room_id: int,
        start_time: int,
        end_time: int,
        post: {
            id: int,
            type: int,
            title: string,
            text: string,
            images: string[],
            audio: {
                url: string,
                length: int,
            }
        },
        comment: string,
    }

## TYPE_NOTIFICATION_EVENT 推送 (点击打开应用) ##
### ST_GE_NOTIFICATION ###

body 格式 (Android)

    {
        text: string, // 提示文字
        rid: int, // 资源ID
    }


Apns 推送的消息格式 (iOS)

    {
        "aps" : {
              "alert” : {
              }
        },
        "notification": {
            "type": int,
            "sub_type": int,
            "rid": int,
            "sent_at": int
        }
    }

## Examples ##

文字聊天消息

    {
        type: 1,
        sub_type: 1001,
        chat_id: 123,
        ttl: 0,
        sent_at: 124555,
        body: {
            mime: "text",
            content: {
                text: "hello world!"
            }
        }
    }


图片聊天消息

    {
        id: "abcd",
        type: 1,
        sub_type: 1001,
        chat_id: 123,
        sender_id: 456,
        ttl: 0,
        sent_at: 124555,
        body: {
            mime: "image",
            content: {
                url: "http://mercury.uwang.me/avatar/1.jpg"
            }
        }
    }

滚屏消息

    {
        id: "abcd",
        type: 3,
        sub_type: 3001,
        chat_id: 0,
        sender_id: 456,
        ttl: 3600,
        sent_at: 124555,
        body: {
            text: "嘎嘎送礼10个钻石给呵呵",
            times: 3
        }
    }

送钻石消息

    {
        id: "abcd",
        type: 1,
        sub_type: 1005,
        chat_id: 1,
        sender_id: 456,
        ttl: 3600,
        sent_at: 124555,
        body: {
            text: 3
        }
    }

更新聊天状态消息

    {
        id: "abcd",
        type: 2,
        sub_type: 2005,
        chat_id: 1,
        sender_id: 456,
        ttl: 3600,
        sent_at: 124555,
        body: {
            id: 29,
            artist_id: 1,
            customer_id: 145,
            last_message: "",
            last_send_at: "0000-00-00 00:00:00",
            status: 2,
            channel_id: null,
        }
    }


