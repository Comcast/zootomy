#! /bin/bash

echo "configuring zookeeper."
zkcfg \
-zooCfgPath=/opt/zookeeper/conf/ \
-zooCfgName=zoo.cfg \
-dynamicConfig=/opt/zookeeper/conf/zoo.cfg.dynamic \
-dataDir=/tmp/zookeeper/zk-$MYID

echo "starting zookeeper."
zkServer.sh start-foreground
