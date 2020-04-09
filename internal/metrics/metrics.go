package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	go func() {
		if err := Serve(":9000"); err != nil {
			fmt.Println(err)
		}
	}()
}

func Serve(addr string) error {

	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}

	return nil
}
