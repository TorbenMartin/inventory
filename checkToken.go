package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"
	"net"
)


//////////////////update login token function//////////////////
func updatetoken(token string, userid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	change, err := db.Prepare("UPDATE login SET session= ?, sessiontime=NOW() where id= ? ")
	if err != nil {
		panic(err.Error())
	}

	defer change.Close()

	if _, err := change.Exec(token,userid); err != nil {
		panic(err.Error())
	}
}


//////////////////delete admin login token after X minutes function//////////////////
func deletetoken() {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	change, err := db.Query("UPDATE login SET session=\"\" where rechte = 1 and sessiontime < (NOW() - interval 3 minute) ")
	if err != nil {
		panic(err.Error())
	}

	defer change.Close()
}

//////////////////checktoken function//////////////////
func checktoken(w http.ResponseWriter, r *http.Request) string{

	ip,_,_ := net.SplitHostPort(r.RemoteAddr)
	agent := r.Header.Get("User-Agent")
	t := time.Now()
	formatted := fmt.Sprintf("%02d.%02d.%d", t.Day(), t.Month(), t.Year())
	tokena := `l`+md5hash(formatted + " " + ip + " " + agent)+``
	
		var (
			id string
			rechte string
			username string
			password string
			l bool
		)
      		c, err := r.Cookie(tokena)
        	if err != nil {
        	} else {

         	value := c.Value
         	splitsting := strings.Split(value, "_")    	         	
		db, err := sql.Open("mysql", sqlcred)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()    	
         	err = db.QueryRow("SELECT id,rechte,username,password FROM login where id = ?", splitsting[0]).Scan(&id,&rechte,&username,&password)
		if err != nil {
			panic(err.Error())
		}
			
		if (id != "")  {
			token := md5hash(formatted + " " + ip + " " + agent + " " + username + " " + password)
			if (token == splitsting[1]) {
			
			
			c := http.Cookie{
			Name:   tokena,
			Value:  value,
			MaxAge: 3600}
			http.SetCookie(w, &c)
			
			updatetoken(token, splitsting[0])
			deletetoken()

				l = true
			} else {
				deletetoken()
				logoutHandler(w, r)
				l = false
			}
		}
         }
if l == true {
return ``+id+`_`+rechte+``
} else {
return `0_0`
}

}      
