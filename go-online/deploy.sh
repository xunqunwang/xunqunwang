#!/bin/bash

#部署微服务脚本
#启动命令说明
# ./deploy.sh debug start 调试模式下启动所有微服务
# ./deploy.sh release start 生产部署模式下启动所有微服务
# ./deploy.sh debug stop 调试模式下关闭所有微服务
# ./deploy.sh release stop 生产部署模式下关闭所有微服务

prjPath=`pwd`
bash $prjPath/app/domain/identify/build.sh $1 $2
cd $prjPath
bash $prjPath/app/user/build.sh $1 $2
cd $prjPath
bash $prjPath/app/group/build.sh $1 $2
cd $prjPath
bash $prjPath/app/interface/build.sh $1 $2
exit