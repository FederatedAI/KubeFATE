FROM centos/python-36-centos7

WORKDIR /

COPY sql/ .
COPY upgrade-mysql.py .
COPY shutdown-flow.py .

RUN pip install --upgrade pip && \
    pip install mysql-connector-python && \
    pip install kubernetes