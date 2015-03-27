IM Server (Golang)
=======
#About
Instant messaging service like Weixin/QQ

#Features
1. Realtime message delivery
2. Offline message store and resend
3. History message store
4. Support both Private and Group messages

#Architecture

## Components
1. HTTP API for message sending
2. Websocket connection to push realtime message

## Data storage
1. Redis (auth token, user message queue, apn device token)
2. RethinkDb (channels, messages)

# Cross build
For linux  

GOOS=linux GOARCH=amd64 go build  
GOOS=linux GOARCH=amd64 go install  
$GOPATH/bin/linux_amd64

# Testing

amazing testing with: Goconvey