FROM python:3.7

COPY .Examples /Examples
RUN pip install notebook fate-client
RUN python setup.py sdist bdist_wheel && pip install dist/*.whl 
RUN mkdir /Examples/Pipeline/logs

CMD flow init -c /data/projects/fate/conf/service_conf.yaml && pipeline init -c /data/projects/fate/conf/pipeline_conf.yaml && jupyter notebook --ip=0.0.0.0 --port=20000 --allow-root --debug --NotebookApp.notebook_dir='/Examples' --no-browser --NotebookApp.token='' --NotebookApp.password=''
