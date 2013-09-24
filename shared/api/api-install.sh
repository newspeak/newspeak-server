#!/bin/sh -x

# install newspeak api server on ubuntu server
#
# the script is supposed to work in three environments:
#  * in a vagrant vm
#  * on a amazon aws instance
#  * on a root server
# all environment specific setup is supposed to happen at the beginning of the
# script.
# some commands may fail in some environments - this is intentended behavior.

# newspeak.io
# Copyright (C) 2013 Jahn Bertsch
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License version 3
# as published by the Free Software Foundation.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

export USERNAME=$(whoami)

if [ "$USERNAME" = "root" ]; then
  # root server specific
  export FILES=/root
  echo api > /etc/hostname
  sed -i 's/Ubuntu-1304-raring-64-minimal/api/' /etc/hosts
  service networking restart
  hostname api
else
  # vagrant and aws specific
  export FILES=/vagrant
fi

echo FILES=$FILES
echo USERNAME=$USERNAME

# timezone
echo 'Europe/Berlin' | sudo tee /etc/timezone
sudo dpkg-reconfigure --frontend noninteractive tzdata
sudo /etc/init.d/ntp restart


echo ===== start =====
date

# remove consolekit on aws
sudo apt-get purge consolekit -y

# update and cleanup
sudo apt-get update
sudo apt-get upgrade -y
sudo apt-get autoremove

# install and configure dynamic login message, motd (already done on aws)
sudo apt-get install update-motd update-notifier-common landscape-common -y
sudo rm -f /etc/update-motd.d/10-help-text /etc/update-motd.d/51-cloudguest

# enable unattended upgrades
sudo echo "get unattended-upgrades/enable_auto_updates" | sudo debconf-communicate
sudo apt-get install unattended-upgrades
sudo echo 'unattended-upgrades unattended-upgrades/enable_auto_updates boolean true' | sudo debconf-set-selections
sudo DEBIAN_FRONTEND=noninteractive dpkg-reconfigure -plow unattended-upgrades
echo "get unattended-upgrades/enable_auto_updates" | sudo debconf-communicate


echo ===== maven =====
# (incl maven workaround from https://bugs.launchpad.net/ubuntu/+source/wagon2/+bug/1171056)
sudo apt-get install -y openjdk-7-jdk maven
sudo dpkg -i --force-all /var/cache/apt/archives/libwagon2-java_*_all.deb
sudo apt-get install -y maven


echo ===== usergrid =====
sudo apt-get install -y git-core
cd
rm -rf usergrid-stack
git clone https://github.com/apigee/usergrid-stack.git
cd usergrid-stack
git reset --hard fb4fe7d0fd3cbe950869abf7849d7f558513e801 # 2013-09-20
sudo cp $FILES/api/usergrid-default.properties $HOME/usergrid-stack/config/src/main/resources

# start usergrid compilation
mvn clean install -DskipTests=true -Dusergrid-custom-spring-properties=classpath:/usergrid-custom.properties

# instead of compilation, download precompiled .war
#wget https://*.s3.amazonaws.com/usergrid/ROOT.war
#mkdir -p /var/lib/tomcat7/webapps
#mv ROOT.war /var/lib/tomcat7/webapps

# /config/src/main/resources/usergrid-custom.properties
# /rest/src/main/resources/usergrid-custom.properties

# install usergrid in tomcat
sudo aptitude install -y tomcat7
sudo rm -rf /var/lib/tomcat7/webapps/ROOT
sudo cp $HOME/usergrid-stack/rest/target/ROOT.war /var/lib/tomcat7/webapps

# stop until setup is complete to free up memory
sudo /etc/init.d/tomcat7 stop

cd ..

echo ===== vpn =====
# tinc vpn
sudo aptitude install tinc -y

sudo bash -c "cat >> /etc/network/interfaces" <<EOF

# tinc
auto vpn
iface vpn inet static
        address 192.168.11.1
#        use address 192.168.11.1 on node1 (running cassandra)
#        use address 192.168.12.1 on node2
        netmask 255.255.0.0
        tinc-net newspeak
        tinc-debug 1
#       tinc-mlock yes
#        see http://www.tinc-vpn.org/pipermail/tinc/2012-September/003056.html
        tinc-user nobody
        tinc-logfile /var/log/tinc
EOF

sudo bash -c "cat >> /etc/tinc/nets.boot" <<EOF
newspeak
EOF

# generate key
# open question: is at this stage of the boot process alread enough entropy available from the linux kernel to generate a good key?
# see: "Mining Your Ps and Qs: Detection of Widespread Weak Keys in Network Devices",
# https://factorable.net/weakkeys12.extended.pdf
sudo mkdir -p /etc/tinc/newspeak
cd /etc/tinc/newspeak
sudo openssl genrsa -out rsa_key.priv 8192
sudo openssl rsa -in rsa_key.priv -pubout > rsa_key.pub

sudo bash -c "cat > /etc/tinc/newspeak/tinc.conf" <<EOF
Device=/dev/net/tun
Interface=vpn
AddressFamily=ipv4
Mode=switch
Name=cassandraMain
PrivateKeyFile=/etc/tinc/newspeak/rsa_key.priv
#ConnectTo=cassandraMirror1
EOF

sudo mkdir -p /etc/tinc/newspeak/hosts

sudo bash -c "cat > /etc/tinc/newspeak/hosts/cassandraMain" <<EOF
Address=192.168.11.1 # use elastic ip if on aws
Port=54321
Subnet=192.168.0.0/16
Compression=9
EOF

sudo cat /etc/tinc/newspeak/rsa_key.pub >> /etc/tinc/newspeak/hosts/cassandraMain

