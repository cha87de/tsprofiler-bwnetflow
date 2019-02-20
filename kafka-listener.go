package main

import (
	"log"
	"runtime/debug"

	flow "omi-gitlab.e-technik.uni-ulm.de/bwnetflow/bwnetflow_api/go"
	ff "omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/kafkaconnector/flowfilter"
)

var flowFilter *ff.FlowFilter

func runKafkaListener() {
	// initialize filters: prepare filter arrays
	flowFilter = ff.NewFlowFilter(*filterCustomerIDs, *filterIPsv4, *filterIPsv6, *filterPeers)

	// handle kafka flow messages in foreground
	for {
		flow := <-kafkaConn.ConsumerChannel()
		if flowFilter.FilterApplies(flow) {
			handleFlow(flow)
		}
	}
}

func handleFlow(flow *flow.FlowMessage) {
	// handle panic while flow processing
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered panic in handleFlow", r)
			debug.PrintStack()
			log.Printf("failed flow: %+v\n", flow)
		}
	}()

	// the only action here: dump the flow
	// dumpFlow(flow)
	profileFlow(flow)
}
