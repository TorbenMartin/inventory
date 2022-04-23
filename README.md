# inventory
a small IT inventory management system as a opensource project.



-----------------------------------------------------------------



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

	- for linux: go build -o lager *.go
	- for windows: GOOS=windows GOARCH=amd64 go build -o lager.exe *.go

(5) normal run:

	- run source: go run *.go
	
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

