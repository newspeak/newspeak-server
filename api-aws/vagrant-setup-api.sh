#!/bin/sh -x

cd $(dirname $0)

# source aws credentials
. ~/.ec2/api.sh

# prepare files to be copied over to aws instance. all files in this directory will be rsynced to aws.
rm -rf newspeak live
mkdir -p newspeak/src
cp -r ../shared/newspeak/src/newspeak newspeak/src
cp -r ../shared/live .

# start up instance
vagrant up --provider=aws

# assign elastic ip
export INSTANCE=$(ec2-describe-instances --region eu-west-1 | grep running | awk '{print $2}')
export IP=$(ec2-describe-addresses | awk '{print $2}')
echo "assigning ip $IP to instance $INSTANCE"
ec2-associate-address -i $INSTANCE $IP

# start api server install script
vagrant ssh -c "/vagrant/api/api-install.sh"

# add to load balancer
export INSTANCE=$(ec2-describe-instances --region eu-west-1 | grep running | awk '{print $2}')
elb-register-instances-with-lb loadbalancer --instances $INSTANCE --region eu-west-1

# remove old instance from load balancer
export INSTANCE=$(elb-describe-instance-health loadbalancer --region eu-west-1 | grep OutOfService | head -n1 | awk '{print $2}')
elb-deregister-instances-from-lb loadbalancer --instances $INSTANCE --region eu-west-1
