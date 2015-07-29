mode=$1
if [ "$mode" = "-prod" ]; then
        echo "prod mode"
        file_name="apns_`date "+%m_%d_%H_%M_%S"`"
        scp ~/go/bin/linux_amd64/apns ubuntu@gx2:/data/pusher/apns/$file_name
        ssh ubuntu@gx2 sudo supervisorctl stop apns
        ssh ubuntu@gx2 rm /data/pusher/apns/apns
        ssh ubuntu@gx2 ln -s /data/pusher/apns/$file_name /data/pusher/apns/apns
        ssh ubuntu@gx2 sudo supervisorctl start apns
else
        echo "dev mode"
        ssh ubuntu@gx1 sudo supervisorctl stop apns
        scp ~/go/bin/linux_amd64/apns ubuntu@gx1:/data/pusher/dev/apns
        ssh ubuntu@gx1 sudo supervisorctl start apns
fi

