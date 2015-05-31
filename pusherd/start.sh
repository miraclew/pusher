v=`git describe --long`
echo -e "package main\n\nconst BINARY_VERSION = \"$v\"" > v.go

#go build && go install && rm pusherd && pusherd -rethinkDb="mercury" -rethinkAddr="192.168.33.10:28015" -redisAddr="192.168.33.10:6379"
go build && go install && rm pusherd &&
pusherd -rethinkDb="sun" -rethinkAddr="192.168.33.10:28015" -redisAddr="192.168.33.10:6379"
