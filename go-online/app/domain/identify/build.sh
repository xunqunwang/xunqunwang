#!/bin/bash

servicePath=`pwd`

# Stop the running app
appdir=app/domain/identify/cmd
app=domain-identify
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
		nohup $servicePath/$appdir/$app -appid domain.identify -conf ../configs -log.dir /data/log_xunqunwang/domain/identify/ >/data/log_xunqunwang/domain/identify/identify.nohup &
    else
		nohup $servicePath/$appdir/$app -appid domain.identify -conf ../conf -log.dir ./../../../../log/domain/identify/ >./../../../../log/domain/identify/identify.nohup &
    fi
    echo "end..." 
fi
exit
