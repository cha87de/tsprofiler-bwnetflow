# tsprofiler-bwnetflow
Integration of tsprofiler into bwnetflow

```
docker pull cha87de/tsprofiler-bwnetflow:master
docker run -d \
    --name tsprofiler-bwnetflow \
    cha87de/tsprofiler-bwnetflow:master \
    /bin/tsprofiler-bwnetflow \
        --kafka.brokers BROKERS \
        --kafka.consumer_group GROUP \
        --kafka.topic flow-messages-enriched \
        --kafka.user USER \
        --kafka.pass PASS \
        --filter.customerid 10109
```