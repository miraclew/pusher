mode=$1
if [ "$mode" = "-prod" ]; then
        echo "prod mode"
        scp ~/go/bin/linux_amd64/apns ubuntu@gx2:/data/pusher/apns_`date "+%m_%d_%H_%M_%S"`
else
        echo "dev mode"
        ssh ubuntu@gx1 -e "sudo supervisorctl stop apns"
        scp ~/go/bin/linux_amd64/apns ubuntu@gx1:/data/pusher/dev/apns
        ssh ubuntu@gx1 -e "sudo supervisorctl start apns"
fi

