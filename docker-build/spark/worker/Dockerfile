ARG SOURCE_PREFIX=federatedai
ARG SOURCE_TAG=1.5.0-release
FROM ${SOURCE_PREFIX}/spark-base:${SOURCE_TAG}

COPY worker.sh /

ENV SPARK_WORKER_WEBUI_PORT 8081
ENV SPARK_WORKER_LOG /spark/logs
ENV SPARK_MASTER "spark://spark-master:7077"

EXPOSE 8081

CMD ["/bin/bash", "/worker.sh"]