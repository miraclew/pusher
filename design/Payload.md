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
    ST_UM_FIREND_ADD = 1003; // 好友申请
    ST_UM_FIREND_CONFIRM = 1004; // 好友确认
    ST_UM_DONATE_COIN = 1005; // 赠送钻石
    ST_UM_RATING = 1006; // 评价
    ST_UM_PAY = 1007; // 付费后给艺人发消息

    TYPE_CHAT_EVENT = 2 // 对话事件
    ST_CE_CREATED = 2001
    ST_CE_MEMBER_JOIN = 2002
    ST_CE_MEMBER_LEAVE = 2003
    ST_CE_SYS_MSG = 2004 // 系统发送的文字消息
    ST_CE_CHAT_UPDATE = 2005 // Chat 信息更新

    TYPE_SYSTEM_EVENT = 3 // 系统事件
    ST_SE_ROLL_MSG = 3001 //滚屏消息 


## Body 定义 ##

### 聊天消息 TYPE_USER_MSG/ST_UM_CHAT ###

    {
        mime: string
        content: string
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

### 对话中的系统消息 ST_CE_SYS_MSG ###

    {
        text: string
    }

### 滚屏消息 ST_SE_ROLL_MSG ###

    {
        text: string
        times: int // 滚屏次数
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
