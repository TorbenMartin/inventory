package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"net"
        "encoding/json"
	"os"

	_ "github.com/go-sql-driver/mysql"
)


/*



(1) download the mysql driver:

	- go get github.com/go-sql-driver/mysql


(2) config file example "conf.json":

	{
	   "User": ["username"],
 	   "Password": ["password"],
 	   "Server": ["localhost"],
 	   "Port": ["3306"],
	   "Db": ["lager"]
	}


(3) create and install on mysql database:

	create database lager;
	use lager;

	CREATE TABLE bestand (
		id int NOT NULL AUTO_INCREMENT,
		gertyp varchar(255) NOT NULL,
		modell varchar(255) NOT NULL,
		seriennummer varchar(255) NOT NULL,
		zinfo varchar(255),
		ticketnr varchar(255),
		ausgabename varchar(255),
		ausgabedatum date,
		changed int,
		einkaufsdatum date,
		PRIMARY KEY (id)
	);


	CREATE TABLE zinfodata (
		id int NOT NULL AUTO_INCREMENT,
		zinfoid varchar(255) NOT NULL,
		bestandid varchar(255) NOT NULL,
		daten varchar(255) NOT NULL,
		PRIMARY KEY (id)
	);


	CREATE TABLE modell (
		id int NOT NULL AUTO_INCREMENT,
		modell varchar(255) NOT NULL,
		sperrbestand varchar(255),
		sort varchar(255),
		PRIMARY KEY (id)
	);

	CREATE TABLE zinfo (
		id int NOT NULL AUTO_INCREMENT,
		gertyp varchar(255) NOT NULL,
		zinfoname varchar(255) NOT NULL,
		sort varchar(255),
		PRIMARY KEY (id)
	);

	CREATE TABLE gertyp (
		id int NOT NULL AUTO_INCREMENT,
		gertyp varchar(255) NOT NULL,
		sort varchar(255),
		PRIMARY KEY (id)
	);


	CREATE TABLE login (
		id int NOT NULL AUTO_INCREMENT,
		username varchar(255) NOT NULL,
		password varchar(255) NOT NULL,
		rechte int,
		aktiv int NOT NULL,
		session varchar(255),
		sessiontime timestamp,
		PRIMARY KEY (id)
	);


	INSERT INTO login (username, password, rechte, aktiv) VALUES ("admin","21232f297a57a5a743894a0e4a801fc3","1","1");

(4) compile:

	- for linux: go build -o lager main.go
	- for windows: GOOS=windows GOARCH=amd64 go build -o lager.exe main.go

(5) normal run:

	- run source: go run main.go
	
	- run binary on linux:
		chmod +x lager
		./lager
	
	-run binary on windows:
		just click it with your fkn mouse

	- open webbrowser on http://localhost:50000

(6) run as a linux service:
	
	cp lager /home/username
	chmod +x /home/username
	
	cat >> /etc/systemd/system/lager.service << EOF
	[Unit]
	Description=Service 1
	DefaultDependencies=no
	#After=network.target
	#Wants=network-online.target systemd-networkd-wait-online.service

	StartLimitIntervalSec=500
	StartLimitBurst=5

	[Service]
	Restart=on-failure
	RestartSec=5s
	WorkingDirectory= /home/username

	Type=simple
	User=username
	Group=username
	ExecStart= /home/username/./lager
	TimeoutStartSec=0
	RemainAfterExit=yes

	[Install]
	WantedBy=default.target
	EOF	
	
	systemctl status lager
	systemctl start lager
	systemctl status lager
	
*/


//////////////////////////////////////////////////////////////////////////////////////
//				source start                                       //
//////////////////////////////////////////////////////////////////////////////////////

// config
type Configuration struct {
    User    []string
    Password   []string
    Server   []string
    Port   []string
    Db   []string
}

var sqlcred = getconfig()

func getconfig() string {
file, _ := os.Open("conf.json")
defer file.Close()
decoder := json.NewDecoder(file)
configuration := Configuration{}
errcode := decoder.Decode(&configuration)
if errcode != nil {
  fmt.Println("error:", errcode)
}

return ``+strings.Join(configuration.User," ")+`:`+strings.Join(configuration.Password," ")+`@tcp(`+strings.Join(configuration.Server," ")+`:`+strings.Join(configuration.Port," ")+`)/`+strings.Join(configuration.Db," ")+``
}


