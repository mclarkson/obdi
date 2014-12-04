#!/bin/bash
#
# Obdi - a REST interface and GUI for deploying software
# Copyright (C) 2014  Mark Clarkson
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.


export PATH=/bin:/usr/bin

if ! which salt >& /dev/null; then
    echo '{"Error":"Salt binary not found."}'
    exit 1
fi

[[ -z $1 ]] && {
  echo '{"Error":"First argument, salt_id, was not set"}'
  exit 1
}

salt_id=$1
shift

[[ -z $1 ]] && {
  echo '{"Error":"No grains were specified"}'
  exit 1
}

for i in "$@"; do
  cmd="salt -t 60 --output=json $salt_id grains.setval ${i%%,*} ${i##*,}"
  output=`env $cmd`

  [[ -z $output ]] && {
    echo "{\"Error\",\"No output from command. $cmd\""
    exit 1
  }

  numlines=`echo "$output" | wc -l`
  [[ $numlines -ne 1 ]] && {
    echo "{\"Error\",\"Expected 1 line of output. Got $numlines.\"}"
    echo "{\"Output\",\"$output\"}" | sed -n '1{h;n}; {H} ; $ {g;s/\n/\\n/g;p}'
    exit 2
  }

  echo "$output"
done

exit 0
