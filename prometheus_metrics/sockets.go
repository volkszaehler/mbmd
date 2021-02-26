package prometheus_metrics

import "github.com/prometheus/client_golang/prometheus"

// socketCollectors contains all Prometheus metrics about web sockets
//
// Implements collectable interface
type socketCollectors struct {}

var (
	WebSocketClientConnectionClose = prometheus.NewCounter(
		*newCounterOpts(
			"websocket_client_connection_closed_total",
			"Total amount of closed client connection to a web socket",
		),
	)

	WebSocketClientMessageSendSuccess = prometheus.NewCounter(
		*newCounterOpts(
			"websocket_client_message_send_success_total",
			"Total amount of messages sent to a web socket client",
		),
	)

	WebSocketClientMessageSendFailure = prometheus.NewCounter(
		*newCounterOpts(
			"websocket_client_message_send_failure_total",
			"",
		),
	)

	WebSocketClientCreationFailure = prometheus.NewCounterVec(
		*newCounterOpts(
			"websocket_client_creation_failure_total",
			"",
		),
		[]string{"creation_type"},
	)

	WebSocketClientCreationSuccess = prometheus.NewCounterVec(
		*newCounterOpts(
			"websocket_client_creation_success_total",
			"",
		),
		[]string{"creation_type"},
	)

	WebSocketMessageBytesSent = prometheus.NewCounter(
		*newCounterOpts(
			"websocket_message_bytes_sent_total",
			"",
		),
	)
)

func (c socketCollectors) Collect() []prometheus.Collector {
	return []prometheus.Collector{
		WebSocketClientConnectionClose,
		WebSocketClientMessageSendSuccess,
		WebSocketClientMessageSendFailure,
		WebSocketClientCreationFailure,
		WebSocketClientCreationSuccess,
		WebSocketMessageBytesSent,
	}
}
