newspeak api server
===================

for creating a predictable working environment during development and testing, the server is run as a virtual machine (vm). for managment of the vm, [virtualbox](http://virtualbox.com) and [vagrant](http://vagrantup.com) are used.

vagrant supports deployment to [amazon aws](http://aws.amazon.com) in a way similar to when used with virtualbox.

this leads to the following directory structure:

 * `newspeak-server/api-vm`: api server for use with vagrant/virtualbox vm
 * `newspeak-server/api-aws`: api server for use with amazon aws
 * `newspeak-server/shared`: files shared between the above machine setups (code, config files, ..)

each of the directories `api-vm` and `api-aws` contain a file called `Vagrantfile` with server configuration.

the directory `newspeak-server/shared/newspeak` contains the actual server sources. see the readme file there for more information.

dependencies
------------

to run the server locally, installation of the following software is needed:

 * [virtualbox](http://virtualbox.com)
 * [vagrant](http://vagrantup.com) - only a recent version (1.2.x) is supported! download from [here](http://downloads.vagrantup.com).
 * install vagrant plugins with: `vagrant plugin install vagrant-aws vagrant-vbguest`

to test sending push messages to ios devices, you have to be member in the [ios developer program](https://developer.apple.com/programs/ios/). receiving push messages is not supported by the ios simulator so will need a real ios device for that. however, you can still test large parts of the system without device or developer program membership.

quick start
-----------

for the impatient: set up an api server vm instance with two commands:

    git clone git@github.com:newspeak/newspeak-server.git
    newspeak-server/api-vm/vagrant-setup-api.sh

api server overview
-------------------

the api server consists of the following components:


                           +--------------------------------+
                           |                                |
        +------+   +-------+-------+   +----------+   +-----+-----+
        | inet |---| go api server |---| usergrid |---| cassandra |
        +------+   +-------+-------+   +----------+   +-----------+
                           |
                           |           +---------+
                           +-----------| uniqush |
                                       +---------+


main entry point is the `go api server`, listening on port 80 (http) and exposing a [rest](https://en.wikipedia.org/wiki/Representational_state_transfer) inspired api. it is written in [go](http://golang.org). the `go api server` sources are located at `newspeak-server/shared/newspeak`. some additional go resources, for reference only:

 * [go syntax](http://golang.org/doc/effective_go.html)
 * the [interfaces](http://golangtutorials.blogspot.de/2011/06/interfaces-in-go.html) concept in go
 * [inheritance](http://diveintogo.blogspot.de/2010/03/classic-inheritance-in-go-lang.html)

the `go api server` in turn uses a http based (rest) api to communicate with usergrid and uniqush. this loosely follows the "api facade" design pattern described in this [pdf](http://offers.apigee.com/api-design-ebook-bw). the components are:

 * [usergrid](http://usergrid.com): manages user signup, login, groups 
 * [uniqush](http://uniqush.org): send push messages to android and ios
 * [cassandra](http://cassandra.apache.org): scalable data store

the `go api server` may also communicate with cassandra directly.

the cassandra nodes communicate with each other using a vpn built with [tinc](http://www.tinc-vpn.org).


command line usage of the api server
------------------------------------

run these commands from your host machine. adjust the ip `192.168.1.77` to match the ip assigned to your api server vm.

 * set api server ip and port: `export HOST='192.168.1.77:80'`

 * request access token: `export ACCESS_TOKEN=$(curl $HOST/tokens --silent --data username=user --data password=user|grep AccessToken|cut -f 2 -d : |cut -f 2 -d \"); echo ACCESS_TOKEN=$ACCESS_TOKEN`

 * register device with push service (adjust device token): `curl $HOST/devices --include --data accessToken=$ACCESS_TOKEN --data username=user --data deviceToken=047595d0fae972fbed0cf1b3a41c7a549e0c475be1a9de5dc7f97cc185ade65d`

 * send push message: `curl $HOST/messages --include --data accessToken=$ACCESS_TOKEN --data recipient=user --data message="hello world"`


usergrid documentation
----------------------

[usergrid homepage](http://usergrid.com), [usergrid on github](https://github.com/apigee/usergrid-stack), [usergrid readme](http://github.com/apigee/usergrid-stack/blob/master/README.md)

usergrid has an admin interface which is entirely written in javascript and therefore completely runs in the browser. that means you can use the admin interface hostet on github to access your local machine. the url is (adjust ip to your setup)


    http://apigee.github.io/usergrid-portal/?api_url=http://192.168.1.77:8080


uniqush documentation
---------------------

[uniqush homepage](http://uniqush.org)

to send push messages with uniqush on the command line without proxying through the `go api server`, use these curl commands from within the api server vm:

 * register push service provider: `curl http://127.0.0.1:9898/addpsp -d service=newspeak -d pushservicetype=apns -d cert=/etc/newspeak/apns-certs/cert.pem -d key=/etc/newspeak/apns-certs/priv-noenc.pem -d sandbox=true`

 * register device (adjust device token): `curl http://127.0.0.1:9898/subscribe -d service=newspeak -d subscriber=user -d pushservicetype=apns -d devtoken=047595d0fae972fbed0cf1b3a41c7a549e0c475be1a9de5dc7f97cc185ade65d`

 * send message: `curl http://127.0.0.1:9898/push -d service=newspeak -d subscriber=user -d msg="hello world"`


apple push notification certificate creation
--------------------------------------------

the `go api server` requires a valid apple push notification service certificate corresponding to your app's bundle id.

create your certificate in the ios dev center according to [this](https://developer.apple.com/library/mac/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/Chapters/ProvisioningDevelopment.html#//apple_ref/doc/uid/TP40008194-CH104-SW1) guide in section 'creating the ssl certificate and keys'.

if you have the certificate imported into your keychain, follow these instructions on how to convert the apple push notification certificates from .p12 into .pem format as required by uniqush.

 * in osx keychain application filter by 'certificate' and select the 'Apple IOS Push Services' certificate you want to use.

 * expand disclosure triangle of certificate. you'll see the private key corresponding to the certificate.

 * export private key as 'priv.p12' and export certificate as 'cert.p12'. if you can't export as p12 you're not trying to export the correct thing!

 * set an empty password while exporting.

 * part 1: convert to certificate to pem: (enter empty import password)
openssl pkcs12 -clcerts -nokeys -out cert.pem -in cert.p12

 * part 2: convert private key to pem (enter empty import password, enter "1234" as pem pass phrase):
openssl pkcs12 -nocerts -out priv.pem -in priv.p12

 * remove password from above file (enter passphrase "1234"):
openssl rsa -in priv.pem -out priv-noenc.pem

 * you may now delete the file `priv.pem`, `priv.p12` and `cert.p12` since they are no longer required.
rm priv.pem priv.p12 cert.p12


place the resulting certificate `cert.pem` and private key `priv-noenc.pem` for the sandbox configuration in

 * `newspeak-server/shared/api/apns-certs-sandbox`

and for the production configuration in

 * `newspeak-server/shared/api/apns-certs-production`


ubuntu images
-------------

list of available [ubuntu aws instances](https://cloud-images.ubuntu.com/locator/ec2)

list of available [ubuntu vagrant base boxes](https://cloud-images.ubuntu.com/vagrant)


license
-------

    newspeak.io
    Copyright (C) 2013 Jahn Bertsch

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License version 3
    as published by the Free Software Foundation.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
