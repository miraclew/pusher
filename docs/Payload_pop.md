## 基本结构 ##

客户端收到的消息体是一个JSON字符串:

    {
        id: string
        type: int,
        sub_type: int,
        format: int,
        chat_id: int,
        sender_id: int,
        ttl: int, // 消息有效时间
        timestamp: int, // 发送时间(push服务收到该消息时间)(毫秒)
        title: string, // 标题
        body: string, // 内容
        extra: string // 额外数据
    }

Body是透传的，内容为API中的body

## Format 定义 ##
    FORMAT_TEXT = 1;
    FORMAT_IMAGE = 2;
    FORMAT_AUDIO = 3;
    FORMAT_VIDEO = 4;
    FORMAT_LINK = 5;
    FORMAT_GIFT = 6;

## 消息类型定义 ##
Type & Subtype:

    TYPE_CHAT = 1;  // 用户发送的消息
    ST_CHAT_MSG = 1001;  // 聊天

    TYPE_NOTIFY             = 2 // 对话事件
    ST_NOTIFY_NEW_FOLLOWER  = 2001 // 新粉丝
    ST_NOTIFY_POST_LIKE      = 2101 // 新点赞
    ST_NOTIFY_POST_COMMENT   = 2102  // 新评论

    TYPE_LIVE_EVENT             = 8
    ST_LIVE_CHAT                = 8001;// 评论
    ST_LIVE_SEND_GIFT           = 8002; // 打赏
    ST_LIVE_LIKE                = 8003;// 点赞
    ST_LIVE_NUM_UPDATE          = 8004;// 数字更新
    ST_LIVE_VIEWER_ENTER        = 8005;// 观众进入
    ST_LIVE_VIEWER_EXIT         = 8006;// 观众离开
    ST_LIVE_RECV_GIFT           = 8007;// 收到打赏
    ST_LIVE_END                 = 8008; // 直播结束
    ST_LIVE_PAUSE               = 8010; // 直播暂停
    ST_LIVE_RESUME              = 8011; // 直播恢复
    ST_LIVE_OPEN                = 8012; // 直播开播

## Body 定义 ##

### 聊天消息 TYPE_CHAT/ST_UM_CHAT ###

    body = "消息内容"


## TYPE_NOTIFY ##
### ST_NOTIFY_NEW_FOLLOWER ###
	body = "提示内容"

### ST_NOTIFY_POST_LIKE ###
	body = {
		'post_id': int
	}

### ST_NOTIFY_POST_COMMENT ###
	body = {
		'post_id': int,
		'comment_id': int
	}

## TYPE_LIVE_EVENT 直播消息推送 ##
### ST_LIVE_CHAT ###

    {
        "channel_id" : 124555,
        "content": "hello"
    }

### ST_LIVE_SEND_GIFT ###

    {
        "channel_id" : 124555,
        'gift_id': 1,
        'quantity': 10,
        'text': "aaa赠送了你礼物",
        'icon': "http://s.impop.me/gift1.png",
        'total_coins': 111,
    }

### ST_LIVE_LIKE ###

    {
        "channel_id" : 124555,
        "user_id" : 124555,
        "like_count": 5
    }

### ST_LIVE_NUM_UPDATE ###

    {
        "channel_id" : 124555,
        "online_viewers" : 123,
        "total_viewers" : 10000,
    }

### ST_LIVE_VIEWER_ENTER ###

    {
        "channel_id" : 124555,
        "user_id" : 124555,
        "nickname" : "abc",
        "avatar" : "abc",
        "online_viewers" : 123,
        "total_viewers" : 10000,
    }

### ST_LIVE_VIEWER_EXIT ###

    {
        "channel_id" : 124555,
        "user_id" : 124555,
        "online_viewers" : 123,
        "total_viewers" : 10000,
    }

### ST_LIVE_END ###

    {
        "channel_id" : 124555,
        "coins_spend" : 122,
        "coins_earn": 12,
    }

## Examples ##

文字聊天消息

    {
        type: 1,
        sub_type: 1001,
        chat_id: 123,
        ttl: 0,
        timestamp: 124555,
        body: "hello world!"
    }

图片聊天消息

    {
        id: "abcd",
        type: 1,
        sub_type: 1001,
        chat_id: 123,
        sender_id: 456,
        ttl: 0,
        timestamp: 124555,
        body: "http://mercury.uwang.me/avatar/1.jpg"
    }

