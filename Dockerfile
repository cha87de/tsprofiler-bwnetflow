FROM ubuntu:latest
RUN apt-get update ; apt-get install -y ca-certificates
ADD tsprofiler-bwnetflow /bin/tsprofiler-bwnetflow
CMD [ "/bin/tsprofiler-bwnetflow" ]