#!/bin/bash

servicePath=`pwd`

# Stop the running app
appdir=app/group/cmd
app=group
cd $appdir
p=$(pidof $servicePath/$appdir/$app )
if [ $p ]
then
    echo "kill $app pid $p"
    kill  $p
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
		nohup $servicePath/$appdir/$app -appid group -conf ../configs -log.dir /data/log_xunqunwang/group/ >/data/log_xunqunwang/group/group.nohup &
    else
		nohup $servicePath/$appdir/$app -appid group -conf ../conf -log.dir ./../../../log/group/ >./../../../log/group/group.nohup &
    fi
    echo "end..." 
fi
exit
