#!/bin/bash

servicePath=`pwd`

# Stop the running app
appdir=app/user/cmd
app=user
cd $appdir
p=$(pidof $servicePath/$appdir/$app )
if [ $p ]
then
    echo "kill $app pid $p"
    kill $p
fi

if [[ -n "$2" && "$2" = "start" ]]
then
    # Compile
    echo "$app building..."
    go build -o $app
    
    # Start
    echo "$app starting..."
    if [[ -n "$1" && "$1" = "release" ]]
    then
		nohup $servicePath/$appdir/$app -appid user -conf ../configs -log.dir /data/log_xunqunwang/user/ >/data/log_xunqunwang/user/user.nohup &
    else
		nohup $servicePath/$appdir/$app -appid user -conf ../conf -log.dir ./../../../log/user/ >./../../../log/user/user.nohup &
    fi
    echo "end..." 
fi
exit
