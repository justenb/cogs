# cogs 

"swiss army knife" for secure shell based server operations
<img src="http://i62.tinypic.com/1qlmv4.png" width="590" height="400"  />

<br><b>get uptime from host(s)</b></br>
cogs -hosts 127.0.0.1 -key ~/.ssh/id_rsa.pub -cmd uptime

<br><b>download files using glob matching </b></br>
cogs -hosts 127.0.0.1 -key ~/.ssh/id_rsa.pub -get /var/www/html/webfiles/*

<br><b>put file on host(s)</b></br>
cogs -hosts 127.0.0.1 -key ~/.ssh/id_rsa.pub -put index.html -rpath /var/www/html/index.html

<br><b>get stream filtered dmesg output</b></br>
cogs -hosts 127.0.0.1 -key ~/.ssh/id_rsa.pub -cmd dmesg -filter panic

<br><b>download filtered log messages</b></br>
cogs -hosts 127.0.0.1 -key ~/.ssh/id_rsa.pub -get /var/www/html/webfiles/logs/* -filter critical
