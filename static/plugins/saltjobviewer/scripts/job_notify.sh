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

PIDFILE=/var/run/job_notify.sh.pid
ME="$0"

# ---------------------------------------------------------------------------
quit()
# ---------------------------------------------------------------------------
{
  rm -f $PIDFILE
}

# ---------------------------------------------------------------------------
write_pid_or_exit()
# ---------------------------------------------------------------------------
{
  local me=`basename $ME`
  if test -e $PIDFILE; then
    pid=`cat $PIDFILE`
    # Stale pid?
    if ps ax | grep -qs "^$pid.*$me"; then
      exit 0
    fi
  fi
  echo "$$" >$PIDFILE
  trap quit EXIT
}

write_pid_or_exit

# ---------------------------------------------------------------------------
# JSON decoding
# ---------------------------------------------------------------------------

throw () {
  echo "$*" >&2
  exit 1
}

BRIEF=1
LEAFONLY=1
PRUNE=1

awk_egrep () {
  local pattern_string=$1

  gawk '{
    while ($0) {
      start=match($0, pattern);
      token=substr($0, start, RLENGTH);
      print token;
      $0=substr($0, start+RLENGTH);
    }
  }' pattern=$pattern_string
}

tokenize () {
  local GREP
  local ESCAPE
  local CHAR

  if echo "test string" | egrep -ao --color=never "test" &>/dev/null
  then
    GREP='egrep -ao --color=never'
  else
    GREP='egrep -ao'
  fi

  if echo "test string" | egrep -o "test" &>/dev/null
  then
    ESCAPE='(\\[^u[:cntrl:]]|\\u[0-9a-fA-F]{4})'
    CHAR='[^[:cntrl:]"\\]'
  else
    GREP=awk_egrep
    ESCAPE='(\\\\[^u[:cntrl:]]|\\u[0-9a-fA-F]{4})'
    CHAR='[^[:cntrl:]"\\\\]'
  fi

  local STRING="\"$CHAR*($ESCAPE$CHAR*)*\""
  local NUMBER='-?(0|[1-9][0-9]*)([.][0-9]*)?([eE][+-]?[0-9]*)?'
  local KEYWORD='null|false|true'
  local SPACE='[[:space:]]+'

  $GREP "$STRING|$NUMBER|$KEYWORD|$SPACE|." | egrep -v "^$SPACE$"
}

parse_array () {
  local index=0
  local ary=''
  read -r token
  case "$token" in
    ']') ;;
    *)
      while :
      do
        parse_value "$1" "$index"
        index=$((index+1))
        ary="$ary""$value" 
        read -r token
        case "$token" in
          ']') break ;;
          ',') ary="$ary," ;;
          *) throw "EXPECTED , or ] GOT ${token:-EOF}" ;;
        esac
        read -r token
      done
      ;;
  esac
  [ "$BRIEF" -eq 0 ] && value=`printf '[%s]' "$ary"` || value=
  :
}

parse_object () {
  local key
  local obj=''
  read -r token
  case "$token" in
    '}') ;;
    *)
      while :
      do
        case "$token" in
          '"'*'"') key=$token ;;
          *) throw "EXPECTED string GOT ${token:-EOF}" ;;
        esac
        read -r token
        case "$token" in
          ':') ;;
          *) throw "EXPECTED : GOT ${token:-EOF}" ;;
        esac
        read -r token
        parse_value "$1" "$key"
        obj="$obj$key:$value"        
        read -r token
        case "$token" in
          '}') break ;;
          ',') obj="$obj," ;;
          *) throw "EXPECTED , or } GOT ${token:-EOF}" ;;
        esac
        read -r token
      done
    ;;
  esac
  [ "$BRIEF" -eq 0 ] && value=`printf '{%s}' "$obj"` || value=
  :
}

parse_value () {
  local jpath="${1:+$1,}$2" isleaf=0 isempty=0 print=0
  case "$token" in
    '{') parse_object "$jpath" ;;
    '[') parse_array  "$jpath" ;;
    # At this point, the only valid single-character tokens are digits.
    ''|[!0-9]) throw "EXPECTED value GOT ${token:-EOF}" ;;
    *) value=$token
       isleaf=1
       [ "$value" = '""' ] && isempty=1
       ;;
  esac
  [ "$value" = '' ] && return
  [ "$LEAFONLY" -eq 0 ] && [ "$PRUNE" -eq 0 ] && print=1
  [ "$LEAFONLY" -eq 1 ] && [ "$isleaf" -eq 1 ] && [ $PRUNE -eq 0 ] && print=1
  [ "$LEAFONLY" -eq 0 ] && [ "$PRUNE" -eq 1 ] && [ "$isempty" -eq 0 ] && print=1
  [ "$LEAFONLY" -eq 1 ] && [ "$isleaf" -eq 1 ] && \
    [ $PRUNE -eq 1 ] && [ $isempty -eq 0 ] && print=1
  [ "$print" -eq 1 ] && printf "[%s]\t%s\n" "$jpath" "$value"
  :
}

parse () {
  read -r token
  parse_value
  read -r token
  case "$token" in
    '') ;;
    *) throw "EXPECTED EOF GOT $token" ;;
  esac
}

decode_json() {
  tokenize | parse
}

# ---------------------------------------------------------------------------
# Salt-run to output json
# ---------------------------------------------------------------------------

