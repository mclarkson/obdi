#!/bin/bash

[[ -z $1 ]] && {
	echo '{"Error":"Argument missing. Expected a branch name"}'
	exit 1
}

BRANCH="$1"

REPOHOME="/srv/freesat/git-repo"
REPONAME="fdp-mgmt-salt"
REPOTMP="/var/cache/obdi/gitrepo"

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

    # Clone the repo
	cd "$REPOTMP"
	git clone "$REPOHOME/$REPONAME.git"
}

get_latest_version_string() {
	cd $REPOTMP/$REPONAME

	git branch -a | \
		sed -n 's#^\s*remotes/origin/test_\(\)#\1#p' | \
		sort -nt . -k1,1 -k2,2 -k3,3 -k4,4 | \
		tail -n 1
}

tmp() {
	git checkout $BRANCH
	git branch test_0.1.10
	git push origin test_0.1.10
	cd ../fdp-mgmt-salt.git/
	git branch -a
	/etc/init.d/salt-master restart
}

main() {
	delete_repo_cache
	clone_repo_cache
	get_latest_version_string
}

main
