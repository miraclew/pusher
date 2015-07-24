mode=$1
if [ "$mode" = "-prod" ]; then
        echo "prod mode"
        scp ~/go/bin/linux_amd64/connector ubuntu@gx2:/data/pusher/connector_`date "+%m_%d_%H_%M_%S"`
else
        echo "dev mode"
        scp ~/go/bin/linux_amd64/connector ubuntu@gx1:/data/pusher/dev/connector
fi

