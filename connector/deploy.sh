mode=$1
if [ "$mode" = "-prod" ]; then
        echo "prod mode"
        scp ~/go/bin/linux_amd64/connector ubuntu@gx2:/data/pusher/connector_`date "+%m_%d_%H_%M_%S"`
else
        echo "dev mode"
        ssh ubuntu@gx1 -e "sudo supervisorctl stop connector:"
        scp ~/go/bin/linux_amd64/connector ubuntu@gx1:/data/pusher/dev/connector
        ssh ubuntu@gx1 -e "sudo supervisorctl start connector:"
fi

