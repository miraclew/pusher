# Client message payload

## 客户端发送消息格式 ##

客户端发送的消息体是一个JSON字符串:

    {
        type: int,
        ack_msg_id: string // optional
    }


## Type Defines

目前只有一种格式的消息

TYPE_ACK_MSG = 6001

