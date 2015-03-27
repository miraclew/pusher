API 定义

## Channel ##

### POST /channels ###

Input:

    members: string // 逗号分割用户ID

Output:

    {
        code: int, 
        channel: {
            id: string
        }
    }
    
## Message ##
### POST /channel_msg ###

Input:

    type: int // 推送类型: 1. Simple
    sender_id: int // 发送者ID
    channel_id: string // channel ID
    payload: string, // 消息体
    push_offline: int, // 是否推送离线设备(默认:1)
    extra: string // 附加信息

Output:

    {
        code: int, 
        message: {
            id: string
        }
    }


### POST /private_msg ###

无需手动创建channel发送消息

Input:

    sender_id: int // 发送者ID
    receiver_id: int // 接收者ID
    payload: string, // 消息体
    extra: string // 附加信息

Output:

    {
        code: int, 
        message: {
            id: string
        }
    }