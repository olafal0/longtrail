# Name: registry.gitlab.com/olafalo/longtrail/builder
# docker build -t registry.gitlab.com/olafalo/longtrail/builder .
# docker push registry.gitlab.com/olafalo/longtrail/builder

FROM golang:latest

RUN apt-get update
RUN apt-get install -y zip jq python-pip
RUN pip install awscli --upgrade
