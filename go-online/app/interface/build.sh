#!/bin/bash

servicePath=`pwd`

# Stop the running app
appdir=app/interface/cmd
app=interface
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
		nohup $servicePath/$appdir/$app -conf ../conf -log.dir /data/log_xunqunwang/interface/ >/data/log_xunqunwang/interface/interface.nohup &
    else
		nohup $servicePath/$appdir/$app -conf ../conf -log.dir ./../../../log/interface/ >./../../../log/interface/interface.nohup &
    fi
    echo "end..." 
fi
exit
