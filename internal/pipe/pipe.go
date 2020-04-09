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
	msgSent    string = "message_sent"
	bytesSent  string = "bytes_sent"
	recvErr    string = "recv_error"
	recvErrEOF string = "recv_error_eof"
)

func (s *Server) Ingest(srv Pipe_IngestServer) error {

	counters := metrics.Counters(namespace, subsystem, []string{
		msgSent, bytesSent, recvErr, recvErrEOF,
	})

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
		counters[msgSent].Inc()

		// record bytes
		counters[bytesSent].Inc()
	}
}
