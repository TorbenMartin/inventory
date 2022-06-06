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
	err := http.ListenAndServeTLS(":50000", "certs/server.crt", "certs/server.key", nil)
	if err != nil {
		panic(err.Error())
	}

	

}

