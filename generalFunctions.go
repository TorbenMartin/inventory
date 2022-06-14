package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
	"io/ioutil"
        "encoding/json"
	"net/http"
//	"fmt"
)

//////////////////add item to bestand function//////////////////
func sqlinsert (gertyp string, modell string, sn string, zinfo string, einkaufsdatum string, userid string) {
	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	split := strings.Split(einkaufsdatum, ".")
	newdate := ``+split[2]+`-`+split[1]+`-`+split[0]+``
	
	insert, err1 := db.Prepare("INSERT INTO bestand (gertyp, modell ,seriennummer, zinfo, ausgabename, ausgabedatum, changed, ticketnr,einkaufsdatum) VALUES (?,?,?,?,\"#LAGER#\",NOW(),?,\"#LAGER#\",?) ")
	if err1 != nil {
		panic(err1.Error())
	}
	defer insert.Close()
	
	if _, err2 := insert.Exec(gertyp, modell, sn, zinfo, userid, newdate); err2 != nil {
		panic(err2.Error())
	}	
	
}
//////////////////zinfodatainput function//////////////////
func zinfodatainput(zinfoid string, bestandid string, daten string) {
	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	insert, err := db.Prepare("INSERT INTO zinfodata (zinfoid, bestandid ,daten) VALUES (?,?,?) ")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	if _, err := insert.Exec(zinfoid,bestandid,daten); err != nil {
		panic(err.Error())
	}
}
//////////////////zinfodataedit function//////////////////
func zinfodataedit(zinfoid string, bestandid string, daten string) {
	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	edit, err := db.Prepare("UPDATE zinfodata SET daten= ? where bestandid= ? and zinfoid= ? ")
	if err != nil {
		panic(err.Error())
	}
	defer edit.Close()

	if _, err := edit.Exec(daten,bestandid,zinfoid); err != nil {
		panic(err.Error())
	}
}
//////////////////add user with default user permissions function//////////////////
func adduser(username string, password string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	insert, err := db.Prepare("INSERT INTO login (username, password, rechte, aktiv) VALUES (?,?,0,1) ")

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	if _, err := insert.Exec(username,md5hash(``+password+``+keysalt+``)); err != nil {
		panic(err.Error())
	}
}
//////////////////add modell function//////////////////
func addmodell(modell string, gertypaddmodell string, sperrbestand string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()
	
	insert, err := db.Prepare("INSERT INTO modell (modell, gertyp, sperrbestand) VALUES (?,?,?) ")

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	if _, err := insert.Exec(modell,gertypaddmodell,sperrbestand); err != nil {
		panic(err.Error())
	}

}
//////////////////add field function//////////////////
func addfield(gertypid string, field string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	insert, err := db.Prepare("INSERT INTO zinfo (gertyp, zinfoname) VALUES (?,?) ")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	if _, err := insert.Exec(gertypid,field); err != nil {
		panic(err.Error())
	}
}

