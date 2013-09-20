#!/bin/sh -x

# call this script from within vm to compile and run newspeak

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
export GOPATH=$(pwd)
go install newspeak

# hot restart - does not work yet
#PID=$(ps ax|grep newspeak|head -n1|awk '{print $1}')
#kill -SIGHUP $PID

/etc/init.d/newspeak stop
/etc/init.d/newspeak start
tail -f /var/log/newspeak/newspeak.log 
