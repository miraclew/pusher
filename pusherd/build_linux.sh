v=`git describe --long`

echo -e "package main\n\nconst BINARY_VERSION = \"$v\"" > v.go

GOOS=linux GOARCH=amd64 go build
GOOS=linux GOARCH=amd64 go install
rm pusherd
cp $GOPATH/bin/linux_amd64/pusherd ~/ubuntu/pusher
