mode=$1
if [ "$mode" = "-prod" ]; then
        echo "prod mode"
        file_name="connector_`date "+%m_%d_%H_%M_%S"`"
        scp ~/go/bin/linux_amd64/connector ubuntu@gx2:/data/pusher/connector/$file_name
        ssh ubuntu@gx2 sudo supervisorctl stop connector:
        ssh ubuntu@gx2 rm /data/pusher/connector/connector
        ssh ubuntu@gx2 ln -s /data/pusher/connector/$file_name /data/pusher/connector/connector
        ssh ubuntu@gx2 sudo supervisorctl start connector:
else
        echo "dev mode"
        ssh ubuntu@gx1 sudo supervisorctl stop connector:
        scp ~/go/bin/linux_amd64/connector ubuntu@gx1:/data/pusher/dev/connector
        ssh ubuntu@gx1 sudo supervisorctl start connector:
fi

