#!/bin/sh

# scp to a vagrant machine.

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

OPTIONS=`vagrant ssh-config | awk -v ORS=' ' '{print "-o " $1 "=" $2}'`
 
scp ${OPTIONS} "$@" || echo "Transfer failed. Did you use 'default:' for the vagrant host? If using aws: did you source an aws config?"
