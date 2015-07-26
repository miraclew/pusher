mode=$1
if [ "$mode" = "-prod" ]; then
		echo "prod mode"
        scp ~/go/bin/linux_amd64/router ubuntu@gx2:/data/pusher/router_`date "+%m_%d_%H_%M_%S"`
else
		echo "dev mode"
        ssh ubuntu@gx1 sudo supervisorctl stop router:
        scp ~/go/bin/linux_amd64/router ubuntu@gx1:/data/pusher/dev/router
        ssh ubuntu@gx1 sudo supervisorctl start router:
fi

