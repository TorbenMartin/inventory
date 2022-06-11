package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"html"
	"net"
	"time"
)

//////////////////admin site function//////////////////
func adminHandler(w http.ResponseWriter, r *http.Request) {

	secheader(w)
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
				err = db.QueryRow("SELECT gertyp FROM zinfo where zinfoname = ? and zinfoname = ? ", r.FormValue("field"), r.FormValue("gertypid")).Scan(&gertypcheck)
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

				err = db.QueryRow("SELECT COUNT(*) FROM zinfodata where zinfoid = ? ", r.FormValue("delfieldpid")).Scan(&delfieldcount)
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
				
				
				
				err = db.QueryRow("SELECT username FROM login where username = ? ", r.FormValue("addusername")).Scan(&ucheck)
				if err != nil {
					//panic(err.Error())
				}

				if ucheck == "" {
					adduser(r.FormValue("addusername"), r.FormValue("addpassword"))
				}
			}

			if r.FormValue("deluser") == "löschen" && r.FormValue("deluserid") != "" {

				ip,_,_ := net.SplitHostPort(r.RemoteAddr)
				agent := r.Header.Get("User-Agent")
				t := time.Now()
				formatted := fmt.Sprintf("%02d.%02d.%d", t.Day(), t.Month(), t.Year())
				tokena := `l`+md5hash(formatted + " " + ip + " " + agent)+``

		      		c, err := r.Cookie(tokena)
				if err != nil {
					//panic(err.Error())
				} else {

				value := c.Value
				splitsting := strings.Split(value, "_")  

				if r.FormValue("deluserid") != splitsting[0] {					
					deluser(r.FormValue("deluserid"))	
				}

				}
				
			}

			if r.FormValue("edituser") == "Speichern" && r.FormValue("edituserid") != "" {
				edituser(r.FormValue("edituserid"), r.FormValue("addusername"), r.FormValue("addpassword"))
			}

			if r.FormValue("addmodell") == "Modell hinzufügen" && r.FormValue("gertypaddmodell") != "" && r.FormValue("addmodellname") != "" && r.FormValue("sperrbestand") != "" {
				err = db.QueryRow("SELECT modell FROM modell where modell = ? ", r.FormValue("addmodellname")).Scan(&mcheck)
				if err != nil {
					//panic(err.Error()) 
				}

				if mcheck == "" {
					addmodell(r.FormValue("addmodellname"), r.FormValue("gertypaddmodell"), r.FormValue("sperrbestand"))
				}
			}

			if r.FormValue("delmodell") == "löschen" && r.FormValue("delmodellid") != "" {

				var (
					delmodellcount int
				)

				err = db.QueryRow("SELECT COUNT(*) FROM bestand where modell = ? ", r.FormValue("delmodellid")).Scan(&delmodellcount)
				if err != nil {
					//panic(err.Error())
				}

				if delmodellcount > 0 {

					anzmodell := strconv.Itoa(delmodellcount)

					fmt.Fprintln(w, `
					<script>alert("Fehler -> Modell hat noch `+anzmodell+` Einträge im Bestand!!!");</script>
				`)
				} else {
					delmodell(r.FormValue("delmodellid"))
				}
			}

			if r.FormValue("addgertyp") == "Gerätetyp hinzufügen" && r.FormValue("addgertypname") != "" {
				err = db.QueryRow("SELECT gertyp FROM bestand where gertyp = ? ", r.FormValue("addgertypname")).Scan(&gcheck)
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

				err = db.QueryRow("SELECT COUNT(*) FROM modell where gertyp = ? ", r.FormValue("delgertypid")).Scan(&delgertypcount)
	
				if err != nil {
					//panic(err.Error())
				}
				if delgertypcount > 0 {

					anzgertyp := strconv.Itoa(delgertypcount)

					fmt.Fprintln(w, `
					<script>alert("Fehler -> Gerätetyp hat noch `+anzgertyp+` Einträge bei Modellen!!!");</script>
				`)
				} else {
				
				err = db.QueryRow("SELECT COUNT(*) FROM zinfo where gertyp = ? ", r.FormValue("delgertypid")).Scan(&delgertypcount)
	
				if err != nil {
					//panic(err.Error())
				}
				if delgertypcount > 0 {

					anzgertyp := strconv.Itoa(delgertypcount)

					fmt.Fprintln(w, `
					<script>alert("Fehler -> Gerätetyp hat noch `+anzgertyp+` Einträge \nbei den Zusatz-Infos!!!");</script>
				`)
				} else {
				
					delgertyp(r.FormValue("delgertypid"))

				}

				}

			}

			if r.FormValue("changerechte") != "" && r.FormValue("changerechteuid") != "" && r.FormValue("changerechterid") != "" {

				ip,_,_ := net.SplitHostPort(r.RemoteAddr)
				agent := r.Header.Get("User-Agent")
				t := time.Now()
				formatted := fmt.Sprintf("%02d.%02d.%d", t.Day(), t.Month(), t.Year())
				tokena := `l`+md5hash(formatted + " " + ip + " " + agent)+``

		      		c, err := r.Cookie(tokena)
				if err != nil {
					//panic(err.Error())
				} else {

				value := c.Value
				splitsting := strings.Split(value, "_")  

				if r.FormValue("changerechteuid") != splitsting[0] {					
					changerechte(r.FormValue("changerechteuid"), r.FormValue("changerechterid"))
				}

				}



			}

			if r.FormValue("editmodell") == "Speichern" && r.FormValue("editmodellid") != "" && r.FormValue("addmodellname") != "" && r.FormValue("sperrbestand") != "" && r.FormValue("gertypaddmodell") != "" {
				editmodell(r.FormValue("editmodellid"), r.FormValue("addmodellname"), r.FormValue("sperrbestand"), r.FormValue("gertypaddmodell"))
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

		fmt.Fprintln(w, `<br><button onclick="window.location.href = './user';">Statistik</button>`)
		fmt.Fprintln(w, `<button style="background-color: red" onclick="window.location.href = './logout';">Logout</button>&nbsp;&nbsp;`)
//<span id="timer" style="font-weight: bold;color:red"></span>
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
			fmt.Fprintln(w, `<button onclick="window.location.href = './user?gertyp=`+id+`';">`+html.EscapeString(gertyp)+`</button>`)
		}
		err = rows.Err()
		if err != nil {
			panic(err.Error())
		}

		//loop model menu end

		fmt.Fprintln(w, `&nbsp;`+menuend)
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

		<table border="0" width="700" id="siteuser">
			<tr>
				<td>
					<form action="/admin#siteuser" method="post">	
					
		
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
					<td width="200" valign="top">`+html.EscapeString(username)+`</td><td width="50" valign="top"></td>
					<td width="195" align="right">
						<div style="float: left; width: 60px;">						
							<button onclick="changeuser('`+id+`','`+username+`')">ändern</button>							
						</div>
						<div style="float: left; width: 60px;">						
							<form action="/admin#siteuser" method="post">
								<input type="hidden" name="changerechteuid" value="`+id+`">
								<input type="hidden" name="changerechterid" value="`+rechte+`">
								<input type="submit" name="changerechte" id="changerechte" value="`+rechtemsg+`">
							</form>
						</div>						
						<div style="float: left; width: 60px;">
							<form action="/admin#siteuser" method="post" onSubmit="return confirm('Soll der Eintrag gelöscht werden?');">
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
	<table border="0" width="700" id="sitegertyp">
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
	<table border="0" width="700">
		<tr>
			<td>

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

			<table>
				<tr>
					<td width="315">`+html.EscapeString(gertyp)+`</td>
					<td>
						<div style="float: left; width: 60px;">						
							<button onclick="changegertyp('`+id+`','`+gertyp+`')">ändern</button>							
						</div>	
						<div style="float: right; width: 60px;">
							<form action="/admin" method="post" onSubmit="return confirm('Soll der Eintrag gelöscht werden?');">
								<input type="hidden" name="delgertypid" value="`+id+`">
								<input type="submit" name="delgertyp" id="delgertyp" value="löschen" style="background-color: red">
							</form>
						</div>
					</td>
				</tr>
			</table>
		
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

		//list modelle start


		var (
			gertypcount int
		)

		err = db.QueryRow("SELECT COUNT(*) FROM gertyp").Scan(&gertypcount)
		if err != nil {
			//panic(err.Error())
		}

		if (gertypcount > 0) {

		fmt.Fprintln(w, `<hr>`)

		fmt.Fprintln(w, `
		<script>
			function changemodell(id,modell,sperrbestand, gertyp){
				document.getElementById('editmodellid').value = id;
				document.getElementById('addmodellname').value = modell;
				document.getElementById('sperrbestand').value = sperrbestand;
				document.getElementById('gertypaddmodell').value = gertyp;				
			}
		</script>

	<br><br>
	<center>
	<table border="0" width="700" id="sitemodell">
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
									<td width="70">Gerätetyp:</td>
									<td>
										<select name="gertypaddmodell" id="gertypaddmodell">
		`)
		
		rows7, err := db.Query("SELECT id, gertyp FROM gertyp")
		if err != nil {
			panic(err.Error())
		}
		defer rows7.Close()

		for rows7.Next() {
			err := rows7.Scan(&id, &gertyp)
			if err != nil {
				panic(err.Error())
			}
		fmt.Fprintln(w, `												
			<option value="`+id+`">`+gertyp+`</option>
		`)
		}	


		fmt.Fprintln(w, `												
										</select>

									
									</td>
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
	<table border="0" width="700">
		<tr>
			<td>
		`)

		rows3, err := db.Query("SELECT id, modell, sperrbestand, gertyp FROM modell")
		if err != nil {
			panic(err.Error())
		}
		defer rows3.Close()

		for rows3.Next() {
			err := rows3.Scan(&id, &modell, &sperrbestand, &gertyp)
			if err != nil {
				panic(err.Error())
			}

		fmt.Fprintln(w, `
		
		<table>
			<tr>
				<td width="120">`+html.EscapeString(getgertyp(gertyp))+`</td>
				<td width="100">`+html.EscapeString(modell)+`</td>
				<td width="95">`+html.EscapeString(sperrbestand)+`</td>
				<td>
				
						<div style="float: left; width: 60px;">						
							<button onclick="changemodell('`+id+`','`+modell+`','`+sperrbestand+`','`+gertyp+`')">ändern</button>							
						</div>							
						<div style="float: right; width: 60px;">
							<form action="/admin#sitemodell" method="post" onSubmit="return confirm('Soll der Eintrag gelöscht werden?');">
								<input type="hidden" name="delmodellid" value="`+id+`">
								<input type="submit" name="delmodell" id="delmodell" value="löschen" style="background-color: red">
							</form>
						</div>
				
				</td>
			</tr>
		</table>
		
		`)

		 }

		fmt.Fprintln(w, `
		</td>
		</tr>
	</table>
	</center>

		`)


		}

		//list modelle end

		fmt.Fprintln(w, `<hr>`)
		
		//list zinfo start
		fmt.Fprintln(w, `<br><br>`)

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

		fmt.Fprintln(w, `
		<center>
		<table border="0" width="700" id="sitefield">
		<tr>
			<td width="120">Gerätetyp: </td>
			<td>`+html.EscapeString(gertyp)+`</td>
			<form action="/admin#sitefield" method="post">
			<td align="right"><input name="field" size="14"><input type="hidden" name="gertypid" value="`+gertypid+`" ></td>
			<td width="235" align="center"><input type="submit" name="addfield" value="hinzufügen"></td>
			</form>
			
		</tr>
		`)


		//dynamic fields start
		rows6, err := db.Query("SELECT id, zinfoname FROM zinfo where gertyp= ? ", gertypid)
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
		<tr>
			<form action="/admin#sitefield" method="post" onSubmit="return confirm('Soll der Eintrag gelöscht werden?');">
			
			<td width="120">Zusatzfeld:</td>
			<td>`+html.EscapeString(zinfoname)+`</td>
			<td align="right">
			
				<input type="hidden" name="delfieldpid" value="`+id+`">
				<input type="submit" name="delfield" id="delfield" value="löschen" style="background-color: red">
		
			</td>

			</form>		
		</tr>
		`)

		}

		//dynamic fields end

	
		fmt.Fprintln(w, `		
		</table>
		</center><br><br>
		`)


		}


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

	fmt.Fprintln(w, autologout)
	/*
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
	`)*/
}