//////////////////add gertyp function//////////////////
func addgertyp(gertyp string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	insert, err := db.Prepare("INSERT INTO gertyp (gertyp) VALUES (?) ")

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	if _, err := insert.Exec(gertyp); err != nil {
		panic(err.Error())
	}
}
//////////////////change modellname and sperrbestand by id function//////////////////
func editmodell(editmodellid string, addmodellname string, sperrbestand string, gertypnew string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	var (
		gertypold string
	)
	err1 := db.QueryRow("SELECT gertyp FROM modell where id = ? ", editmodellid).Scan(&gertypold)
	if err1 != nil {
		panic(err1.Error())
	}
	
	change, err := db.Prepare("UPDATE modell SET modell = ? , sperrbestand = ? , gertyp = ? where id = ? ")
	if err != nil {
		panic(err.Error())
	}
	defer change.Close()
	if _, err := change.Exec(addmodellname,sperrbestand,gertypnew,editmodellid); err != nil {
		panic(err.Error())
	}

		var (
			bestandcount int
		)

		err2 := db.QueryRow("SELECT COUNT(*) FROM bestand where gertyp = ? ", gertypold).Scan(&bestandcount)
		if err2 != nil {
			panic(err2.Error())
		}
		
		
		if (bestandcount > 0) {
		

			changeb, err := db.Prepare("UPDATE bestand SET gertyp = ? where gertyp = ? ")
			if err != nil {
				panic(err.Error())
			}
			
			defer changeb.Close()
			

			if _, err := changeb.Exec(gertypnew,gertypold); err != nil {
				panic(err.Error())
			}

		}

}
//////////////////change gertypname by id function//////////////////
func editgertyp(editgertypid string, addgertypname string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	change, err := db.Prepare("UPDATE gertyp SET gertyp= ? where id= ? ")

	if err != nil {
		panic(err.Error())
	}
	defer change.Close()

	if _, err := change.Exec(addgertypname,editgertypid); err != nil {
		panic(err.Error())
	}
}
//////////////////with user permissions change values from bestand by id function//////////////////
func sqledit2(changeitemid string, gertyp string, ausgegebenan string, ausgegebenam string, ticketnr string, userid string) {

	split := strings.Split(ausgegebenam, ".")
	ausgegebenamb := ``+split[2]+`-`+split[1]+`-`+split[0]+``

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	change, err := db.Prepare("UPDATE bestand SET gertyp= ?, ausgabename= ? , ausgabedatum= ?, changed= ?, ticketnr= ? where id= ? ")

	if err != nil {
		panic(err.Error())
	}
	defer change.Close()

	if _, err := change.Exec(gertyp,ausgegebenan,ausgegebenamb,userid,ticketnr,changeitemid); err != nil {
		panic(err.Error())
	}

}
//////////////////with admin permissions change values from bestand by id function//////////////////
func sqledit(modell string, changeitemid string, gertyp string, sn string, ausgegebenan string, ausgegebenam string, ticketnr string, einkaufsdatum string, userid string) {

	split := strings.Split(einkaufsdatum, ".")
	einkaufsdatumb := ``+split[2]+`-`+split[1]+`-`+split[0]+``

	splitb := strings.Split(ausgegebenam, ".")
	ausgegebenamb := ``+splitb[2]+`-`+splitb[1]+`-`+splitb[0]+``

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()
	
	change, err := db.Prepare("UPDATE bestand SET gertyp= ?, modell= ?, seriennummer= ?, ausgabename= ? , ausgabedatum= ?, changed= ?, ticketnr= ?, einkaufsdatum= ? where id= ? ")

	if err != nil {
		panic(err.Error())
	}
	defer change.Close()

	if _, err := change.Exec(gertyp,modell,sn,ausgegebenan,ausgegebenamb,userid,ticketnr,einkaufsdatumb,changeitemid); err != nil {
		panic(err.Error())
	}
}
//////////////////change username or password function//////////////////
func edituser(edituserid string, addusername string, addpassword string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	change, err := db.Prepare("UPDATE login SET username= ? where id= ? ")
	if err != nil {
		panic(err.Error())
	}

	defer change.Close()

	if _, err := change.Exec(addusername,edituserid); err != nil {
		panic(err.Error())
	}

	if addpassword != "" {
		change, err := db.Prepare("UPDATE login SET password= ? where id= ? ")

		if err != nil {
			panic(err.Error())
		}

		defer change.Close()

		if _, err := change.Exec(md5hash(``+addpassword+``+keysalt+``),edituserid); err != nil {
			panic(err.Error())
		}

	}
}
//////////////////change user permissions function//////////////////
func changerechte(changerechteuid string, changerechterid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
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

	change, err := db.Prepare("UPDATE login SET rechte= ? where id= ? ")
	if err != nil {
		panic(err.Error())
	}

	defer change.Close()

	if _, err := change.Exec(changemsg,changerechteuid); err != nil {
		panic(err.Error())
	}
}
//////////////////delete item by id function//////////////////
func delitem(delitemid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()
	
	getback, err := db.Prepare("DELETE from bestand where id= ? ")
	if err != nil {
		panic(err.Error())
	}
	defer getback.Close()

	if _, err := getback.Exec(delitemid); err != nil {
		panic(err.Error())
	}
	
	getbackb, err := db.Prepare("DELETE from zinfodata where bestandid= ? ")
	if err != nil {
		panic(err.Error())
	}
	defer getbackb.Close()		

	if _, err := getbackb.Exec(delitemid); err != nil {
		panic(err.Error())
	}
	
}
//////////////////delete user by id function//////////////////
func deluser(deluserid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	getback, err := db.Prepare("DELETE from login where id= ? ")
	if err != nil {
		panic(err.Error())
	}

	defer getback.Close()

	if _, err := getback.Exec(deluserid); err != nil {
		panic(err.Error())
	}
}
//////////////////delete modell by id function//////////////////
func delmodell(delmodellid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	getback, err := db.Prepare("DELETE from modell where id= ? ")
	if err != nil {
		panic(err.Error())
	}
	defer getback.Close()

	if _, err := getback.Exec(delmodellid); err != nil {
		panic(err.Error())
	}
}
//////////////////delete gertyp by id function//////////////////
func delgertyp(delgertypid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	getback, err := db.Prepare("DELETE from gertyp where id= ? ")
	if err != nil {
		panic(err.Error())
	}
	defer getback.Close()

	if _, err := getback.Exec(delgertypid); err != nil {
		panic(err.Error())
	}
}
//////////////////delete delfield by id function//////////////////
func delfield(delfieldpid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	getback, err := db.Prepare("DELETE from zinfo where id= ? ")
	if err != nil {
		panic(err.Error())
	}
	defer getback.Close()

	if _, err := getback.Exec(delfieldpid); err != nil {
		panic(err.Error())
	}
}
//////////////////get model name by id function//////////////////
func getmodinfo(modell string) string {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	var (
		modellb string
	)

	err1 := db.QueryRow("SELECT modell FROM modell where id = ? ", modell).Scan(&modellb)
	if err1 != nil {
		panic(err1.Error())
	}

	return modellb
}
//////////////////get item back by id function//////////////////
func getback(getbackid string, userid string) {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	getback, err := db.Prepare("UPDATE bestand SET ticketnr='#LAGER#', ausgabename='#LAGER#', ausgabedatum=NOW(), changed= ? where id= ? ")
	if err != nil {
		panic(err.Error())
	}
	defer getback.Close()

	if _, err := getback.Exec(userid,getbackid); err != nil {
		panic(err.Error())
	}
}
//////////////////get gertypname by id function//////////////////
func getgertypinfo(gertyp string) string {

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	var (
		gertypname string
	)

	err1 := db.QueryRow("SELECT gertyp FROM gertyp where id = ? ", gertyp).Scan(&gertypname)
	if err1 != nil {
		panic(err1.Error())
	}
	return gertypname
}
//////////////////get user by id function//////////////////
func getuser(getuserid string) string {

	var getuseridcheck int

	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	err1 := db.QueryRow("SELECT COUNT(*) FROM login where id = ? ", getuserid).Scan(&getuseridcheck)
	if err1 != nil {
		panic(err1.Error())
	}

	
	if getuseridcheck > 0 {

		err2 := db.QueryRow("SELECT username FROM login where id = ? ", getuserid).Scan(&getuserid)
		if err2 != nil {
			panic(err2.Error())
		}
	
		return getuserid

	} else
	{
	
		return `User `+getuserid+` wurde gelöscht.`
	
	}

}

//////////////////get user by id function//////////////////
func getgertyp(getgertypid string) string {
	db, err := sql.Open("mysql", sqlcred)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	err1 := db.QueryRow("SELECT gertyp FROM gertyp where id = ? ", getgertypid).Scan(&getgertypid)
	if err1 != nil {
		panic(err1.Error())
	}
	return getgertypid
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


//////////////////struct recapv3//////////////////
type jDaten struct {
	Success     bool
	ChallengeTS string
	Hostname    string
	Score       float64
	Action      string
}


//////////////////get recapv3 function//////////////////
func OnPage(id string, aktion string)(string) {

	var link string = `https://www.google.com/recaptcha/api/siteverify?secret=`+pivcap+`&response=`+id+``

	res, err := http.Get(link)
	if err != nil {
		panic(err.Error())
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err.Error())
	}

	var daten jDaten
	Data := []byte(string(content))
	errb := json.Unmarshal(Data, &daten) 
	if errb != nil {
		panic(err.Error())
	}
	
	var output string = "0"
	
	if (daten.Action == aktion && daten.Success == true && daten.Score > scorecap) {
		output = "1"	
	}
	
//	fmt.Println(link)
//	fmt.Println(daten)
//	fmt.Println(output)	
	return output
}
