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
    kill -9 $p
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
		nohup $servicePath/$appdir/$app -conf "user-release.toml" >/data/log_xunqunwang/user/user.nohup &
    else
		nohup $servicePath/$appdir/$app -conf "user-dev.toml" >./user.nohup &
    fi
    echo "end..." 
fi
exit
