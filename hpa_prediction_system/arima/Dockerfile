FROM ubuntu:18.04

RUN apt-get update -y && \
    apt-get install -y python3.7 python3-pip python3.7-dev
RUN pip3 install --upgrade pip
RUN pip3 install --no-cache-dir numpy matplotlib pandas sklearn statsmodels

RUN mkdir /arima

COPY . /arima

EXPOSE 50000

ENTRYPOINT ['/arima/arima.py']
