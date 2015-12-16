#! /bin/bash

echo "configuring zookeeper."
zkcfg -dataDir=/tmp/zookeeper/ -filepath=/opt/zookeeper/conf/ -filename=zoo.cfg.dynamic

echo "starting zookeeper."
zkServer.sh start-foreground
