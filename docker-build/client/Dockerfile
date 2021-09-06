ARG SOURCE_PREFIX=federatedai
ARG SOURCE_TAG=1.5.0-release
FROM ${SOURCE_PREFIX}/python:${SOURCE_TAG} as data

FROM python:3.7

COPY Examples /Examples
RUN pip install notebook fate-client pandas sklearn
RUN mkdir /Examples/Pipeline/logs
copy --from=data /data/projects/fate/examples /Examples/examples

CMD flow init -c /data/projects/fate/conf/service_conf.yaml && pipeline init -c /data/projects/fate/conf/pipeline_conf.yaml && jupyter notebook --ip=0.0.0.0 --port=20000 --allow-root --debug --NotebookApp.notebook_dir='/Examples' --no-browser --NotebookApp.token='' --NotebookApp.password=''
