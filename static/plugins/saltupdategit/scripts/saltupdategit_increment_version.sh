#!/bin/bash

BRANCH="$1"
POS="$2"

# Where the cloned repo is temporarily stored
REPOTMP="/var/cache/obdi/gitrepo"

# Where the salt master config file is
SALT_MASTER_FILE="/etc/salt/master"

delete_repo_cache() {

    # REPOTMP does not exist so no need to delete
	[[ ! -e "$REPOTMP/$REPONAME" ]] && return

    # REPOTMP is not a directory
	[[ ! -d "$REPOTMP/$REPONAME" ]] && {
		echo -n '{"Error":"The location, `'"$REPOTMP/$REPONAME"'`, is not a'
		echo ' directory. Aborting."}'
		exit 1
	}

	[[ ! -e "$REPOTMP/$REPONAME/.git" ]] && {
		echo -n '{"Error":"Directory, `'"$REPOTMP/$REPONAME"'`, does not'
		echo -n ' contain a git repository, so was probably used for something'
		echo ' else. Will not delete the directory. Aborting."}'
		exit 1
	}

	rm -rf "$REPOTMP/$REPONAME/"
}

clone_repo_cache() {

	# Create the repo cache directory
	mkdir -p "$REPOTMP"

	# Update the repo first to pull upstream changes
	cd "$REPOHOME/$REPONAME.git"
	git fetch

	# Clone the repo
	cd "$REPOTMP"
	git clone "$REPOHOME/$REPONAME.git" >/dev/null
}

get_latest_version_string() {
	cd $REPOTMP/$REPONAME

	git branch -a | \
		sed -n 's#^\s*remotes/origin/'"$BRANCH"'_\(\)#\1#p' | \
		sort -nt . -k1,1 -k2,2 -k3,3 -k4,4 | \
		tail -n 1
}

increment_version_string() {
	# Arg1 - version, e.g.      0.1.10

	local version="$1"
	local -i pos=$POS
	local a
	
	IFS="." read -a a < <( echo "$version" )
	for i in 1 2 3; do
		if [[ $i == $pos ]]; then
			echo "$version" | awk -F. '{ printf( "%s", $'"$pos"'+1 ); }'
			a[$i]=0
			a[$i+1]=0
		else
			echo -n "${a[$i-1]}"
		fi
		[[ $i -lt 3 ]] && echo -n "."
	done
}

create_branch() {
	# Arg1 - version

	cd $REPOTMP/$REPONAME

	git checkout $BRANCH
	git branch ${BRANCH}_$1
	git push origin ${BRANCH}_$1

	/etc/init.d/salt-master restart
}

main() {
	local version new_version

	[[ -z $BRANCH ]] && {
		echo '{"Error":"Argument 1 missing. Expected a branch name."}'
		exit 1
	}

	[[ -z $POS ]] && {
		echo -n '{"Error":"Argument 2 missing. Expected a position in the'
        echo ' version string to increment."}'
		exit 1
	}

	# This script expects the lines to be in the following format
	# in the SALT_MASTER_FILE
	#
	#   ...
	#   gitfs_remotes:
	#     - file:///srv/fdp-mgmt-salt.git
	#   ...
	#
	# So in this case, 't', below, is loaded with: '/srv/fdp-mgmt-salt.git'

	t=`sed -n '/gitfs_remotes:/,/^[^ ]/ {/^ \+/{s#^.*://\(.*\)#\1#p}}' $SALT_MASTER_FILE`

	[[ -z $t ]] && {
		echo -n '{"Error":"Getting the gitfs_remotes setting from'
		echo ', $SALT_MASTER_FILE, failed (Code 1). Aborting."}'
		exit 1
	}
	REPOHOME="`dirname $t`"
	[[ -z $REPOHOME ]] && {
		echo -n '{"Error":"Getting the gitfs_remotes setting from'
		echo ', $SALT_MASTER_FILE, failed (Code 2). Aborting."}'
		exit 1
	}

	b=`basename $t`
	REPONAME="${b%%.git}"
	[[ -z $REPONAME ]] && {
		echo -n '{"Error":"Getting the gitfs_remotes setting from'
		echo ', $SALT_MASTER_FILE, failed (Code 3). Aborting."}'
		exit 1
	}

	delete_repo_cache

	clone_repo_cache

	version=`get_latest_version_string`

	# New behaviour. Start from version 0.0.0 if there are no versions.
    #[[ -z $version ]] && {
    #    echo -n '{"Error":"Version is empty. Could not find any versions'
    #    echo ' (branches) for '"$BRANCH"'"}'
	#	exit 1
	#}

    [[ -z $version ]] && {
		version="0.0.0"
	}

	new_version=`increment_version_string "$version"`

	create_branch "$new_version"
}

main

# vim:ts=4:sw=4:noet
