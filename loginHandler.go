package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"net"
	"strconv"
	"io/ioutil"
        "encoding/json"
)

// declaring a struct
type jDaten struct {
	Success     bool
	ChallengeTS string
	Hostname    string
	Score       float64
	Action      string
}


//////////////////get recapv3 function//////////////////
func OnPage(link string)(string) {

	res, err := http.Get(link)
	if err != nil {
		return "0"
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "0"
	}

	var daten jDaten
	Data := []byte(string(content))
	errb := json.Unmarshal(Data, &daten) 
	if errb != nil {
		return "0"
	}

	if (daten.Action == "login" && daten.Success == true && daten.Score > scorecap) {
		return "1"
	} else{
		return "0"
	}

return "0"
}

//////////////////login function//////////////////
func loginHandler(w http.ResponseWriter, r *http.Request) {
	
	secheader(w)
	if r.Method == "POST" {
		r.ParseForm()

		var linkcap string = `https://www.google.com/recaptcha/api/siteverify?secret=`+pivcap+`&response=`+r.FormValue("recaptcha_token")+``

		if (r.FormValue("login") == "Login" && r.FormValue("password") != "" && r.FormValue("username") != "" && OnPage(linkcap) == "1") {

			db, err := sql.Open("mysql", sqlcred)
			if err != nil {
				panic(err.Error())
			}
			defer db.Close()
			var (
				id       int
				username string
				password string
				rechte1  int
				aktiv    int
			)
			err = db.QueryRow("SELECT id, username, password, rechte, aktiv FROM login where username = ? ", r.FormValue("username")).Scan(&id, &username, &password, &rechte1, &aktiv)
			if err != nil {
				w.Header().Set("Content-Type", "text/html")
				fmt.Fprintln(w, "<meta http-equiv=\"refresh\" content=\"0; URL=/\">")
			}
			if ((r.FormValue("username") == username) && (md5hash(``+r.FormValue("password")+``+keysalt+``) == password)) {
				
				ip,_,_ := net.SplitHostPort(r.RemoteAddr)
				agent := r.Header.Get("User-Agent")
				t := time.Now()
				formatted := fmt.Sprintf("%02d.%02d.%d", t.Day(), t.Month(), t.Year())
				
				tokena := `l`+md5hash(formatted + " " + ip + " " + agent)+``
				token := ``
				uidstring := ``

				token = md5hash(formatted + " " + ip + " " + agent + " " + username + " " + password)
				uidstring = strconv.Itoa(id)				

				updatetoken(token,uidstring)			

				c := http.Cookie{
				Name:   tokena,
				Value:  ``+uidstring+`_`+token+``,
				MaxAge: 3600}
				http.SetCookie(w, &c)


				deletetoken()
				
				w.Header().Set("Content-Type", "text/html")
				fmt.Fprintln(w, "<meta http-equiv=\"refresh\" content=\"0; URL=/user\">")				
					
			}
		}
	}
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintln(w, globalmeta)
	fmt.Fprintln(w, `
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		
		<script src="https://www.google.com/recaptcha/api.js?render=`+pubcap+`"></script>
		<script>
			grecaptcha.ready(function() {
			grecaptcha.execute('`+pubcap+`', {action: 'login'}).then(function(token) {
			document.getElementById('recaptcha_token').value = token;
				});
			});
		</script>
	`)

	fmt.Fprintln(w, `
	
		<style>
		
			a:link {
				color: gray;
			}

			a:visited {
				color: gray;
			}

			a:hover {
				color: 008000;
			}

			a:active {
				color: gray;
			}


			body {
				margin: 0;
				padding: 0;
				background-image: url("./img/bg.jpeg"); //add "" if you want
				background-repeat: no-repeat;
				background-attachment: fixed;
				background-size: cover;
				color: black;
			}
			.glass{
				background: linear-gradient(135deg, rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0));
				backdrop-filter: blur(10px);
				-webkit-backdrop-filter: blur(10px);
				border-radius: 10px;
				border:1px solid rgba(255, 255, 255, 0.18);
				box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.37);

				height: 160px;
			}
		</style>
	

		<!-- 
		Dieses Programm hat Torben MARTIN, im Zusammenhang einer Umschulung zum IT Administrator erstellt.
		Programmiert habe ich dieses Projekt zu Hause in meiner Freizeit aus eigenem Antrieb herraus.
		 -->
		 
		
		<table border="0" width="100%" height="100%">
		<tr>
		<td>
		<form action="/" method="post">
		<input type="hidden" name="recaptcha_token" id="recaptcha_token">
		<center><div class="glass">
		<table border="0" ><br>
    			<tr>
        			<td >
           			 <table border="0" >
            				<tr>
            				    	<td width="70">Username:</td>
						<td><input type="text" name="username" autofocus></td>
            				</tr>
            				<tr>
                				<td width="70">Passwort:</td>
						<td><input type="password" name="password"></td>
            				</tr>
            				<tr>
                				<td width="70"></td>
						<td align="right">
						<center>
							<input type="submit" name="login" value="Login">
							<br><br>
						</center>
							<a href="/user">[Statistik]</a><a href="/user?txt=1">[txt]</a>
						
						</td>
            				</tr>
           			 </table>
        			</td>
    			</tr>
		</table>
		</div>
		</center>	
		</form>
		</td>
		</tr>
		</table>
		
	`)
}

