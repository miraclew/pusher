v=`git describe --long`

echo -e "package main
import (
    \"fmt\"
    \"runtime\"
)

const BINARY_VERSION = \"$v\"

func Version(app string) string {
    return fmt.Sprintf(\"%s %s (built w/%s)\", app, BINARY_VERSION, runtime.Version())
}
" > version.go

GOOS=linux GOARCH=amd64 go build -o pconnector
GOOS=linux GOARCH=amd64 go install -o pconnector
cp $GOPATH/bin/linux_amd64/pconnector ~/ubuntu/pconnector
rm pconnector
