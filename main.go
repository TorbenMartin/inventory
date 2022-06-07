package main

import (
	"net/http"
	"os"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

// cipher key for user password in database
const keysalt string = "drfvert5z5434t54654ewqr3465wrcf5v4tbnzujrrwcrc3456zer"


// cipher key 32 char for database password in config file
const key string = "q-BkdGka2YpGd2,eBb=-Ab5?qhd:M-7'"

//////////////////main process function//////////////////
func main() {

	fileName := "requests.log"
	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	defer logFile.Close()
	log.SetOutput(logFile)

	mux := http.NewServeMux()
	mux.HandleFunc("/", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/user", userHandler)
	mux.HandleFunc("/admin", adminHandler)
	mux.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("./img"))))

	http.ListenAndServeTLS(":50000", "certs/cert.pem", "certs/privkey.pem", RequestLogger(mux))

}

func secheader(w http.ResponseWriter) {

	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("cache-control", "private, max-age=0")
	w.Header().Set("bfcache-opt-in", "unload")
	w.Header().Set("server", "gws")
	w.Header().Set("expires", "-1")
	w.Header().Set("strict-transport-security", "max-age=31536000; includeSubDomains")
	w.Header().Set("x-frame-options", "sameorigin")

}

func RequestLogger(targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		targetMux.ServeHTTP(w, r)

		// log request by who(IP address)
		requesterIP := r.RemoteAddr

		log.Printf(
			"%s\t%s\t%s\t\t\t%s\t%s\t%s\t%s\t",
			r.Method,
			requesterIP,
			r.Header,
			r.TLS,
			r.Form,
			r.URL,
		)
	})
}
