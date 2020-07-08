FROM python:3.7

WORKDIR /fml_manager
COPY . /fml_manager

RUN python setup.py sdist bdist_wheel && pip install dist/*.whl 
RUN pip install notebook

CMD jupyter notebook --ip=0.0.0.0 --port=20000 --allow-root --debug --NotebookApp.notebook_dir='/fml_manager/Examples' --no-browser --NotebookApp.token='' --NotebookApp.password=''
