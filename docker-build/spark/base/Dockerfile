# Refer to https://github.com/big-data-europe/docker-spark/tree/2.4.1-hadoop2.7/base
ARG SOURCE_PREFIX=federatedai
ARG SOURCE_TAG=1.5.0-release
FROM ${SOURCE_PREFIX}/python:${SOURCE_TAG}

ENV ENABLE_INIT_DAEMON true
ENV INIT_DAEMON_BASE_URI http://identifier/init-daemon
ENV INIT_DAEMON_STEP spark_master_init

ENV SPARK_VERSION=2.4.1
ENV HADOOP_VERSION=2.7

WORKDIR /

RUN set -eux && \
    rpm --rebuilddb && \
    rpm --import /etc/pki/rpm-gpg/RPM* && \
    yum -y install gcc gcc-c++ make openssl-devel gmp-devel mpfr-devel libmpc-devel\
    libmpcdevel libaio numactl autoconf automake libtool libffi-devel  \
    snappy snappy-devel zlib zlib-devel bzip2 bzip2-devel lz4-devel libasan lsof sysstat telnet psmisc wget && \
    yum install -y which java-1.8.0-openjdk java-1.8.0-openjdk-devel && \
    yum clean all && \
    wget https://archive.apache.org/dist/spark/spark-${SPARK_VERSION}/spark-${SPARK_VERSION}-bin-hadoop${HADOOP_VERSION}.tgz && \
    tar -xvzf spark-${SPARK_VERSION}-bin-hadoop${HADOOP_VERSION}.tgz && \
    mv spark-${SPARK_VERSION}-bin-hadoop${HADOOP_VERSION} spark && \
    rm spark-${SPARK_VERSION}-bin-hadoop${HADOOP_VERSION}.tgz && \
    cd /

ENV PYTHONPATH=$PYTHONPATH:/data/projects/fate/python
ENV JAVA_HOME=/usr/lib/jvm/jre-1.8.0-openjdk
