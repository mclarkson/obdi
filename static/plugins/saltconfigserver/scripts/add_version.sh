#!/bin/bash

[[ -z $1 ]] && {
        echo '{"Error":"Argument missing. Expected a branch name"}'
        exit 1
}

BRANCH="$1"

REPOHOME="/srv/freesat/git-repo"
REPONAME="fdp-mgmt-salt"
REPOTMP="/var/cache/obdi/gitrepo"

clone_repo() {

        # Create the repo cache directory
        mkdir -p $REPOTMP

    # Clone the repo
        cd $REPOTMP
        git clone "$REPOHOME/$REPONAME.git"
}

get_next_version_string() {
        cd $REPOTMP/$REPONAME

        git branch -a | \
                sed -n 's#^\s*remotes/origin/test_\(\)#\1#p' | \
                sort -nt . -k1,1 -k2,2 -k3,3 -k4,4
}

git checkout $BRANCH
git branch test_0.1.10
git push origin test_0.1.10
cd ../fdp-mgmt-salt.git/
git branch -a
/etc/init.d/salt-master restart

