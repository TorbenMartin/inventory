package main

import (
	"fmt"
	"net/http"
	"time"
	"net"
)

//////////////////logout function//////////////////
func logoutHandler(w http.ResponseWriter, r *http.Request) {

	secheader(w)
	ip,_,_ := net.SplitHostPort(r.RemoteAddr)
	agent := r.Header.Get("User-Agent")
	t := time.Now()
	formatted := fmt.Sprintf("%02d.%02d.%d", t.Day(), t.Month(), t.Year())
	tokena := `l`+md5hash(formatted + " " + ip + " " + agent)+``

	c := http.Cookie{
               Name:   tokena,
               MaxAge: -1}
	http.SetCookie(w, &c)
         
	deletetoken()



	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<meta http-equiv=\"refresh\" content=\"0; URL=/\">")
	
}
