#!/bin/bash

export PATH=/bin:/usr/bin

[[ -z "$1" ]] && {
    echo '{"Error":"Argument missing. Expected Arg1 - JobID."}'
    exit 1
}

[[ ! -z $2 ]] && {
    echo '{"Error":"Too many arguments. Expected Arg1 - server name, Arg2 JobID."}'
    exit 1
}

if ! which salt >& /dev/null; then
    echo '{"Error":"Salt binary not found."}'
    exit 1
fi

JOBID="$1"

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

# Get minion name (only one is expected and checked for)

result=$(salt_run_json jobs.list_job $JOBID)
minion=$( echo "$result" | decode_json | sed -n "s/\[\"Minions\",0\]\s*//p")

[[ -z $minion ]] && {
    echo -n '{"Error":"Could not find the job. Maybe it is has been deleted.'
    echo -n " Consider setting the 'keep_jobs' variable in /etc/salt/master"
    echo -n " to a higher value to keep job history for longer, since job"
    echo ' history is not saved in the Obdi databases."}'
    exit 1
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
        exit 0
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
            exit 0
        }

        # No output from salt-run or salt list_job so there's a problem

        echo -n '{"Error":"The job has no result output and is not running on the minion.'
        echo -n ' If the salt server or minion is under heavy load then the job'
        echo -n ' might not have started yet, so it might be worth querying again.'
        echo -n ' If the problem persists then the minion might have been down or died when'
        echo ' the job was submitted. In that case check the minion then submit the job again."}'
        exit 1
    }

    # The job is still running on the minion so check again...

done

# Script never reaches here so the next line is redundant
exit 0

