mode=$1
if [ "$mode" = "-prod" ]; then
		echo "prod mode"
        file_name="router_`date "+%m_%d_%H_%M_%S"`"
        scp ~/go/bin/linux_amd64/router ubuntu@gx2:/data/pusher/router/$file_name
        ssh ubuntu@gx2 sudo supervisorctl stop router:
        ssh ubuntu@gx2 rm /data/pusher/router/router
        ssh ubuntu@gx2 ln -s /data/pusher/router/$file_name /data/pusher/router/router
        ssh ubuntu@gx2 sudo supervisorctl start router:
else
		echo "dev mode"
        ssh ubuntu@gx1 sudo supervisorctl stop router:
        scp ~/go/bin/linux_amd64/router ubuntu@gx1:/data/pusher/dev/router
        ssh ubuntu@gx1 sudo supervisorctl start router:
fi

