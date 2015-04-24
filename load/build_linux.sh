GOOS=linux GOARCH=amd64 go build
GOOS=linux GOARCH=amd64 go install
rm load
cp $GOPATH/bin/linux_amd64/load ~/ubuntu/pusher/
