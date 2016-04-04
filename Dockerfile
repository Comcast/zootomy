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

FROM java:8-jre
MAINTAINER Mike Lloyd <mike_lloyd@kevin.michael.lloyd@gmail.com>

# Maintainer's Note:
# While there are multiple layers, which technically sits outside the Docker officia recommendations, more layers equals
# more visility since this is essentially a build process. This allows for easier extensibility for other developers and
# users to extend this to fit their needs.

# enable incremental garbage collection and set the heap to max at 2GB.
ENV _JAVA_OPTIONS "-Xmx2G -Xincgc"

# get zookeeper 3.5.0-alpha
RUN curl -fLk http://apache.cs.utah.edu/zookeeper/zookeeper-3.5.0-alpha/zookeeper-3.5.0-alpha.tar.gz | tar xzf - -C /opt
RUN mv /opt/zookeeper-3.5.0-alpha /opt/zookeeper

# create the configs and default data dir of /tmp/zookeeper
RUN mkdir /tmp/zookeeper/

# zkcfg is the configuration program for Zookeeper. the reason
# this is not a dockerfile ENTRYPOINT command is due to the
# instruction order of dockerfile commands. ENTRYPOINT is run
# as a prepend to the next CMD function.
# for more information: https://goo.gl/p7tzlz
COPY zkcfg /usr/local/bin
COPY run.sh /opt/zookeeper/bin
COPY prestage.sh /opt/zookeeper/bin

ENV PATH=/opt/zookeeper/bin:${PATH} \
    ZOO_LOG4J_PROP="INFO, CONSOLE"

EXPOSE 2181 2888 3888 8080

CMD run.sh
