#! /bin/bash

echo "configuring zookeeper."
zkcfg \
-zooCfgPath=/opt/zookeeper/conf/ \
-zooCfgName=zoo.cfg \
-dynamicConfig=/opt/zookeeper/conf/zoo.cfg.dynamic \
-dataDir=/tmp/zookeeper/zk-$MYID \
-autopurge.snapRetainCount=10 \
autopurge.purgeInterval=24

if [ -n "${LAB+1}" ]; then
	echo "temporarily starting zookeeper in the background to stage data."
	zkServer.sh start
	sleep 10
	zkCli.sh -- -cmd create /$BUCKET_1
	zkCli.sh -- -cmd create /$BUCKET_2
	sleep 10
	zkServer.sh stop
	echo "stopped zookeeper and restarting in the foreground."
	zkServer.sh start-foreground
else
	echo "starting zookeeper."
	zkServer.sh start-foreground
fi
