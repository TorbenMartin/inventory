package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
	"net"

)

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
