ARG SOURCE_PREFIX=federatedai
ARG SOURCE_TAG=1.5.0-release
FROM ${SOURCE_PREFIX}/python:${SOURCE_TAG}

RUN rpm --rebuilddb && \
    rpm --import /etc/pki/rpm-gpg/RPM* && \
    yum install -y  which java-1.8.0-openjdk wget && \
    yum clean all && \
    wget https://archive.apache.org/dist/hadoop/common/hadoop-2.7.4/hadoop-2.7.4.tar.gz && \
    tar -xf ./hadoop-2.7.4.tar.gz -C /data/projects/ && rm ./hadoop-2.7.4.tar.gz

RUN wget https://archive.apache.org/dist/spark/spark-2.4.1/spark-2.4.1-bin-hadoop2.7.tgz && \
    tar -xf ./spark-2.4.1-bin-hadoop2.7.tgz -C /data/projects/ && rm ./spark-2.4.1-bin-hadoop2.7.tgz

ENV JAVA_HOME=/usr/lib/jvm/jre-1.8.0-openjdk
ENV SPARK_HOME=/data/projects/spark-2.4.1-bin-hadoop2.7/
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/data/projects/hadoop-2.7.4/lib/native
ENV PATH=$PATH:/data/projects/spark-2.4.1-bin-hadoop2.7/bin:/data/projects/hadoop-2.7.4/bin

RUN pip install pyspark