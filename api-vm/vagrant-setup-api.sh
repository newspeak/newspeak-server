#!/bin/sh -x

# prepare vm and install api server.
#
# depends on the vagrant plugin vbguest. install it with:
# `vagrant plugin install vagrant-vbguest`

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

cd $(dirname $0)

vagrant up --no-provision

# remove outdated guest additions and chef
vagrant ssh -c "sudo apt-get -y purge virtualbox-guest-dkms virtualbox-guest-utils virtualbox-guest-x11 chef"

# remove some more unneeded packages
vagrant ssh -c "sudo apt-get -y --purge autoremove"

# rename default user
#user=newspeak
#vagrant ssh -c "sudo usermod  -l $user ubuntu"
#vagrant ssh -c "sudo groupmod -n $user ubuntu"
#vagrant ssh -c "sudo usermod  -d /home/$user -m $user"
#vagrant ssh -c "sudo mv /etc/sudoers.d/90-cloud-init-users /etc/sudoers.d/90-cloud-init-$user"
#vagrant ssh -c "sudo perl -pi -e \"s/ubuntu/$user/g;\" /etc/sudoers.d/90-cloud-init-$user"

# update guest additions
vagrant vbguest

# do clean reboot. you have to manually check if guest addition versions match on startup
vagrant halt
vagrant up --no-provision

# copy over cached files to speed up compilation
../shared/scripts/vagrant/scp.sh -r ../cache/deb/*.deb default:/home/vagrant
../shared/scripts/vagrant/scp.sh -r ../cache/maven/* default:/home/vagrant/.m2
vagrant ssh -c "sudo mv /home/vagrant/deb/*.deb /var/cache/apt/archives; rm -rf /home/vagrant/deb"

# update virtualbox guest additions
../shared/scripts/vagrant/setup-basebox.sh

# start api server installation
vagrant ssh -c "/vagrant/api/api-install.sh"

# update cache to speed up subsequent installs
rm -rf ../cache
mkdir -p ../cache/deb ../cache/maven
../shared/scripts/vagrant/scp.sh -r default:/var/cache/apt/archives/*.deb ../cache/deb
../shared/scripts/vagrant/scp.sh -r default:/home/vagrant/.m2/* ../cache/maven
