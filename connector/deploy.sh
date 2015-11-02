host=$1
app=connector
if [[ "$host" == "" ]]; then
        echo "deploy host is required"
else
        file_name=$app"_`date "+%m_%d_%H_%M_%S"`"
        scp ~/go/bin/linux_amd64/$app ubuntu@$host:/data/pusher/bin/$file_name
        echo "stop $app"
        ssh ubuntu@$host sudo supervisorctl stop $app:
        ssh ubuntu@$host rm /data/pusher/bin/$app
        ssh ubuntu@$host ln -s /data/pusher/bin/$file_name /data/pusher/bin/$app
        echo "start $app"
        ssh ubuntu@$host sudo supervisorctl start $app:
fi