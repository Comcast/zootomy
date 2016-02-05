FROM java:8-jre
MAINTAINER Mike Lloyd <mike_lloyd@cable.comcast.com>

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

ENV PATH=/opt/zookeeper/bin:${PATH} \
    ZOO_LOG4J_PROP="INFO, CONSOLE"

EXPOSE 2181 2888 3888

CMD run.sh