//////////////////main process function//////////////////
func main() {

	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/admin", adminHandler)
	http.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("./img"))))

	err := http.ListenAndServe(":50000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

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
         	err = db.QueryRow("SELECT id,rechte,username,password FROM login where id = \"" + splitsting[0] + "\"").Scan(&id,&rechte,&username,&password)
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
 

   
//////////////////user site function//////////////////
func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	deletetoken()

	checktokenstring := strings.Split(checktoken(w,r), "_")    

	if (checktokenstring[0] != "0") {

		db, err := sql.Open("mysql", sqlcred)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		///post start
		if r.Method == "POST" {
			r.ParseForm()

			if (r.FormValue("additem") == "Hinzufügen" && r.FormValue("sn") != "" && r.FormValue("modell") != "" && r.FormValue("einkaufsdatum")!="") {
			
				var (
					sncount    int
					zinfoid    string
					bestandid    string
				)


				err = db.QueryRow("SELECT COUNT(*) FROM bestand where seriennummer = \"" + r.FormValue("sn") + "\" ").Scan(&sncount)
				if err != nil {
				}


				if (sncount == 0) {
					sqlinsert(r.URL.Query().Get("gertyp"), r.FormValue("modell"), r.FormValue("sn"), r.URL.Query().Get("gertyp"), r.FormValue("einkaufsdatum"), checktokenstring[0])


					err = db.QueryRow("SELECT id FROM bestand where gertyp=\""+r.URL.Query().Get("gertyp")+"\" and seriennummer = \"" + r.FormValue("sn") + "\" ").Scan(&bestandid)
					if err != nil {
					}
					
					//loadzinfo1 start					
					if (bestandid != "") {
						rows7, err := db.Query("SELECT id FROM zinfo")
						if err != nil {
							log.Fatal(err)
						}
						defer rows7.Close()
						for rows7.Next() {
						err := rows7.Scan(&zinfoid)
								if err != nil {
								log.Fatal(err)
							}

							if (r.FormValue("zinfo"+zinfoid) !="") {
								zinfodatainput(zinfoid, bestandid, r.FormValue("zinfo"+zinfoid))
							}
						
						}
					}		
					//loadzinfo1 end
					


				} else {
					fmt.Fprintln(w, `
				<script>alert("Fehler -> Gerät schon vorhanden!!!");</script>
			`)
				}
			}

			if checktokenstring[1] == "1" {
				if (r.FormValue("additem") == "Speichern" && r.FormValue("sn") != "" && r.FormValue("changeitemid") != "" && r.FormValue("ausgegebenan") != "" && r.FormValue("ausgegebenam") != "" && r.FormValue("ticketnr") != "" && r.FormValue("einkaufsdatum")!="" && r.FormValue("modell") != "") {
					sqledit(r.FormValue("modell"), r.FormValue("changeitemid"), r.URL.Query().Get("gertyp"), r.FormValue("sn"), r.FormValue("ausgegebenan"), r.FormValue("ausgegebenam"), r.FormValue("ticketnr"), r.FormValue("einkaufsdatum"),checktokenstring[0])
				
				
				
					var (
					zinfoid    string
					bestandid    string
				)
			
				
				
				
					err = db.QueryRow("SELECT id FROM bestand where seriennummer = \"" + r.FormValue("sn") + "\" ").Scan(&bestandid)
					if err != nil {
					}
					
					//loadzinfo1 start					
					if (bestandid != "") {
					
						rows9, err := db.Query("SELECT id FROM zinfo where gertyp=\"" + r.URL.Query().Get("gertyp") + "\" ")
						if err != nil {
							log.Fatal(err)
						}
						defer rows9.Close()
						for rows9.Next() {
						err := rows9.Scan(&zinfoid)
								if err != nil {
								log.Fatal(err)
							}

							if (r.FormValue("zinfo"+zinfoid) !="") {
							
								zinfodataedit(zinfoid, bestandid, r.FormValue("zinfo"+zinfoid))
							}
						
						}
					}		
					//loadzinfo1 end				
				
		
				}
			} else {
				if (r.FormValue("additem") == "Speichern" && r.FormValue("changeitemid") != "" && r.FormValue("ausgegebenan") != "" && r.FormValue("ausgegebenam") != "" && r.FormValue("ticketnr") != "") {
					sqledit2(r.FormValue("changeitemid"), r.URL.Query().Get("gertyp"), r.FormValue("ausgegebenan"), r.FormValue("ausgegebenam"), r.FormValue("ticketnr"),checktokenstring[0])

				}
			}

			if r.FormValue("additem") == "zurück" && r.FormValue("getbackid") != "" {
				getback(r.FormValue("getbackid"),checktokenstring[0])
			}

			if checktokenstring[1] == "1" {
				if r.FormValue("additem") == "löschen" && r.FormValue("delitemid") != "" {
					delitem(r.FormValue("delitemid"))
				}
			}

		}
		///post end


	if r.URL.Query().Get("txt") != "1" {
		fmt.Fprintln(w, globalmeta)

		///style start
		fmt.Fprintln(w, style)
		///style end
	}


		///menu1 start

		fmt.Fprintln(w, menustart)

		if checktokenstring[1] == "1" {
			fmt.Fprintln(w, `<li><a href="/admin">[Adminbereich]</a></li>`)
		}

		fmt.Fprintln(w, `<li><a href="/user">[Statistik]</a></li>`)
		fmt.Fprintln(w, `<li><a href="/logout"><span style="color:red">[Logout]</span> <span id="timer" style="font-weight: bold;color:red"></span></a></li>`)

		fmt.Fprintln(w, menuend)
		///menu1 end

		///load vars start
		var (
			id           string
			modell       string
			gertyp       string
			seriennummer string
			zinfo        string
			ausgabename  string
			ausgabedatum string
			changed      string
			ticketnr     string
			gertypcount  string
			zinfoid  string
			zinfoname  string
			gertypid  string
			fsize int
			daten string
			einkaufsdatum string
		)
		///load vars end

		///menu2 start
		fmt.Fprintln(w, menustart)

		rows, err := db.Query("SELECT id, gertyp FROM gertyp")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &gertyp)
			if err != nil {
				log.Fatal(err)
			}

			if id != "" && gertyp != "" {
				fmt.Fprintln(w, `<li><a href="/user?gertyp=`+id+`">`+gertyp+`</a></li>`)
			}
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintln(w, menuend)
		///menu2 end

		///logo start
		fmt.Fprintln(w, `
<br><center><img src="./img/logo.png" height="60"></center>
`)
		///logo end

		///content start

		if len(r.URL.Query().Get("gertyp")) > 0 {


	

			err = db.QueryRow("SELECT COUNT(*) FROM gertyp where id = \"" + r.URL.Query().Get("gertyp") + "\" ").Scan(&gertypcount)
			if err != nil {
			}

			if gertypcount != "0" {

				///search start
				if r.URL.Query().Get("gertyp") != "" {
					fmt.Fprintln(w, `
		<br><br>
		<table width="95%">
			<tr width="100%" align="right">
				<td>
					<form action="/user?gertyp=`+r.URL.Query().Get("gertyp")+`" method="post">
						<input required type="text" size="21" name="search" placeholder="YYYY-MM-DD oder Suchwort" value="">
						<input type="submit" name="searchitem" id="searchitem" value="suchen" style="background-color: lightblue">
					</form>
				</td>
			</tr>
		</table>
	`)
				}
				///search end




				///additem block start
				fmt.Fprintln(w, `
	<center>
<details name="addmen" id="addmen" close>
  <summary></summary><br>
	<table border="0" width="80%">
		<tr>
			<td valign="top" align="left">`)

				fmt.Fprintln(w, `			
				<center>


				<form action="/user?gertyp=`+r.URL.Query().Get("gertyp")+`" method="post">
				<input type="hidden" name="changeitemid" id="changeitemid" value="">
				<table>
					<tr>
						<td>Modell</td>
						<td>
						<select required name="modell" id="modell" width="180" style="width: 180px">
			`)

				rows2, err := db.Query("SELECT id, modell FROM modell")
				if err != nil {
					log.Fatal(err)
				}
				defer rows2.Close()
				for rows2.Next() {
					err := rows2.Scan(&id, &modell)
					if err != nil {
						log.Fatal(err)
					}

					if r.FormValue("modell") != "" {
						if r.FormValue("modell") == id {
							fmt.Fprintln(w, `							<option name="modellid`+id+`" value="`+id+`" required selected>`+modell+`</option>`)
						} else {
							fmt.Fprintln(w, `							<option name="modellid`+id+`" value="`+id+`" required>`+modell+`</option>`)
						}
					} else {
						fmt.Fprintln(w, `							<option name="modellid`+id+`" value="`+id+`" required>`+modell+`</option>`)
					}

				}
				err = rows2.Err()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(w, `
						</select>
						</td>
					</tr>
					<tr><td> Seriennummer </td><td><input type="text" name="sn" id="sn"  required tabIndex="1" autofocus></td></tr>`)
			
					//loadzinfo start
					
					var tabcount int = 1
					
					rows6, err := db.Query("SELECT id, gertyp ,zinfoname FROM zinfo")
					if err != nil {
						log.Fatal(err)
					}
					defer rows6.Close()
					for rows6.Next() {
					err := rows6.Scan(&zinfoid, &gertypid, &zinfoname)
						if err != nil {
							log.Fatal(err)
						}
						
						if (gertypid == r.URL.Query().Get("gertyp")) {
						tabcount = (tabcount + 1) 
						tab := strconv.Itoa(tabcount)
							fmt.Fprintln(w, `					<tr><td> `+ zinfoname +` </td><td><input type="text" name="zinfo`+ zinfoid +`" id="zinfo`+ zinfoid +`" required tabIndex="`+tab+`"></td></tr>`)
						}
						
					}			
					//loadzinfo end

					t := time.Now()
					formatted := fmt.Sprintf("%02d.%02d.%d", t.Day(), t.Month(), t.Year())
					
					
					if r.FormValue("einkaufsdatum") != "" {
					fmt.Fprintln(w, `<tr><td> Einkaufsdatum </td><td><input type="text" name="einkaufsdatum"  id="einkaufsdatum" placeholder="DD/MM/YYYY" pattern="(?:((?:0[1-9]|1[0-9]|2[0-9])\.(?:0[1-9]|1[0-2])|(?:30)\.(?!02)(?:0[1-9]|1[0-2])|31\.(?:0[13578]|1[02]))\.(?:19|20)[0-9]{2})" value="`+r.FormValue("einkaufsdatum")+`" required></td></tr>`)					
					} else {
					fmt.Fprintln(w, `					<tr><td> Einkaufsdatum </td><td><input type="text" name="einkaufsdatum"  id="einkaufsdatum" placeholder="DD/MM/YYYY" pattern="(?:((?:0[1-9]|1[0-9]|2[0-9])\.(?:0[1-9]|1[0-2])|(?:30)\.(?!02)(?:0[1-9]|1[0-2])|31\.(?:0[13578]|1[02]))\.(?:19|20)[0-9]{2})" value="`+formatted+`" required></td></tr>`)					
					}
					
					
					fmt.Fprintln(w, `					<tr><td> ausgegeben an </td><td><input type="text" name="ausgegebenan"  id="ausgegebenan" disabled></td></tr>
					<tr><td> Ticket Nummer </td><td><input type="text" name="ticketnr"  id="ticketnr" disabled></td></tr>
					<tr><td> ausgegeben am </td><td><input type="text" name="ausgegebenam"  id="ausgegebenam" disabled></td></tr>
					<tr><td></td><td><input type="submit" name="additem" id="additem" value="Hinzufügen"> <input type="submit" name="additem" id="edititem" value="Speichern" disabled> <button onclick="cancel()">Abbrechen</button></td></tr>
				</table>
				</form>
				</center>
			`)

				fmt.Fprintln(w, `
			</td>
		</tr>
	</table>
</details>
	</center>
	`)


				if (r.FormValue("additem") == "Hinzufügen") || (r.FormValue("changeitemid") != "") {
					fmt.Fprintln(w, `
			<script>
			document.getElementById('addmen').open = true;
			document.getElementById('sn').focus();

			</script>
			`)
				}
				///additem block end
			var (
			sitea int
			siteb int
			siteaa string
			sitebb string
			)
			
			sitea = 0
			siteb = 1
			
			if(r.URL.Query().Get("site") != "") {

					sitea, _ = strconv.Atoi(r.URL.Query().Get("site"))
					sitea = sitea - 1
					
					if (sitea < 0) {
					sitea = 0
					}
					
					siteb, _ = strconv.Atoi(r.URL.Query().Get("site"))
					siteb = siteb + 1
					
					if (siteb < 0) {
					siteb = 1
					}

					
			}


			siteaa = strconv.Itoa(sitea) 
			sitebb = strconv.Itoa(siteb) 
		
			fmt.Fprintln(w, `<br><br><center><table border="0" width= "80%" ><tr align="right"><td>
			
			
			<a href="./user?gertyp=`+r.URL.Query().Get("gertyp")+`&rows=`+r.URL.Query().Get("rows")+`&site=`+siteaa+`" target="_self"> [ < ] </a>
			<a href="./user?gertyp=`+r.URL.Query().Get("gertyp")+`&rows=`+r.URL.Query().Get("rows")+`&site=`+sitebb+`" target="_self"> [ > ] </a>
			
			<select onChange="window.open(this.value, '_self')">
				<option value =""></option>
				<option value ="./user?gertyp=`+r.URL.Query().Get("gertyp")+`&rows=&site=">*</option> 
				<option value="./user?gertyp=`+r.URL.Query().Get("gertyp")+`&rows=10&site=`+r.URL.Query().Get("site")+`">10</option>
				<option value="./user?gertyp=`+r.URL.Query().Get("gertyp")+`&rows=30&site=`+r.URL.Query().Get("site")+`">30</option>
				<option value="./user?gertyp=`+r.URL.Query().Get("gertyp")+`&rows=50&site=`+r.URL.Query().Get("site")+`">50</option>
				<option value="./user?gertyp=`+r.URL.Query().Get("gertyp")+`&rows=100&site=`+r.URL.Query().Get("site")+`">100</option>
			</select>
			
			</td></tr></table></center>`)	
			
			
				
				///table bestand data start
				var changeheader []string
				var changeheaderb []string
				fmt.Fprintln(w, `<br>`)

				mid := r.URL.Query().Get("gertyp")
				if len(mid) > 0 {

					

					
		fmt.Fprintln(w, `					
		<center>

		<table border="0" id="foobar-table">
			<thead>

				<tr>	
					<th valign="top" align="center" data-type="string" width="100"> `+getgertypinfo(r.URL.Query().Get("gertyp"))+`: </th>
					<th valign="top" align="center" data-type="number" width="140"> Seriennummer: </th>`)
					
					
					//loadzinfo start
					rows8, err := db.Query("SELECT id, zinfoname FROM zinfo where gertyp= \"" + r.URL.Query().Get("gertyp") + "\" ")
					if err != nil {
						log.Fatal(err)
					}
					defer rows8.Close()
					
					var zinfoarray []string
					
					for rows8.Next() {
					err := rows8.Scan(&zinfoid, &zinfoname)
						if err != nil {
							log.Fatal(err)
						}
						
							if (zinfoname != "") {												
								fsize = 18 + (len(zinfoname) * 10)							
								fsizeb := strconv.Itoa(fsize)

								fmt.Fprintln(w, `					<th valign="top" align="center" data-type="string" width="`+ fsizeb +`"> `+ zinfoname +`: </th>`)
								zinfoarray = append(zinfoarray, zinfoid)
							}
						
					}			
					//loadzinfo end
					
					
					
					fmt.Fprintln(w, `					<th valign="top" align="center" data-type="string" width="150"> Einkaufsdatum: </th>
					<th valign="top" align="center" data-type="string" width="150"> ausgegeben an: </th>
					<th valign="top" align="center" data-type="string" width="100"> Ticket NR: </th>
					<th valign="top" align="center" data-type="number" width="120"> geändert am: </th>
					<th valign="top" align="center" data-type="string" width="183"> Optionen: </th>
				</tr>

			</thead>
		<tbody>
		`)
					///loop bestand start

					//search start output
					
					if r.FormValue("searchitem") == "suchen" && trimHtml(r.FormValue("search")) != "" {
					
					var searchcount int
					mid := r.URL.Query().Get("gertyp")
					err = db.QueryRow("SELECT COUNT(*) FROM zinfo where gertyp = \"" + mid + "\" ").Scan(&searchcount)
					
					if (searchcount > 0) {

					search := trimHtml(r.FormValue("search"))
										
					sqlq := `SELECT DISTINCT bestand.id, bestand.gertyp, bestand.modell, bestand.seriennummer, bestand.zinfo, bestand.ausgabename, bestand.ausgabedatum, bestand.changed, bestand.ticketnr, bestand.einkaufsdatum
						FROM bestand 
						JOIN zinfo ON bestand.gertyp
						JOIN zinfodata ON bestand.id
						where zinfodata.bestandid = bestand.id and bestand.gertyp = `+mid+` and (zinfodata.daten like '%`+search+`%' or bestand.seriennummer like '%`+search+`%' or bestand.ausgabename like '%`+search+`%' or bestand.ausgabedatum like '%`+search+`%' or bestand.ticketnr like '%`+search+`%' or bestand.einkaufsdatum like '%`+search+`%') ;`

						rows3, err := db.Query(sqlq)
						
						if err != nil {
							log.Fatal(err)
						}

						for rows3.Next() {
							err := rows3.Scan(&id, &gertyp, &modell, &seriennummer, &zinfo, &ausgabename, &ausgabedatum, &changed, &ticketnr, &einkaufsdatum)
							if err != nil {
								log.Fatal(err)
							}

							fmt.Fprintln(w, `
	<tr>
		<td valign="top" align="center">`+getmodinfo(modell)+`</td>
		<td valign="top" align="center">`+seriennummer+`</td>

	`)

///zinfo table start



		//loadzinfo start
		for index := range zinfoarray {
		daten=""						
		err = db.QueryRow("SELECT daten FROM zinfodata where bestandid = \"" + id + "\" and zinfoid = \"" + zinfoarray[index] + "\" ").Scan(&daten)
		if err != nil {
		}
								
			if (daten !="") {
				fmt.Fprintln(w, `<td valign="top" align="center">`+ daten +`</td>`)
			} else {
				fmt.Fprintln(w, `<td valign="top" align="center"></td>`)
			}
		}
									
		//loadzinfo end


	split := strings.Split(einkaufsdatum, "-")
	einkaufsdatumb := ``+split[2]+`.`+split[1]+`.`+split[0]+``
	
	splitb := strings.Split(ausgabedatum, "-")
	ausgabedatumb := ``+splitb[2]+`.`+splitb[1]+`.`+splitb[0]+``

///zinfo table end
							fmt.Fprintln(w, `
		<td valign="top" align="center">`+einkaufsdatumb+`</td>`)

if ( ausgabename == "#LAGER#" && ticketnr == "#LAGER#") {
	
		fmt.Fprintln(w, `<td valign="top" align="center">`+ausgabename+`</td>`)
		fmt.Fprintln(w, `<td valign="top" align="center"></td>`)
	
} else {

if ( ausgabename != "#LAGER#" && ticketnr == "#LAGER#") {
		fmt.Fprintln(w, `<td valign="top" align="center">`+ausgabename+`</td>`)
		fmt.Fprintln(w, `<td valign="top" align="center"></td>`)
}

if ( ausgabename == "#LAGER#" && ticketnr != "#LAGER#") {
		fmt.Fprintln(w, `<td valign="top" align="center"></td>`)
		fmt.Fprintln(w, `<td valign="top" align="center">`+ticketnr+`</td>`)
}


if ( ausgabename != "#LAGER#" && ticketnr != "#LAGER#") {
		fmt.Fprintln(w, `<td valign="top" align="center">`+ausgabename+`</td>`)
		fmt.Fprintln(w, `<td valign="top" align="center">`+ticketnr+`</td>`)
}

}		
		fmt.Fprintln(w, `
		<td valign="top" title="`+getuser(changed)+`" align="center">`+ausgabedatumb+`</td>
		<td valign="top" align="center" width="183">
  
		<div style="float: left; width: 60px;">
			<button onclick="changeitem('`+id+`','`+modell+`','`+seriennummer+`','`+zinfo+`','`+ausgabename+`','`+ausgabedatumb+`','`+changed+`','`+ticketnr+`')">ändern</button>
		</div>
		<div style="float: left; width: 60px; ">
			<form action="/user?gertyp=`+r.URL.Query().Get("gertyp")+`" method="post">
				<input type="hidden" name="getbackid" value="`+id+`">
				<input type="submit" name="additem" id="backitem" value="zurück">
			</form> 										
		</div>`)

							if checktokenstring[1] == "1" {
								fmt.Fprintln(w, `
		<div style="float: left; width: 60px; ">
			<form action="/user?gertyp=`+r.URL.Query().Get("gertyp")+`" method="post">
				<input type="hidden" name="delitemid" value="`+id+`">
				<input type="submit" name="additem" id="delitem" value="löschen" style="background-color: red">
			</form> 
		</div>`)
							}
							fmt.Fprintln(w, `  
</td>
</tr>`)

						}
						err = rows3.Err()
						if err != nil {
							log.Fatal(err)
						}
	
						} else {

					search := trimHtml(r.FormValue("search"))
					
					sqlq := `SELECT DISTINCT bestand.id, bestand.gertyp, bestand.modell, bestand.seriennummer, bestand.zinfo, bestand.ausgabename, bestand.ausgabedatum, bestand.changed, bestand.ticketnr, bestand.einkaufsdatum
						FROM bestand 
						JOIN zinfo ON bestand.gertyp
						JOIN zinfodata ON bestand.id
						where bestand.gertyp = `+mid+` and (zinfodata.daten like '%`+search+`%' or bestand.seriennummer like '%`+search+`%' or bestand.ausgabename like '%`+search+`%' or bestand.ausgabedatum like '%`+search+`%' or bestand.ticketnr like '%`+search+`%' or bestand.einkaufsdatum like '%`+search+`%') ;`

						rows3, err := db.Query(sqlq)
						
						if err != nil {
							log.Fatal(err)
						}

						for rows3.Next() {
							err := rows3.Scan(&id, &gertyp, &modell, &seriennummer, &zinfo, &ausgabename, &ausgabedatum, &changed, &ticketnr, &einkaufsdatum)
							if err != nil {
								log.Fatal(err)
							}

							fmt.Fprintln(w, `
	<tr>
		<td valign="top" align="center">`+getmodinfo(modell)+`</td>
		<td valign="top" align="center">`+seriennummer+`</td>

	`)

///zinfo table start

		//loadzinfo start
		for index := range zinfoarray {
		daten=""						
		err = db.QueryRow("SELECT daten FROM zinfodata where bestandid = \"" + id + "\" and zinfoid = \"" + zinfoarray[index] + "\" ").Scan(&daten)
		if err != nil {
		}
								
			if (daten !="") {
				fmt.Fprintln(w, `<td valign="top" align="center">`+ daten +`</td>`)
			} else {
				fmt.Fprintln(w, `<td valign="top" align="center"></td>`)
			}
		}
									
		//loadzinfo end


	split := strings.Split(einkaufsdatum, "-")
	einkaufsdatumb := ``+split[2]+`.`+split[1]+`.`+split[0]+``
	
	splitb := strings.Split(ausgabedatum, "-")
	ausgabedatumb := ``+splitb[2]+`.`+splitb[1]+`.`+splitb[0]+``

///zinfo table end
							fmt.Fprintln(w, `
		<td valign="top" align="center">`+einkaufsdatumb+`</td>`)

if ( ausgabename == "#LAGER#" && ticketnr == "#LAGER#") {
	
		fmt.Fprintln(w, `<td valign="top" align="center">`+ausgabename+`</td>`)
		fmt.Fprintln(w, `<td valign="top" align="center"></td>`)
	
} else {

if ( ausgabename != "#LAGER#" && ticketnr == "#LAGER#") {
		fmt.Fprintln(w, `<td valign="top" align="center">`+ausgabename+`</td>`)
		fmt.Fprintln(w, `<td valign="top" align="center"></td>`)
}

if ( ausgabename == "#LAGER#" && ticketnr != "#LAGER#") {
		fmt.Fprintln(w, `<td valign="top" align="center"></td>`)
		fmt.Fprintln(w, `<td valign="top" align="center">`+ticketnr+`</td>`)
}


if ( ausgabename != "#LAGER#" && ticketnr != "#LAGER#") {
		fmt.Fprintln(w, `<td valign="top" align="center">`+ausgabename+`</td>`)
		fmt.Fprintln(w, `<td valign="top" align="center">`+ticketnr+`</td>`)
}

}


		
		fmt.Fprintln(w, `
		<td valign="top" title="`+getuser(changed)+`" align="center">`+ausgabedatumb+`</td>
		<td valign="top" align="center" width="183">
  
		<div style="float: left; width: 60px;">
			<button onclick="changeitem('`+id+`','`+modell+`','`+seriennummer+`','`+zinfo+`','`+ausgabename+`','`+ausgabedatumb+`','`+changed+`','`+ticketnr+`')">ändern</button>
		</div>
		<div style="float: left; width: 60px; ">
			<form action="/user?gertyp=`+r.URL.Query().Get("gertyp")+`" method="post">
				<input type="hidden" name="getbackid" value="`+id+`">
				<input type="submit" name="additem" id="backitem" value="zurück">
			</form> 										
		</div>`)

							if checktokenstring[1] == "1" {
								fmt.Fprintln(w, `
		<div style="float: left; width: 60px; ">
			<form action="/user?gertyp=`+r.URL.Query().Get("gertyp")+`" method="post">
				<input type="hidden" name="delitemid" value="`+id+`">
				<input type="submit" name="additem" id="delitem" value="löschen" style="background-color: red">
			</form> 
		</div>`)
							}
							fmt.Fprintln(w, `  
</td>
</tr>`)

						}
						err = rows3.Err()
						if err != nil {
							log.Fatal(err)
						}
				
						}
										
						//search end output
										
					} else {
						


						
						rows3, err := db.Query("SELECT id, gertyp, modell, seriennummer, zinfo, ausgabename, ausgabedatum, changed, ticketnr, einkaufsdatum FROM bestand where gertyp = \"" + mid + "\" ")
						if err != nil {
							log.Fatal(err)
						}
					
						for rows3.Next() {
							err := rows3.Scan(&id, &gertyp, &modell, &seriennummer, &zinfo, &ausgabename, &ausgabedatum, &changed, &ticketnr, &einkaufsdatum)
							if err != nil {
								log.Fatal(err)
							}

							fmt.Fprintln(w, `
<tr>
<td valign="top" align="center">`+getmodinfo(modell)+`</td>
<td valign="top" align="center">`+seriennummer+`</td>
`)

					var changezinfo []string
					var vorhanden int = 0
					var sstring string
					//loadzinfo start
							for index := range zinfoarray {
								daten=""
								err = db.QueryRow("SELECT daten FROM zinfodata where bestandid = \"" + id + "\" and zinfoid = \"" + zinfoarray[index] + "\" ").Scan(&daten)
								if err != nil {
								}
								
								if (daten !="") {
									
									fmt.Fprintln(w, `<td valign="top" align="center">`+ daten +`</td>`)
									changezinfo = append(changezinfo, `,'`+daten+`'`)
									
									for indexb := range changeheader {
									
										sstring = `,zinfo`+zinfoarray[index]+``									
																
										if (sstring == changeheader[indexb]) {									
											vorhanden = 1										
										}
									
									}	
									
									if (vorhanden == 0 ){
										changeheader = append(changeheader, `,zinfo`+zinfoarray[index]+``)
										changeheaderb = append(changeheaderb, `zinfo`+zinfoarray[index]+``)
									} 
									vorhanden = 0
									
								} else {
								
								fmt.Fprintln(w, `<td valign="top" align="center"></td>`)
								
								}
							}
									
					//loadzinfo end

	split := strings.Split(einkaufsdatum, "-")
	einkaufsdatumb := ``+split[2]+`.`+split[1]+`.`+split[0]+``
	
	splitb := strings.Split(ausgabedatum, "-")
	ausgabedatumb := ``+splitb[2]+`.`+splitb[1]+`.`+splitb[0]+``

							fmt.Fprintln(w, `
<td valign="top" align="center">`+einkaufsdatumb+`</td>`)


if ( ausgabename == "#LAGER#" && ticketnr == "#LAGER#") {
	
		fmt.Fprintln(w, `<td valign="top" align="center">`+ausgabename+`</td>`)
		fmt.Fprintln(w, `<td valign="top" align="center"></td>`)
	
} else {

if ( ausgabename != "#LAGER#" && ticketnr == "#LAGER#") {
		fmt.Fprintln(w, `<td valign="top" align="center">`+ausgabename+`</td>`)
		fmt.Fprintln(w, `<td valign="top" align="center"></td>`)
}

if ( ausgabename == "#LAGER#" && ticketnr != "#LAGER#") {
		fmt.Fprintln(w, `<td valign="top" align="center"></td>`)
		fmt.Fprintln(w, `<td valign="top" align="center">`+ticketnr+`</td>`)
}


if ( ausgabename != "#LAGER#" && ticketnr != "#LAGER#") {
		fmt.Fprintln(w, `<td valign="top" align="center">`+ausgabename+`</td>`)
		fmt.Fprintln(w, `<td valign="top" align="center">`+ticketnr+`</td>`)
}

}



fmt.Fprintln(w, `
<td valign="top" title="`+getuser(changed)+`" align="center">`+ausgabedatumb+`</td>
<td valign="top" align="center" width="183">
  
<div style="float: left; width: 60px; ">
<button onclick="changeitem('`+id+`','`+modell+`','`+seriennummer+`','`+zinfo+`','`+ausgabename+`','`+ausgabedatumb+`','`+changed+`','`+ticketnr+`','`+einkaufsdatumb+`'`+strings.Join(changezinfo, "")+`)">ändern</button>
</div>
<div style="float: left; width: 60px; ">
<form action="/user?gertyp=`+r.URL.Query().Get("gertyp")+`" method="post">
<input type="hidden" name="getbackid" value="`+id+`">
<input type="submit" name="additem" id="backitem" value="zurück">
</form> 										
</div>`)

							if checktokenstring[1] == "1" {
								fmt.Fprintln(w, `		
<div style="float: left; width: 60px; ">
<form action="/user?gertyp=`+r.URL.Query().Get("gertyp")+`" method="post">
<input type="hidden" name="delitemid" value="`+id+`">
<input type="submit" name="additem" id="delitem" value="löschen" style="background-color: red">
</form> 
</div>`)
							}

							fmt.Fprintln(w, ` 
</td>
</tr>`)
						}

						err = rows3.Err()
						if err != nil {
							log.Fatal(err)
						}

					}
					///loop bestand end

					fmt.Fprintln(w, `
		</tbody>
		</table>
	</center>
	`)
				}

				///table bestand data end

				///load javascript start
				fmt.Fprintln(w, `
<script>
// HTML-Tabelle sortieren
const table = document.getElementById("foobar-table");
const headers = table.querySelectorAll("th");
const tableBody = table.querySelector("tbody");
const rows = tableBody.querySelectorAll("tr");
const directions = Array.from(headers).map(function(header) {
 return "";
});

const transform = function(index, content) {
 const type = headers[index].getAttribute("data-type");
 switch (type) {
     case "number":
  return parseFloat(content);
     case "string":
     default:
  return content;
 }
};

const sortColumn = function(index) {
 const direction = directions[index] || "asc";
  const multiplier = (direction === "asc") ? 1 : -1;
  const newRows = Array.from(rows);

  newRows.sort(function(rowA, rowB) {
     const cellA = rowA.querySelectorAll("td")[index].textContent;
     const cellB = rowB.querySelectorAll("td")[index].textContent;

     const a = transform(index, cellA);
     const b = transform(index, cellB);

     switch (true) {
      case a > b: return 1 * multiplier;
      case a < b: return -1 * multiplier;
      case a === b: return 0;
       }
   });

   [].forEach.call(rows, function(row) {
       tableBody.removeChild(row);
   });

   directions[index] = direction === "asc" ? "desc" : "asc";
   newRows.forEach(function(newRow) {
       tableBody.appendChild(newRow);
   });
 };

[].forEach.call(headers, function(header, index) {
   header.addEventListener("click", function() {
       sortColumn(index);
   });
});


function cancel(){
document.getElementById('changeitemid').value = "";
document.getElementById('sn').value = "";
document.getElementById('ticketnr').value = "";`)
for indexb := range changeheaderb {
fmt.Fprintln(w, `document.getElementById('`+changeheaderb[indexb]+`').value = "";`)																							
}
fmt.Fprintln(w, `
document.getElementById('ausgegebenan').value = "";
document.getElementById('ausgegebenam').value = "";
document.getElementById('edititem').disabled = !this.checked;
document.getElementById('sn').disabled = this.checked;
document.getElementById('einkaufsdatum').disabled = this.checked;
document.getElementById('modell').disabled = this.checked;
document.getElementById('additem').disabled = this.checked;
document.getElementById('ausgegebenan').disabled = !this.checked;
document.getElementById('ausgegebenam').disabled = !this.checked;
document.getElementById('ticketnr').disabled = !this.checked;
document.getElementById('addmen').open = false;
}



function changeitem(id,modellid,seriennummer,zinfo,ausgegebenan,ausgegebenam,changed,ticketnr,einkaufsdatum`+strings.Join(changeheader, "")+`){
document.getElementById('modell').namedItem('modellid' + modellid).selected=true;
document.getElementById('changeitemid').value = id;
document.getElementById('einkaufsdatum').value = einkaufsdatum;
document.getElementById('sn').value = seriennummer;`)
for indexb := range changeheaderb {
fmt.Fprintln(w, `document.getElementById('`+changeheaderb[indexb]+`').value = `+changeheaderb[indexb]+`;`)																							
}
fmt.Fprintln(w, `
document.getElementById('ticketnr').value = ticketnr;
document.getElementById('ausgegebenan').value = ausgegebenan;
document.getElementById('ausgegebenam').value = ausgegebenam;
document.getElementById('edititem').disabled = this.checked;
document.getElementById('additem').disabled = !this.checked;
document.getElementById('ausgegebenan').disabled = this.checked;
document.getElementById('ausgegebenam').disabled = this.checked;
document.getElementById('ticketnr').disabled = this.checked;
document.getElementById('addmen').open = true;
`)

				if checktokenstring[1] == "1" {
					fmt.Fprintln(w, `
document.getElementById('sn').disabled = this.checked;													
document.getElementById('zinfo').disabled = this.checked;							
document.getElementById('modell').disabled = this.checked;
document.getElementById('einkaufsdatum').disabled = this.checked;
`)
				} else {
					fmt.Fprintln(w, `
document.getElementById('sn').disabled = !this.checked;													
document.getElementById('zinfo').disabled = !this.checked;							
document.getElementById('modell').disabled = !this.checked;
document.getElementById('einkaufsdatum').disabled = this.checked;
`)
				}

fmt.Fprintln(w, `

//<![CDATA[
 $(document).ready(function(){
 
     // hide #back-top 
     $("#back-top2").hide();
    
     // fade in #back-top
     $(function () { 
         $(window).scroll(function () {
             if ($(this).scrollTop() > 200) {
                 $('#back-top2').fadeIn();
             } else {
                 $('#back-top2').fadeOut();
             }
         });
 
         // scroll body to 0px 
         $('#back-top2 a').click(function () {
             $('body,html').animate({
                 scrollTop: 0
             }, 500);
             return false;
         });
     });
 });
//]]>

`)





				fmt.Fprintln(w, `}
</script>
`)
				///load javascript end
			}

		}
		///content end

		if checktokenstring[1] == "1" {
			fmt.Fprintln(w, autologout+`
<script>
function Timer(s) {
 document.getElementById("timer").innerHTML=s;
 s--;
 if (s > 0) { // Sekunden (Endwert)
  window.setTimeout('Timer(' + s + ')', 999);
 }
}
Timer(180);
</script>
`)
		}

	}

	///stats start

	//stats1 start

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	if r.URL.Query().Get("gertyp") == "" {

		if r.URL.Query().Get("txt") != "1" {
			fmt.Fprintln(w, `<br>
<center>
<table border="0">
<tr>
<td>
<div class="glass">
<center>
<table border="0">
<tr>
<td>`)
		}

		var (
			id1     string
			gertyp1 string
			counta  int
			countb  int
			count   int
		)

		rows4, err := db.Query("SELECT id, gertyp FROM gertyp")
		if err != nil {
			log.Fatal(err)
		}

		if r.URL.Query().Get("txt") != "1" {
			fmt.Fprintln(w, `<br><center><table>`)
		}

		for rows4.Next() {
			err := rows4.Scan(&id1, &gertyp1)
			if err != nil {
				log.Fatal(err)
			}

			/////////

			err = db.QueryRow("SELECT COUNT(*) FROM bestand where gertyp = \"" + id1 + "\" ").Scan(&counta)
			if err != nil {
			}

			err = db.QueryRow("SELECT COUNT(*) FROM bestand where gertyp = \"" + id1 + "\" and ausgabename =\"#LAGER#\" ").Scan(&countb)
			if err != nil {
			}

			count = counta - countb

			countstr := strconv.Itoa(count)
			countbstr := strconv.Itoa(counta)

			if r.URL.Query().Get("txt") == "1" {
				fmt.Fprintln(w, ``+gertyp1+` ausgegeben: `+countstr+`/`+countbstr+``)
			} else {
				fmt.Fprintln(w, `<tr><td width ="300">`+gertyp1+` ausgegeben:</td><td>`+countstr+`/`+countbstr+`</td></tr>`)
			}

			/////////

		}
		err = rows4.Err()
		if err != nil {
			log.Fatal(err)
		}

		if r.URL.Query().Get("txt") != "1" {
			fmt.Fprintln(w, `</table></center>`)
		}

		//oldplaceboxstatistik

	}

	//stats1 end

	//stats2 start

	if r.URL.Query().Get("gertyp") == "" {

		var (
			id2          string
			modell2      string
			counta2      int
			count2       int
			countges     int
			sperrbestand int
			gertyp       string
		)

		rows5, err := db.Query("SELECT id, modell, sperrbestand FROM modell")
		if err != nil {
			log.Fatal(err)
		}

		if r.URL.Query().Get("txt") != "1" {


			fmt.Fprintln(w, `<br><center><table>`)
		}

		for rows5.Next() {
			err := rows5.Scan(&id2, &modell2, &sperrbestand)
			if err != nil {
				log.Fatal(err)
			}

			/////////

			err = db.QueryRow("SELECT COUNT(*) FROM bestand where modell = \"" + id2 + "\" ").Scan(&countges)
			if err != nil {
			}

			err = db.QueryRow("SELECT COUNT(*) FROM bestand where modell = \"" + id2 + "\" and ausgabename =\"#LAGER#\" ").Scan(&counta2)
			if err != nil {
			}

			count2 = countges - (countges - counta2)

			count2str := strconv.Itoa(count2)
			countgesstr := strconv.Itoa(countges)

			if (sperrbestand > 0 && countges > 0 && counta2 > 0) {
				if count2 <= sperrbestand {

					err = db.QueryRow("SELECT gertyp FROM bestand where modell = \"" + id2 + "\" and ausgabename =\"#LAGER#\" limit 1 ").Scan(&gertyp)
					if err != nil {
					}

					if gertyp != "" {
						if r.URL.Query().Get("txt") == "1" {
							fmt.Fprintln(w, ``+getgertypinfo(gertyp)+` `+modell2+` nachbestellen!!! `+count2str+` von `+countgesstr+` im Bestand.`)
						} else {
							fmt.Fprintln(w, `<tr><td width ="300"><span style="color:orange">`+getgertypinfo(gertyp)+` `+modell2+` nachbestellen!!!</span> </td><td>`+count2str+` von `+countgesstr+` im Bestand.</td></tr>`)
						}
					}

				}
			}

			/////////

		}
		err = rows5.Err()
		if err != nil {
			log.Fatal(err)
		}

		if r.URL.Query().Get("txt") != "1" {
			fmt.Fprintln(w, `</table></center>`)



		}

	}

	if r.URL.Query().Get("txt") != "1" {
		fmt.Fprintln(w, `
</td>
</tr>
</table>
</center>
</div>
</td>
</tr>
</table>
</center>
`)
	}
	//stats2 end
	///stats end

