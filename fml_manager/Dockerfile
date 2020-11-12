FROM python:3.7

WORKDIR /fml_manager
COPY . /fml_manager

RUN pip install notebook fate-client
RUN python setup.py sdist bdist_wheel && pip install dist/*.whl 
RUN mkdir /fml_manager/Examples/Pipeline/logs

CMD flow init -c /data/projects/fate/conf/service_conf.yaml && pipeline init -c /data/projects/fate/conf/pipeline_conf.yaml && jupyter notebook --ip=0.0.0.0 --port=20000 --allow-root --debug --NotebookApp.notebook_dir='/fml_manager/Examples' --no-browser --NotebookApp.token='' --NotebookApp.password=''
