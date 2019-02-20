FROM alpine:latest
RUN apk update
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
ADD tsprofiler-bwnetflow /bin/tsprofiler-bwnetflow
CMD [ "/bin/tsprofiler-bwnetflow" ]