fmt.Fprintln(w, `
<div id="back-top2" style="display: block; ">
  <a href="#top">
<span class="obenkreis">
<i class="obenpfeil-1"></i>
<i class="obenpfeil-2"></i>
</span>
   </a>
</div>
`)

}

//////////////////admin site function//////////////////
func adminHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	deletetoken()

	checktokenstring := strings.Split(checktoken(w,r), "_")    

	if (checktokenstring[1] == "1" && checktokenstring[0] != "0" && checktokenstring[1] != "0") {

		db, err := sql.Open("mysql", sqlcred)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		var (
			ucheck		 string
			mcheck		 string
			gcheck		 string
			gertypcheck	 string
		)

		if r.Method == "POST" {
			r.ParseForm()


			if (r.FormValue("addfield") == "hinzufügen" && r.FormValue("field") != "" && r.FormValue("gertypid") != "") {
				err = db.QueryRow("SELECT gertyp FROM zinfo where zinfoname = \"" + r.FormValue("field") + "\" and zinfoname = \"" + r.FormValue("gertypid") + "\" ").Scan(&gertypcheck)
				if err != nil {
					//panic(err.Error())
				}

				if gertypcheck == "" {
					addfield(r.FormValue("gertypid"), r.FormValue("field"))
				}
			}

			if (r.FormValue("delfield") == "löschen" && r.FormValue("delfieldpid") != "") {

				var (
					delfieldcount int
				)

				err = db.QueryRow("SELECT COUNT(*) FROM zinfodata where zinfoid = \"" + r.FormValue("delfieldpid") + "\" ").Scan(&delfieldcount)
				if err != nil {
					//panic(err.Error())
				}


				if (delfieldcount > 0) {

					anzfield := strconv.Itoa(delfieldcount)

					fmt.Fprintln(w, `
					<script>alert("Fehler -> Dieses Feld hat noch `+anzfield+` Einträge!!!");</script>
				`)
				} else {
					delfield(r.FormValue("delfieldpid"))
				}








					

			}







			if r.FormValue("adduser") == "Benutzer hinzufügen" && r.FormValue("addusername") != "" && r.FormValue("addpassword") != "" {
				
				
				
				err = db.QueryRow("SELECT username FROM login where username = \"" + r.FormValue("addusername") + "\" ").Scan(&ucheck)
				if err != nil {
					//panic(err.Error())
				}

				if ucheck == "" {
					adduser(r.FormValue("addusername"), r.FormValue("addpassword"))
				}
			}

			if r.FormValue("deluser") == "löschen" && r.FormValue("deluserid") != "" {
				deluser(r.FormValue("deluserid"))
			}

			if r.FormValue("edituser") == "Speichern" && r.FormValue("edituserid") != "" {
			log.Printf(r.FormValue("edituserid"))
				edituser(r.FormValue("edituserid"), r.FormValue("addusername"), r.FormValue("addpassword"))
			}

			if r.FormValue("addmodell") == "Modell hinzufügen" && r.FormValue("addmodellname") != "" && r.FormValue("sperrbestand") != "" {
				err = db.QueryRow("SELECT modell FROM modell where modell = \"" + r.FormValue("addmodellname") + "\" ").Scan(&mcheck)
				if err != nil {
					//panic(err.Error())
				}

				if mcheck == "" {
					addmodell(r.FormValue("addmodellname"), r.FormValue("sperrbestand"))
				}
			}

			if r.FormValue("delmodell") == "löschen" && r.FormValue("delmodellid") != "" {

				var (
					delmodellcount int
				)

				err = db.QueryRow("SELECT COUNT(*) FROM bestand where modell = \"" + r.FormValue("delmodellid") + "\" ").Scan(&delmodellcount)
				if err != nil {
					//panic(err.Error())
				}

				if delmodellcount > 0 {

					anzmodell := strconv.Itoa(delmodellcount)

					fmt.Fprintln(w, `
					<script>alert("Fehler -> Modell hat noch `+anzmodell+` Einträge!!!");</script>
				`)
				} else {
					delmodell(r.FormValue("delmodellid"))
				}
			}

			if r.FormValue("addgertyp") == "Gerätetyp hinzufügen" && r.FormValue("addgertypname") != "" {
				err = db.QueryRow("SELECT gertyp FROM bestand where gertyp = \"" + r.FormValue("addgertypname") + "\" ").Scan(&gcheck)
				if err != nil {
					//panic(err.Error())
				}

				if gcheck == "" {
					addgertyp(r.FormValue("addgertypname"))
				}
			}

			if r.FormValue("delgertyp") == "löschen" && r.FormValue("delgertypid") != "" {

				var (
					delgertypcount int
				)

				err = db.QueryRow("SELECT COUNT(*) FROM zinfo where gertyp = \"" + r.FormValue("delgertypid") + "\" ").Scan(&delgertypcount)
				if err != nil {
					//panic(err.Error())
				}
				if delgertypcount > 0 {

					anzgertyp := strconv.Itoa(delgertypcount)

					fmt.Fprintln(w, `
					<script>alert("Fehler -> Gerätetyp hat noch `+anzgertyp+` Einträge!!!");</script>
				`)
				} else {
					delgertyp(r.FormValue("delgertypid"))
				}

			}

			if r.FormValue("changerechte") != "" && r.FormValue("changerechteuid") != "" && r.FormValue("changerechterid") != "" {
				changerechte(r.FormValue("changerechteuid"), r.FormValue("changerechterid"))
			}

			if r.FormValue("editmodell") == "Speichern" && r.FormValue("editmodellid") != "" && r.FormValue("addmodellname") != "" && r.FormValue("sperrbestand") != "" {
				editmodell(r.FormValue("editmodellid"), r.FormValue("addmodellname"), r.FormValue("sperrbestand"))
			}

			if r.FormValue("editgertyp") == "Speichern" && r.FormValue("editgertypid") != "" && r.FormValue("addgertypname") != "" {
				editgertyp(r.FormValue("editgertypid"), r.FormValue("addgertypname"))
			}
		}

		//page admin start

		fmt.Fprintln(w, globalmeta)
		fmt.Fprintln(w, style)

		fmt.Fprintln(w, `
		
		<script type='text/javascript'>
		//<![CDATA[
 		$(document).ready(function(){
 
 		    // hide #back-top 
 		    $("#back-top2").hide();
    
  		   // fade in #back-top
  		   $(function () { 
   		      $(window).scroll(function () {
     		        if ($(this).scrollTop() > 200) {
              		   $('#back-top2').fadeIn();
           		  } else {
           		      $('#back-top2').fadeOut();
          		   }
      		   });
 
     		    // scroll body to 0px 
     		    $('#back-top2 a').click(function () {
     		        $('body,html').animate({
      		           scrollTop: 0
      		       }, 500);
       		      return false;
      		   });
   		  });
 		});
		//]]>
		</script>
		
		`)

		//menu1 start

		fmt.Fprintln(w, menustart)

		if checktokenstring[1] == "1" {
			fmt.Fprintln(w, `<li><a href="/user">[Userbereich]</a></li>`)
		}

		fmt.Fprintln(w, `<li><a href="/user">[Statistik]</a></li>`)
		fmt.Fprintln(w, `<li><a href="/logout"><span style="color:red">[Logout]</span> <span id="timer" style="font-weight: bold;color:red"></span></a></li>`)

		fmt.Fprintln(w, menuend)

		//menu1 end

		//menu2 start
		fmt.Fprintln(w, menustart)

		//loop model menu start

		var (
			id           string
			gertyp       string
			gertypid     string
			username     string
			rechte       string
			aktiv        string
			modell       string
			sperrbestand string
			rechtemsg    string
			zinfoname    string
		)

		rows, err := db.Query("SELECT id, gertyp FROM gertyp")
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&id, &gertyp)
			if err != nil {
				panic(err.Error())
			}
			fmt.Fprintln(w, `<li><a href="/user?gertyp=`+id+`">`+gertyp+`</a></li>`)
		}
		err = rows.Err()
		if err != nil {
			panic(err.Error())
		}

		//loop model menu end

		fmt.Fprintln(w, menuend)
		//menu2 end

		//content start

		fmt.Fprintln(w, `
		<center>
		<script>
			function changeuser(id,username){
				document.getElementById('edituserid').value = id;
				document.getElementById('addusername').value = username;
			}
		</script>

		<table border="0" width="700" >
			<tr>
				<td>
					<form action="/admin#siteadd" method="post">	
					
		
					<table border="0" >
						<tr>
							<td >
								<table border="0" >
									<tr>
										<td width="70"></td>
										<td valign="top" align="center"><p id="siteadd">Benutzer hinzufügen</p></td>
									</tr>
									<tr>
										<td width="70">Username:</td>
										<td>
											<input type="hidden" name="edituserid" id="edituserid">
											<input type="text" name="addusername" id="addusername">
										</td>
									</tr>
									<tr>
										<td width="70">Passwort:</td>
										<td><input type="password" name="addpassword"></td>
									</tr>
									<tr>
										<td width="70"></td>
										<td valign="top" align="center"><input type="submit" name="adduser" value="Benutzer hinzufügen"><input type="submit" name="edituser" value="Speichern"></td>
									</tr>
								</table>
							</td>
						</tr>
					</table>
						
					</form>
				</td>
			</tr> 	
		</table>
		</center>
	`)

		//list users start
		fmt.Fprintln(w, `
		<center>
		<table border="0" width="700"><tr><td>
	`)

		rows2, err := db.Query("SELECT id, username, rechte, aktiv FROM login")
		if err != nil {
			panic(err.Error())
		}
		defer rows2.Close()

		for rows2.Next() {
			err := rows2.Scan(&id, &username, &rechte, &aktiv)
			if err != nil {
				panic(err.Error())
			}

			if rechte == "1" {
				rechtemsg = "Admin"
			} else {
				rechtemsg = "User"
			}

			fmt.Fprintln(w, `
			<table border="0" >
				<tr>
					<td width="200" valign="top">`+username+`</td><td width="50" valign="top"></td>
					<td width="195">
						<div style="float: left; width: 60px;">						
							<button onclick="changeuser('`+id+`','`+username+`')">ändern</button>							
						</div>
						<div style="float: left; width: 60px;">						
							<form action="/admin" method="post">
								<input type="hidden" name="changerechteuid" value="`+id+`">
								<input type="hidden" name="changerechterid" value="`+rechte+`">
								<input type="submit" name="changerechte" id="changerechte" value="`+rechtemsg+`">
							</form>
						</div>						
						<div style="float: left; width: 60px;">
							<form action="/admin" method="post">
								<input type="hidden" name="deluserid" value="`+id+`">
								<input type="submit" name="deluser" id="deluser" value="löschen" style="background-color: red">
							</form>
						</div>
					</td>
				</tr>
			</table>								
		`)

		}
		err = rows2.Err()
		if err != nil {
			panic(err.Error())
		}

		fmt.Fprintln(w, `
		</td></tr></table>
		</center>
	`)
		//list users end

		fmt.Fprintln(w, `<hr>`)

		//list modelle start
		fmt.Fprintln(w, `
		<script>
			function changemodell(id,modell,sperrbestand){
				document.getElementById('editmodellid').value = id;
				document.getElementById('addmodellname').value = modell;
				document.getElementById('sperrbestand').value = sperrbestand;
			}
		</script>

	<br><br>
	<center>
	<table border="0" width="700">
		<tr>
			<td>
				<form action="/admin#sitemodell" method="post">	

		
				<table border="0" >
					<tr>
						<td >
							<table border="0" >
								<tr>
									<td width="70"></td>
									<td valign="top" align="center"><p id="sitemodell">Modell hinzufügen</p></td>
								</tr>
								<tr>
									<td width="70">Modell:</td>
									<td>
										<input type="hidden" name="editmodellid" id="editmodellid">
										<input type="text" name="addmodellname" id="addmodellname">
									</td>
								</tr>
								<tr>
									<td width="70">Sperrbestand:</td>
									<td><input type="text" name="sperrbestand" id="sperrbestand" pattern="[0-9]+" value="0"></td>
								</tr>
								<tr>
									<td width="70"></td>
									<td valign="top" align="center"><input type="submit" name="addmodell" value="Modell hinzufügen"><input type="submit" name="editmodell" value="Speichern"></td>
								</tr>
							</table>
						</td>
					</tr>
				</table>	
				</form>
			</td>
		</tr> 	
	</table>
	</center>


	<center>
	<table border="0">
		<tr>
			<td>
	`)

		rows3, err := db.Query("SELECT id, modell, sperrbestand FROM modell")
		if err != nil {
			panic(err.Error())
		}
		defer rows3.Close()

		for rows3.Next() {
			err := rows3.Scan(&id, &modell, &sperrbestand)
			if err != nil {
				panic(err.Error())
			}

			fmt.Fprintln(w, `
			<table border="0" >
				<tr>
					<td width="200" valign="top">`+modell+`</td><td width="50" valign="top">`+sperrbestand+`</td>
					<td width="125">
						<div style="float: left; width: 60px;">						
							<button onclick="changemodell('`+id+`','`+modell+`','`+sperrbestand+`')">ändern</button>							
						</div>							
						<div style="float: right; width: 60px;">
							<form action="/admin" method="post">
								<input type="hidden" name="delmodellid" value="`+id+`">
								<input type="submit" name="delmodell" id="delmodell" value="löschen" style="background-color: red">
							</form>
						</div>
					</td>
				</tr>
			<table>								
		`)

		}
		err = rows3.Err()
		if err != nil {
			panic(err.Error())
		}

		fmt.Fprintln(w, `
			</td>
		</tr>
	</table>
	</center>
	`)

		//list modelle end

		fmt.Fprintln(w, `<hr>`)

		//list gertype start
		fmt.Fprintln(w, `
	<script>
		function changegertyp(id,gertyp){
			document.getElementById('editgertypid').value = id;
			document.getElementById('addgertypname').value = gertyp;
		}
	</script>

	<br><br>
	<center>
	<table border="0" width="700">
		<tr>
			<td>
				<form action="/admin#sitegertyp" method="post">			
					<table border="0" >
						<tr>
							<td >
								<table border="0" >
									<tr>
										<td width="70"></td>
										<td valign="top" align="center"><p id="sitegertyp">Gerätetyp hinzufügen</p></td>
									</tr>
									<tr>
										<td width="70">Gerätetyp:</td>
										<td>
											<input type="hidden" name="editgertypid" id="editgertypid">
											<input type="text" name="addgertypname" id="addgertypname">
										</td>
									</tr>
									<tr>
										<td width="70"></td>
										<td valign="top" align="center">
											<input type="submit" name="addgertyp" value="Gerätetyp hinzufügen"><input type="submit" name="editgertyp" value="Speichern">
										</td>
									</tr>
								</table>
							</td>
						</tr>
					</table>	
				</form>
			</td>
		</tr> 	
	</table>
	</center>


	<center>
	<table border="0">
		<tr>
			<td >
	`)

		rows4, err := db.Query("SELECT id, gertyp FROM gertyp")
		if err != nil {
			panic(err.Error())
		}
		defer rows4.Close()

		for rows4.Next() {
			err := rows4.Scan(&id, &gertyp)
			if err != nil {
				panic(err.Error())
			}

			fmt.Fprintln(w, `
			<table border="0" >
				<tr>
					<td width="200" valign="top">`+gertyp+`</td><td width="50" valign="top"></td>
					<td width="125">
						<div style="float: left; width: 60px;">						
							<button onclick="changegertyp('`+id+`','`+gertyp+`')">ändern</button>							
						</div>							
						<div style="float: right; width: 60px;">
							<form action="/admin" method="post">
								<input type="hidden" name="delgertypid" value="`+id+`">
								<input type="submit" name="delgertyp" id="delgertyp" value="löschen" style="background-color: red">
							</form>
						</div>
					</td>
				</tr>
			<table>								
		`)

		}
		err = rows4.Err()
		if err != nil {
			panic(err.Error())
		}

		fmt.Fprintln(w, `
			</td>
		</tr>
	</table>
	</center>
	`)

		//list gertype end
		
		//list zinfo start
		fmt.Fprintln(w, `<hr><br><br>`)
		fmt.Fprintln(w, `<center><p id="sitefield">Felder hinzufügen</p><table border="0" width="700"><tr>`)

		rows5, err := db.Query("SELECT id, gertyp FROM gertyp")
		if err != nil {
			panic(err.Error())
		}
		defer rows5.Close()

		for rows5.Next() {
			err := rows5.Scan(&gertypid, &gertyp)
		if err != nil {
			panic(err.Error())
		}

		fmt.Fprintln(w, `<tr><td>`+gertyp+`: </td><td>
		<table border="0"><tr>`)


		//dynamic fields start
		rows6, err := db.Query("SELECT id, zinfoname FROM zinfo where gertyp=\"" + gertypid + "\"")
		if err != nil {
			panic(err.Error())
		}
		defer rows6.Close()

		for rows6.Next() {
		err := rows6.Scan(&id, &zinfoname)
		if err != nil {
			panic(err.Error())
		}
		fmt.Fprintln(w, `
		<td><form action="/admin#sitefield" method="post">
			<table border="1">
				<tr><td align="center">`+zinfoname+`</td></tr>
				<tr><td align="center">
					<input type="hidden" name="delfieldpid" value="`+id+`">
					<input type="submit" name="delfield" id="delfield" value="löschen" style="background-color: red">
				</td></tr>
			</table></form>
		</td>`)
		}
		//dynamic fields end

		fmt.Fprintln(w, `
		</tr></table>
		</td><td>
			<form action="/admin#sitefield" method="post">
				<table border="0">
					<tr><td align="center"><input name="field" size="7"><input type="hidden" name="gertypid" value="`+gertypid+`" ></td></tr>
					<tr><td align="center"><input type="submit" name="addfield" value="hinzufügen"></td></tr>
				</table>
			</form>
		</td></tr>`)

		}

		fmt.Fprintln(w, `
		</tr></table>
		</center>`)	
		//list zinfo end
		
		
		//scroll up button start
		fmt.Fprintln(w, `
		<div id="back-top2" style="display: block; ">
 			<a href="#top">
				<span class="obenkreis">
					<i class="obenpfeil-1"></i>
					<i class="obenpfeil-2"></i>
				</span>
  			</a>
		</div>
		`)
		//scroll up button end		
		//content end
		//page admin end
	} else {
		fmt.Fprintln(w, "<meta http-equiv=\"refresh\" content=\"0; URL=/\">")
	}
	
	//logout timer
	fmt.Fprintln(w, autologout+`
		<script>
		function Timer(s) {
			document.getElementById("timer").innerHTML=s;
			s--;
			if (s > 0) { // Sekunden (Endwert)
				window.setTimeout('Timer(' + s + ')', 999);
			}
		}
		Timer(180);
		</script>
	`)
}

