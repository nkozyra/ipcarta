package ipcarta

import (
	"fmt"
	"net/http"
)

func Serve() {
	http.HandleFunc("/ips", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("ip")
		found, msg := Search(q)
		if found {
			out, _ := msg.MarshalJSON()
			fmt.Fprintf(w, string(out))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	http.ListenAndServe(":9999", nil)
}
