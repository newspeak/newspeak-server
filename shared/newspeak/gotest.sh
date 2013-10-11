#!/bin/sh

# run go test recursively or do benchmarks

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
echo "be verbose with: $(basename $0) -v"
echo "run benchmarks with: $(basename $0) -b"
echo
if [ "$1" = "-b" ]; then
  go test ./src/newspeak/... -v -gocheck.b | grep Benchmark
else
  go test ./src/newspeak/... $1
fi
