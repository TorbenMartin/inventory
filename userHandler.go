package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)


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