//////////////////login function//////////////////
func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		r.ParseForm()

		if r.FormValue("login") == "Login" && r.FormValue("password") != "" && r.FormValue("username") != "" {

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
			err = db.QueryRow("SELECT id, username, password, rechte, aktiv FROM login where username = \""+r.FormValue("username")+"\"").Scan(&id, &username, &password, &rechte1, &aktiv)
			if err != nil {
				w.Header().Set("Content-Type", "text/html")
				fmt.Fprintln(w, "<meta http-equiv=\"refresh\" content=\"0; URL=/\">")
			}

			if ((r.FormValue("username") == username) && (md5hash(r.FormValue("password")) == password)) {
				

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


//////////////////global content//////////////////
const style = `
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


.glass{
background: linear-gradient(135deg, rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0));
backdrop-filter: blur(10px);
-webkit-backdrop-filter: blur(10px);
border-radius: 10px;
border:1px solid rgba(255, 255, 255, 0.18);
box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.37);
width: 700px;
min-height: 130px;
}

body {
margin: 0;
padding: 0;
overflow: scroll;
background-image: url("./img/bg1.jpg"); //add "" if you want
background-repeat: no-repeat;
background-attachment: fixed;
background-size: cover;
color: black;
}



.navbar {
  background-color: lightblue;
  height: 50px;
}

nav ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

nav {
  float: right;
  margin: 0;
}

li {
  display: inline;
  padding: 15px;
}

nav a {
  text-decoration: none;
  color: #353535;
  display: inline-block;
  padding-top: 15px;
}

nav a:hover {
  color: #00b2b2;
}

nav a:active {
  font-weight: bold;
}

@media screen and (max-width: 700px) {
  nav {
    width: 100%;
    margin-top: -5px;
  }

  nav li {
    display: block;
    background-color: #e5e5e5;
    text-align: center;
  }
 }

* {
  font-family: Montserrat, sans-serif;
}


summary {
   position: relative;  
   }

summary::marker {
   color: transparent;
}

summary::after {
   content:  "+"; 
   position: absolute;
   color: black;
   font-size: 3em;
   font-weight: bold; 
   right: 1em;
   top: .2em;
   transition: all 0.5s;
} 

details[open] summary::after {
 color: red;
 transform: translate(5px,0) rotate(45deg);
}


table#foobar-table > thead > tr {
 background-color: #e5e5e5;
 color: black;
 cursor: Pointer;
 font-size: 18px;
}




.obenkreis {
  width: 40px;
  height: 40px;
  border-radius: 50px;
  background-color: rgba(51,51,51,0.6);		
  animation: anitop 1s;
}
@keyframes anitop {
  0%{opacity:0}
	100%{opacity:1}
  }
.obenpfeil-1, .obenpfeil-2 {
  border: solid #fff;						
  border-width: 0 3px 3px 0;
  display: inline-block;
  padding: 5px;
    transform: rotate(-135deg);
    -webkit-transform: rotate(-135deg);
	position:absolute;
   left:50%;
   margin-left:-6px;
	transition: all 0.2s ease;
	}
.obenpfeil-1 {
	top:15px;
	}
.obenpfeil-2 {
	top:22px;
	}
.obenkreis:hover .obenpfeil-2 {
	top:8px;
}
#back-top2 {
position: fixed;
bottom: 5%;
right:5%;
z-index: 1000;
}
#back-top2  span{
display: block;
}
@media (max-width: 1680px) {
#back-top2 {
bottom: 5px;
right:5px;
}}

</style>
`

const globalmeta = `
<html lang="de">
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta http-equiv="Cache-Control" content="no-cache, no-store, must-revalidate" />
<meta http-equiv="Pragma" content="no-cache" />
<meta http-equiv="Expires" content="0" />
`

const menustart = `
<div class="navbar">
  <nav>
     <ul>
`
const menuend = `
    </ul>
  </nav>
</div>
`

var autologout string = "<meta http-equiv=\"refresh\" content=\"180; URL=/logout\">"



//////////////////logout function//////////////////
func logoutHandler(w http.ResponseWriter, r *http.Request) {

	ip,_,_ := net.SplitHostPort(r.RemoteAddr)
	agent := r.Header.Get("User-Agent")
	t := time.Now()
	formatted := fmt.Sprintf("%02d.%02d.%d", t.Day(), t.Month(), t.Year())
	tokena := `l`+md5hash(formatted + " " + ip + " " + agent)+``

         c := http.Cookie{
                 Name:   tokena,
                 MaxAge: -1}
         http.SetCookie(w, &c)


	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<meta http-equiv=\"refresh\" content=\"0; URL=/\">")
	
}

//////////////////add item to bestand function//////////////////
func sqlinsert(gertyp string, modell string, sn string, zinfo string, einkaufsdatum string, userid string) {
	split := strings.Split(einkaufsdatum, ".")
	newdate := ``+split[2]+`-`+split[1]+`-`+split[0]+``
	
	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO bestand (gertyp, modell ,seriennummer, zinfo, ausgabename, ausgabedatum, changed, ticketnr,einkaufsdatum) VALUES (\"" + trimHtml(gertyp) + "\",\"" + trimHtml(modell) + "\",\"" + trimHtml(sn) + "\",\"" + trimHtml(zinfo) + "\",\"#LAGER#\",NOW(),\"" + userid + "\",\"#LAGER#\",\"" + trimHtml(newdate) + "\" )")

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

//////////////////zinfodatainput function//////////////////
func zinfodatainput(zinfoid string, bestandid string, daten string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO zinfodata (zinfoid, bestandid ,daten) VALUES (\"" + trimHtml(zinfoid) + "\",\"" + trimHtml(bestandid) + "\",\"" + trimHtml(daten) + "\")")

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}



//////////////////zinfodataedit function//////////////////
func zinfodataedit(zinfoid string, bestandid string, daten string) {
	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	edit, err := db.Query("UPDATE zinfodata SET daten=\"" + trimHtml(daten) + "\" where bestandid=\"" + bestandid + "\" and zinfoid=\"" + zinfoid + "\" ")
	if err != nil {
		panic(err.Error())
	}
	defer edit.Close()
}


//////////////////add user with default user permissions function//////////////////
func adduser(username string, password string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO login (username, password, rechte, aktiv) VALUES (\"" + trimHtml(username) + "\",\"" + md5hash(password) + "\",0,1)")

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

//////////////////add modell function//////////////////
func addmodell(modell string, sperrbestand string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO modell (modell, sperrbestand) VALUES (\"" + trimHtml(modell) + "\",\"" + trimHtml(sperrbestand) + "\")")

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}


//////////////////add field function//////////////////
func addfield(gertypid string, field string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO zinfo (gertyp, zinfoname) VALUES (\"" + trimHtml(gertypid) + "\",\"" + trimHtml(field) + "\")")

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

//////////////////add gertyp function//////////////////
func addgertyp(gertyp string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO gertyp (gertyp) VALUES (\"" + trimHtml(gertyp) + "\")")

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

//////////////////change modellname and sperrbestand by id function//////////////////
func editmodell(editmodellid string, addmodellname string, sperrbestand string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	change, err := db.Query("UPDATE modell SET modell=\"" + trimHtml(addmodellname) + "\" , sperrbestand=\"" + trimHtml(sperrbestand) + "\" where id=\"" + editmodellid + "\"")

	if err != nil {
		panic(err.Error())
	}

	defer change.Close()
}

//////////////////change gertypname by id function//////////////////
func editgertyp(editgertypid string, addgertypname string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	change, err := db.Query("UPDATE gertyp SET gertyp=\"" + trimHtml(addgertypname) + "\" where id=\"" + editgertypid + "\"")

	if err != nil {
		panic(err.Error())
	}

	defer change.Close()
}

//////////////////with user permissions change values from bestand by id function//////////////////
func sqledit2(changeitemid string, gertyp string, ausgegebenan string, ausgegebenam string, ticketnr string, userid string) {

	split := strings.Split(ausgegebenam, ".")
	ausgegebenamb := ``+split[2]+`-`+split[1]+`-`+split[0]+``

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	change, err := db.Query("UPDATE bestand SET gertyp=\"" + trimHtml(gertyp) + "\", ausgabename=\"" + trimHtml(ausgegebenan) + "\" , ausgabedatum=\"" + trimHtml(ausgegebenamb) + "\", changed=\"" + userid + "\", ticketnr=\"" + ticketnr + "\" where id=\"" + changeitemid + "\"")

	if err != nil {
		panic(err.Error())
	}

	defer change.Close()
}

//////////////////with admin permissions change values from bestand by id function//////////////////
func sqledit(modell string, changeitemid string, gertyp string, sn string, ausgegebenan string, ausgegebenam string, ticketnr string, einkaufsdatum string, userid string) {

	split := strings.Split(einkaufsdatum, ".")
	einkaufsdatumb := ``+split[2]+`-`+split[1]+`-`+split[0]+``

	splitb := strings.Split(ausgegebenam, ".")
	ausgegebenamb := ``+splitb[2]+`-`+splitb[1]+`-`+splitb[0]+``

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	change, err := db.Query("UPDATE bestand SET gertyp=\"" + trimHtml(gertyp) + "\", modell=\"" + trimHtml(modell) + "\", seriennummer=\"" + trimHtml(sn) + "\" , ausgabename=\"" + trimHtml(ausgegebenan) + "\" , ausgabedatum=\"" + trimHtml(ausgegebenamb) + "\", changed=\"" + userid + "\", ticketnr=\"" + ticketnr + "\", einkaufsdatum=\"" + einkaufsdatumb + "\" where id=\"" + changeitemid + "\" ")

	if err != nil {
		panic(err.Error())
	}
	defer change.Close()

}

//////////////////change username or password function//////////////////
func edituser(edituserid string, addusername string, addpassword string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	change, err := db.Query("UPDATE login SET username=\"" + trimHtml(addusername) + "\" where id=\"" + edituserid + "\"")
	if err != nil {
		panic(err.Error())
	}

	defer change.Close()

	if addpassword != "" {
		change, err := db.Query("UPDATE login SET password=\"" + md5hash(addpassword) + "\" where id=\"" + edituserid + "\"")

		if err != nil {
			panic(err.Error())
		}

		defer change.Close()
	}
}


//////////////////update login token function//////////////////
func updatetoken(token string, userid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	change, err := db.Query("UPDATE login SET session=\"" + token + "\", sessiontime=NOW() where id=\"" + userid + "\"")
	if err != nil {
		panic(err.Error())
	}

	defer change.Close()
}


//////////////////update login token function//////////////////
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

//////////////////change user permissions function//////////////////
func changerechte(changerechteuid string, changerechterid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var (
		changemsg string
	)

	if changerechterid == "1" {
		changemsg = "0"
	} else {
		changemsg = "1"
	}

	change, err := db.Query("UPDATE login SET rechte=\"" + changemsg + "\" where id=\"" + changerechteuid + "\"")
	if err != nil {
		panic(err.Error())
	}

	defer change.Close()

}

//////////////////delete item by id function//////////////////
func delitem(delitemid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	getback, err := db.Query("DELETE from bestand where id=\"" + delitemid + "\"")
	if err != nil {
		panic(err.Error())
	}
	defer getback.Close()
	
	getbackb, err := db.Query("DELETE from zinfodata where bestandid=\"" + delitemid + "\"")
	if err != nil {
		panic(err.Error())
	}
	defer getbackb.Close()		
}

//////////////////delete user by id function//////////////////
func deluser(deluserid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	getback, err := db.Query("DELETE from login where id=\"" + deluserid + "\"")

	if err != nil {
		panic(err.Error())
	}

	defer getback.Close()
}

//////////////////delete modell by id function//////////////////
func delmodell(delmodellid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	getback, err := db.Query("DELETE from modell where id=\"" + delmodellid + "\"")
	if err != nil {
		panic(err.Error())
	}
	defer getback.Close()
}

//////////////////delete gertyp by id function//////////////////
func delgertyp(delgertypid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	getback, err := db.Query("DELETE from gertyp where id=\"" + delgertypid + "\"")
	if err != nil {
		panic(err.Error())
	}
	defer getback.Close()
}

//////////////////delete delfield by id function//////////////////
func delfield(delfieldpid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	getback, err := db.Query("DELETE from zinfo where id=\"" + delfieldpid + "\"")
	if err != nil {
		panic(err.Error())
	}
	defer getback.Close()
}

//////////////////get model name by id function//////////////////
func getmodinfo(modell string) string {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var (
		modellb string
	)

	err = db.QueryRow("SELECT modell FROM modell where id = \"" + modell + "\"").Scan(&modellb)
	if err != nil {
		panic(err.Error())
	}
	return modellb
}

//////////////////get item back by id function//////////////////
func getback(getbackid string, userid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	getback, err := db.Query("UPDATE bestand SET ticketnr='#LAGER#', ausgabename='#LAGER#', ausgabedatum=NOW(), changed=\"" + userid + "\" where id=\"" + getbackid + "\"")
	if err != nil {
		panic(err.Error())
	}
	defer getback.Close()
}

//////////////////get gertypname by id function//////////////////
func getgertypinfo(gertyp string) string {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var (
		gertypname string
	)

	err = db.QueryRow("SELECT gertyp FROM gertyp where id = \"" + gertyp + "\"").Scan(&gertypname)
	if err != nil {
		panic(err.Error())
	}
	return gertypname
}


//////////////////get user by id function//////////////////
func getuser(getuserid string) string {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.QueryRow("SELECT username FROM login where id = \"" + getuserid + "\"").Scan(&getuserid)
	if err != nil {
		panic(err.Error())
	}
	return getuserid
}


//////////////////gen md5 function//////////////////
func md5hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//////////////////random string function//////////////////
func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

//////////////////remove bad chars function//////////////////
func trimHtml(src string) string {
	//Convert all HTML tags to lowercase
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//Remove STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//Remove SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//Remove all HTML code in angle brackets and replace them with newline characters
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//Remove consecutive newlines
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	//remove all other bads
	re, _ = regexp.Compile("[''(){}<> /%[\\]]")
	src = re.ReplaceAllString(src, "")
	return strings.TrimSpace(src)
	
}
