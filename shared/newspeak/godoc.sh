#!/bin/sh -x

# generate and display documentation for this project at
# http://localhost:6060/pkg/newspeak

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
godoc -http=:6060

# to generate offline documentation for use without godoc, uncomment lines
# below
# based on https://code.google.com/p/go/issues/detail?id=2381#c3
#wget -r -np -N -E -p -k -e robots=off http://localhost:6060/pkg/newspeak
#mv ./localhost\:6060 doc