sudo bash -c "cat >> /etc/default/tinc" <<EOF
EXTRA=""
EOF

sudo ifup vpn


echo ===== cassandra =====
sudo su -c 'echo >> /etc/apt/sources.list'
sudo su -c 'echo "# cassandra" >> /etc/apt/sources.list'
sudo su -c 'echo "deb http://debian.datastax.com/community stable main" >> /etc/apt/sources.list'
wget -O - http://debian.datastax.com/debian/repo_key | sudo apt-key add -
sudo apt-get update
sudo apt-get install -y cassandra

# stop until setup is complete to free up memory
sudo /etc/init.d/cassandra stop

# update cluster config
sudo cp /vagrant/api/cassandra.yaml /etc/cassandra/cassandra.yaml

# delete data containing old cluster config
sudo rm -rf /var/lib/cassandra/data

# reduce memory usage of cassandra and tomcat for aws by overwriting memory default value
# sudo cp $FILES/api/cassandra-env.sh /etc/cassandra/
# sudo cp $FILES/api/tomcat7.default /etc/default/


echo ===== uniqush =====

# install a recent version of go
cd
wget https://go.googlecode.com/files/go1.1.2.linux-amd64.tar.gz
tar xf go*.linux-amd64.tar.gz
rm go*linux-amd64.tar.gz
sudo mv go /usr/local
sudo ln -s /usr/local/go/bin/go /usr/local/bin
sudo ln -s /usr/local/go/bin/gofmt /usr/local/bin
sudo ln -s /usr/local/go/bin/godoc /usr/local/bin

# alternatively, install the version provided by the package manager 
#sudo DEBIAN_FRONTEND='noninteractive' apt-get -o Dpkg::Options::='--force-confnew' -y install golang

# install redis and uniqush
sudo apt-get install -y redis-server mercurial
cd
git clone https://github.com/uniqush/uniqush-push.git
cd uniqush-push
git reset --hard 1.5.0
export GOPATH=$(pwd)
go get code.google.com/p/goconf/conf # requires mercurial
go get github.com/nu7hatch/gouuid
go get github.com/uniqush/log
go get github.com/uniqush/uniqush-push/db
go get github.com/uniqush/uniqush-push/push
go get github.com/uniqush/uniqush-push/srv
go build
sudo mkdir -p /etc/uniqush
sudo cp conf/uniqush-push.conf /etc/uniqush
sudo cp uniqush-push /usr/local/bin

# stop until setup is complete to free up memory
sudo /etc/init.d/redis-server stop


echo ===== go api server =====

sudo apt-get install -y memcached

# initscript
sudo cp $FILES/api/newspeak.initscript /etc/init.d/newspeak
sudo chmod 755 /etc/init.d/newspeak
cd /etc/init.d
sudo update-rc.d newspeak defaults

# log
sudo mkdir -p /var/log/newspeak
sudo touch /var/log/newspeak/newspeak.log
sudo chown -R $USERNAME:$USERNAME /var/log/newspeak

# go dependencies for newspeak
mkdir -p $FILES/newspeak
cd $FILES/newspeak/
export GOPATH=$(pwd)
go get github.com/fitstar/falcore
go get github.com/Mistobaan/go-apns
go get github.com/orfjackal/gospec/src/gospec
go get github.com/peterbourgon/g2s
go get github.com/bradfitz/gomemcache/memcache

# build and install
go install newspeak
sudo cp bin/newspeak /usr/local/bin

# set apns certificates
sudo mkdir -p /etc/newspeak/apns-certs
sudo cp -r $FILES/api/apns-certs-production/* /etc/newspeak/apns-certs

# bitly statsdaemon
cd
mkdir -p statsdaemon
cd statsdaemon
export GOPATH=$(pwd)
go get github.com/bitly/statsdaemon
sudo cp bin/statsdaemon /usr/local/bin


echo ===== usergrid setup =====

sudo /etc/init.d/tinc start
sudo /etc/init.d/tomcat7 start
sudo /etc/init.d/cassandra start

# setup db
sudo apt-get install curl -y
curl --user superuser:superuser http://localhost:8080/system/database/setup

# simple test - should not return any errors
curl localhost:8080/status
curl localhost:8080/test/hello
echo checking if cassandra is available from usergrid:
curl -s localhost:8080/status | grep cassandraAvailable

# create organization
curl -X POST  \
     -d 'organization=newspeak&username=admin&name=admin&email=admin@example.com&password=admin' \
     http://localhost:8080/management/organizations

# get admin auth token
export ADMINTOKEN=$(curl 'http://localhost:8080/management/token?grant_type=password&username=admin&password=admin' | cut -f 2 -d , | cut -f 2 -d : | cut -f 2 -d \")

# create app
curl -H "Authorization: Bearer $ADMINTOKEN" \
     -H "Content-Type: application/json" \
     -X POST -d '{ "name":"newspeak" }' \
     http://localhost:8080/management/orgs/newspeak/apps

# create user
curl -H "Authorization: Bearer $ADMINTOKEN" \
     -X POST "http://localhost:8080/newspeak/newspeak/users" \
     -d '{ "username":"user", "password":"user", "email":"user@example.com" }'

echo ===== restart services =====

sudo /etc/init.d/tomcat7 stop
sudo /etc/init.d/cassandra stop
sudo /etc/init.d/redis-server stop
sudo /etc/init.d/newspeak stop

/etc/init.d/newspeak start
sudo /etc/init.d/cassandra start
sudo /etc/init.d/redis-server start
sudo /etc/init.d/tomcat7 start
sudo /etc/init.d/newspeak start

echo ===== completed =====
date
