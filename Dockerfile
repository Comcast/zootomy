FROM java:8-jre
MAINTAINER Mike Lloyd <mike_lloyd@cable.comcast.com>

ENV _JAVA_OPTIONS "-Xmx2G -Xincgc"

# get zookeeper 3.5.0-alpha
RUN curl -fLk http://apache.cs.utah.edu/zookeeper/zookeeper-3.5.0-alpha/zookeeper-3.5.0-alpha.tar.gz | tar xzf - -C /opt
RUN mv /opt/zookeeper-3.5.0-alpha /opt/zookeeper

# create the configs and data dir of /tmp/zookeeper
RUN mkdir /tmp/zookeeper/
RUN mv /opt/zookeeper/conf/zoo_sample.cfg /opt/zookeeper/conf/zoo.cfg
RUN touch /opt/zookeeper/conf/zoo.cfg.dynamic
RUN echo "quorumListenOnAllIPs=true" >> /opt/zookeeper/conf/zoo.cfg
RUN echo "standaloneEnabled=false" >> /opt/zookeeper/conf/zoo.cfg
RUN echo "dynamicConfigFile=/opt/zookeeper/conf/zoo.cfg.dynamic" >> /opt/zookeeper/conf/zoo.cfg
RUN echo "/zookeeper/bin/zkServer.sh status | egrep 'Mode: (standalone|leading|following|observing)'" > /opt/zookeeper/bin/zkReady.sh

# zkcfg reads the environmental variables and then writes them
# to the configured files. the reason this is not a dockerfile ENTRYPOINT
# command is due to the instruction order of dockerfile commands. ENTRYPOINT
# is run as a prepend to the next CMD function.
# for more information: https://goo.gl/p7tzlz
COPY zkcfg /usr/local/bin
COPY run.sh /

ENV PATH=/opt/zookeeper/bin:${PATH} \
    ZOO_LOG4J_PROP="INFO, CONSOLE"

EXPOSE 2181 2888 3888

CMD bash run.sh
