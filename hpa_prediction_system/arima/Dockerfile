FROM ubuntu:18.04

WORKDIR /usr/bin

RUN apt-get update -y && \
    apt-get install -y python3.7 python3-pip python3.7-dev
RUN pip3 install --upgrade pip
RUN pip3 install --no-cache-dir numpy matplotlib pandas sklearn statsmodels

COPY . .

EXPOSE 50000

ENTRYPOINT ['./arima.py']
