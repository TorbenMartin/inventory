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


(3) create and install mysql database:

	Restore the SQL Backup:

		mysql -u root -p
		source lager.sql;
			
	Backup created by:

		mysqldump -u root -p -x -B lager > lager.sql


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

