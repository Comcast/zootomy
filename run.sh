# Copyright 2016 Comcast Cable Communications Management, LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#! /bin/bash

echo "Configuring zookeeper."
zkcfg \
-zooCfgPath=/opt/zookeeper/conf/ \
-zooCfgName=zoo.cfg \
-dynamicConfig=/opt/zookeeper/conf/zoo.cfg.dynamic \
-dataDir=/tmp/zookeeper/zk-$MYID \
-autopurge.snapRetainCount=10 \
-autopurge.purgeInterval=24

if [ -n "${LAB+1}" ]; then
	parallel ::: prestage.sh "zkServer.sh start-foreground"
else
	echo "Starting zookeeper."
	zkServer.sh start-foreground
fi
