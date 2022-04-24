package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)


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
		log.Fatal("ListenAndServe: ", err)
	}

}
