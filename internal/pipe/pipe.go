package pipe

import (
	"io"

	"github.com/jrieck1991/mark/internal/metrics"
)

type Server struct{}

const (
	// metrics
	namespace  string = "app"
	subsystem  string = "server"
	msgRecv    string = "message_recv"
	bytesRecv  string = "bytes_recv"
	recvErr    string = "recv_error"
	recvErrEOF string = "recv_error_eof"
)

// Ingest data
func (s *Server) Ingest(srv Pipe_IngestServer) error {

	// init metric counters
	counters := metrics.Counters(namespace, subsystem, []string{
		msgRecv, bytesRecv, recvErr, recvErrEOF,
	})

	// receive data
	for {

		data, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				counters[recvErrEOF].Inc()
				continue
			}
			counters[recvErr].Inc()
			continue
		}

		// record message sent
		counters[msgRecv].Inc()

		// record bytes
		counters[bytesRecv].Add(float64(len(data.GetData())))
	}
}
