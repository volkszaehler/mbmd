package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

// socketCollectors contains all Prometheus metrics about web sockets
//
// Implements collectable interface
type socketCollectors struct{}

var (
	WebSocketClientConnectionClose = prometheus.NewCounter(
		newCounterOpts(
			"websocket_client_connections_closed_total",
			"Total number of closed client connections to a web socket",
		),
	)

	WebSocketClientMessageSendSuccess = prometheus.NewCounter(
		newCounterOpts(
			"websocket_client_message_send_successes_total",
			"Total number of messages sent to a web socket client",
		),
	)

	WebSocketClientMessageSendFailure = prometheus.NewCounter(
		newCounterOpts(
			"websocket_client_message_send_failures_total",
			"Total number of message send failures to a web socket client",
		),
	)

	WebSocketClientCreationFailure = prometheus.NewCounterVec(
		newCounterOpts(
			"websocket_client_creation_failures_total",
			"Total number of accepting and failed creation of a web socket client",
		),
		[]string{"creation_type"},
	)

	WebSocketClientCreationSuccess = prometheus.NewCounterVec(
		newCounterOpts(
			"websocket_client_creation_successes_total",
			"Total number of accepting and successful creation of a web socket client",
		),
		[]string{"creation_type"},
	)

	WebSocketMessageBytesSent = prometheus.NewCounter(
		newCounterOpts(
			"websocket_message_bytes_sent_total",
			"Total number of bytes sent to web socket clients",
		),
	)
)

func (socketCollectors) Collect() []prometheus.Collector {
	return []prometheus.Collector{
		WebSocketClientConnectionClose,
		WebSocketClientMessageSendSuccess,
		WebSocketClientMessageSendFailure,
		WebSocketClientCreationFailure,
		WebSocketClientCreationSuccess,
		WebSocketMessageBytesSent,
	}
}
