package main

import (
	"net/http"
	_ "github.com/go-sql-driver/mysql"
)

// cipher key for user password in database
const keysalt string = "drfvert5z5434t54654ewqr3465wrcf5v4tbnzujrrwcrc3456zer"


// cipher key 32 char for database password in config file
const key string = "q-BkdGka2YpGd2,eBb=-Ab5?qhd:M-7'"

//////////////////main process function//////////////////
func main() {

	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/admin", adminHandler)
	http.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("./img"))))

	//err := http.ListenAndServe(":50000", nil)
	err := http.ListenAndServeTLS(":50000", "certs/cert.pem", "certs/privkey.pem", nil)
	if err != nil {
		panic(err.Error())
	}

	

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
