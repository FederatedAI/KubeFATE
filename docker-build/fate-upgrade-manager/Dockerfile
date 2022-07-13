FROM python:3.9.2-alpine3.13

WORKDIR /

RUN pip install mysql-connector-python

COPY upgrade-mysql.py .
COPY sql ./sql
