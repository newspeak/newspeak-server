newspeak on amazon web services
===============================

deploy to amazon web services (aws) using a load balancer.

the setup described here allows for use of multiple aws accounts. most files are prefixed with a string for the account used, in this example the account prefix is `aws-api`.

configuration
-------------

 * create a private key (pk), register it with aws and place in
   
       `~/.ec2/aws-api-pk-ABCDEFGHIJKLMNOPQRSTUVXYZ1234567.pem`

   the string 'ABCDEFGHIJKLMNOPQRSTUVXYZ1234567' should be the id of your key.

 * download the aws certificate (cert) for your account and place in

       `~/.ec2/aws-api-cert-ABCDEFGHIJKLMNOPQRSTUVXYZ1234567.pem`

   the string 'ABCDEFGHIJKLMNOPQRSTUVXYZ1234567' should be the id of your key.

 * create an aws security group with inbound ports 22 (ssh) and 80 (http) and name it: `22+80`. this exact name is important since it is used in configuration files.

 * create a simple script called `~/.ec2/api.sh` to set all required shell environment variables. this way, scripts commited to git do not have to contain any credentials.

       echo '== api =='

       export AWS_ACCOUNT='aws-api'
       export EC2_URL='https://ec2.eu-west-1.amazonaws.com'
       export AWS_ACCESS_KEY="ABCDEFGHIJKLMNOPQRST"
       export AWS_SECRET_KEY="abcdefghijklmnopqrstuvwxyz1234567890ABCD"

       export AWS_ELB_HOME="/usr/local/Library/LinkedKegs/elb-tools/jars"
       export EC2_PRIVATE_KEY="$(/bin/ls "$HOME"/.ec2/aws-api-pk*.pem | /usr/bin/head -1)"
       export EC2_CERT="$(/bin/ls "$HOME"/.ec2/aws-api-cert*.pem | /usr/bin/head -1)"

 * edit `Vagrantfile` to set aws instance type or region.

usage
-----

create a new aws instance and install the api server on it, run:

    cd newspeak-server/api-aws
    ./vagrant-setup-api.sh

log into an running aws instance:

    cd newspeak-server/api-aws
    vagrant ssh

shut down, destroy and delete all data on an aws instance:

    cd newspeak-server/api-aws
    vagrant destroy

