package pipe

import (
	"fmt"
	"io"
)

type Server struct{}

func (s *Server) Ingest(srv Pipe_IngestServer) error {

	for {
		data, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		fmt.Println("data received:", data.GetData())
	}
}
