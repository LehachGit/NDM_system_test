package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
)

func headersHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	xff := r.Header.Get("X-Forwarded-For")
	remoteHost, _, _ := strings.Cut(r.RemoteAddr, ":")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprintf(w, "=== XFF Demo App  hostname=%s ===\n\n", hostname)
	fmt.Fprintf(w, "%-30s %s\n", "Remote-Addr (TCP):", r.RemoteAddr)
	fmt.Fprintf(w, "%-30s %s\n", "X-Real-IP:", r.Header.Get("X-Real-IP"))
	fmt.Fprintf(w, "%-30s %s\n", "X-Forwarded-For:", xff)

	if xff != "" {
		hops := strings.Split(xff, ",")
		fmt.Fprintf(w, "\n--- X-Forwarded-For chain (%d hop(s) + last nginx) ---\n", len(hops))
		for i, ip := range hops {
			label := "proxy"
			if i == 0 {
				label = "CLIENT — original user IP"
			}
			fmt.Fprintf(w, "  hop %-2d  %-18s  %s\n", i+1, strings.TrimSpace(ip), label)
		}
		fmt.Fprintf(w, "  hop %-2d  %-18s  last nginx (TCP)\n", len(hops)+1, remoteHost)
	}

	fmt.Fprintf(w, "\n--- All headers ---\n")
	names := make([]string, 0, len(r.Header))
	for name := range r.Header {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Fprintf(w, "  %-30s %s\n", name+":", strings.Join(r.Header[name], ", "))
	}
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})
	http.HandleFunc("/", headersHandler)

	fmt.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
