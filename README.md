# Obdi

## Screenshot

![](images/obdi-0.1.2.png?raw=true)

More pics here: http://blogger.smorg.co.uk/2015/01/obdi-in-pictures.html

## What is it?

Obdi is an extendable graphical user interface for running scripts on
remote servers.  The most basic installation provides administrative tools to
manage users, data centres and environments, and includes a job manager.

## Status

Development and testing phase. Not ready for use.

## Features

* All GUI items in the screen shot are implemented using plugins.
* Plugins can be REST end-points - they can be used for automation.
* The saltconfigserver plugin contains a number of REST end-points, including an External Node Classifier (ENC).

## Architecture

* Google Go for the back end REST interface.
* Angular JS and bootstrap for the front end.
* Sqlite3 used for data storage.
* Obdi simply runs scripts on remote servers, so it can be used for lots of different tasks.

## Install

For Red Hat based systems RPMs can be built using:
```
cd
yum install rpmdevtools gcc
git clone https://github.com/mclarkson/obdi.git
cd obdi
BUILD_NUMBER=1 ./dist/jenkins-build.sh
```

## Todo

Moved to GitHub Issues