salt_run_json() {

    local fun="$1"
    local arg="$2"

    output=`python < <(cat <<EnD
import json
import os
import sys

import salt.config
import salt.runner

if __name__ == '__main__':
    __opts__ = salt.config.master_config(
            os.environ.get('SALT_MASTER_CONFIG', '/etc/salt/minion'))
    runner = salt.runner.Runner(__opts__)

    stdout_bak = sys.stdout
    with open(os.devnull, 'wb') as f:
        sys.stdout = f
        ret = runner.cmd("$fun", ["$arg"])
    sys.stdout = stdout_bak

    print json.dumps(ret)
EnD
    )`

    echo "$output"
}

# ---------------------------------------------------------------------------
# Job check
# ---------------------------------------------------------------------------

job_header_check() {
    result=$(salt_run_json jobs.list_job $JOBID)
    minion=$( echo "$result" | decode_json | sed -n "s/\[\"Minions\",0\]\s*//p")

    [[ -z $minion ]] && {
        # Job has no minions, makes no sense
        NOOUTPUT=1
        return 1
    }

    result=$(salt_run_json jobs.list_job $JOBID)

    echo "$result"
}

job_check() {
    result=$(salt_run_json jobs.list_job $JOBID)
    minion=$( echo "$result" | decode_json | sed -n "s/\[\"Minions\",0\]\s*//p")

    [[ -z $minion ]] && {
        # Job has no minions, makes no sense
        NOOUTPUT=1
        return 1
    }

    while true; do
        # Try to get a result 10 times (probably around 60 seconds)
        for i in 1 2 3 4 5 6 7 8 9 0; do
            gotresult=$(echo "$result" | decode_json | grep -m1 "^\[\"Result\",")
            [[ ! -z $gotresult ]] && break
            sleep 5
            result=$(salt_run_json jobs.list_job $JOBID)
        done

        [[ ! -z $gotresult ]] && {
            echo "$result"
            return 0
        }

        # Check if it's still running on the minion

        isrunning=$(salt "$minion" --output=json saltutil.find_job $JOBID | \
            decode_json | grep -m1 "^\[\"'"$minion"'\",\"jid\"\]")

        [[ -z $isrunning ]] && {

            # It isn't running on the minion. Do a single check to see
            # if there has been a result since checking the minion

            result=$(salt_run_json jobs.list_job $JOBID)
            gotresult=$(echo "$result" | decode_json | grep -m1 "^\[\"Result\",")

            [[ ! -z $gotresult ]] && {
                echo "$result"
                return 0
            }

            # No output from salt-run or salt list_job so there's a problem
            NOOUTPUT=1
            return 1
        }

        # The job is still running on the minion so check again...

    done
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

STATUS_UNKNOWN=0
STATUS_OK=1
STATUS_ERRORS_FOUND=2
STATUS_NO_OUTPUT=3

# Get the number of hours that jobs are kept.
# This is sent in the status so saltjobstatus can prune the list.
# keep_jobs will be zero if there was an error.
declare -i keep_jobs
keep_jobs=`salt-call --output=newline_values_only config.get keep_jobs 2>/dev/null`

send_status() {
    local status=$1

    echo "Sending status $status"

    [[ ! -e "/etc/obdi/job_status.conf" ]] && {
      echo ".. Could not send. /etc/obdi/job_status.conf is missing."
      return
    }
    source /etc/obdi/job_status.conf

    proto="https"
    opts="-k -s"
    ipport="127.0.0.1:443"
    guid=`curl $opts -d '{"Login":"'"$STATUS_USER"'","Password":"'"$STATUS_PASS"'"}' \
        $proto://$ipport/api/login | grep -o "[a-z0-9][^\"]*"`
    curl $opts -d '{"JobId":"'"$JOBID"'","Status":'"$status"',"KeepJobs":'"$keep_jobs"'}' \
        $proto://$ipport/api/jobstatus/$guid/saltjobviewer/saltjobstatus
    echo
}

#echo "/var/cache/salt/master/jobs/6d/af04e79a4419d80b6ceef81237df5e" | \
dir="/var/cache/salt/master/jobs"
numslashes_base=`echo -n "$dir" | sed 's#/$##' | sed 's#[^/]##g' | wc -c`
inotifywait -q -e isdir --format '%w%f' -m -r $dir/* | \
while read line; do
    # Only look for jid file 2 levels down from the base dir
    numslashes=`echo -n "$line" | sed 's#/$##' | sed 's#[^/]##g' | wc -c`
    level=$((numslashes-numslashes_base))
    [[ $level -eq 2 ]] && {
        for dummy in 1 2 3 4 5; do
            [[ -e "$line/jid" ]] && {
                JOBID=`cat "$line/jid"`
                json=`job_header_check | decode_json`
                if echo $json | grep -qs '\["Function"\].*"state.highstate"'; then
                    json=`job_check | decode_json`
                    if [[ $NOOUTPUT -eq 1 ]]; then
                        send_status $STATUS_NO_OUTPUT
                    else
                        # It's a highstate
                        #nsuccess=`echo "$json" | \
                        #    grep -E -c '\["Result",.*,"result"\]\s*true'`
                        nfail=`echo "$json" | \
                            grep -E -c '\["Result",.*,"result"\]\s*false'`
                        if [[ $nfail -ge 1 ]]; then
                            send_status $STATUS_ERRORS_FOUND
                        else
                            send_status $STATUS_OK
                        fi
                    fi
                fi
                break
            }
            sleep 2
        done
    }
done

