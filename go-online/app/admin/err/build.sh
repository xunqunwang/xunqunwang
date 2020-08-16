#!/bin/bash

servicePath=`pwd`

# Stop the running app
appdir=app/admin/err/cmd
app=admin-err
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
		nohup $servicePath/$appdir/$app -conf ../conf -log.dir /data/log_xunqunwang/admin/err/  >/data/log_xunqunwang/admin/err/err.nohup &
    else
		nohup $servicePath/$appdir/$app -conf ../conf ./../../../../log/admin/err/  >./../../../../log/admin/err/err.nohup &
    fi
    echo "end..." 
fi
exit
