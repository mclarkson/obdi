#!/bin/bash

# Get the version from the spec file (TODO)
VERSION="0.1.0"

# Get the list of top level files before we make changes
files=`echo * | sed 's/rpmbuild//'`

mkdir -p rpmbuild/BUILD/
mkdir -p rpmbuild/RPMS/
mkdir -p rpmbuild/SOURCES/
mkdir -p rpmbuild/SPECS/
mkdir -p rpmbuild/SRPMS/

cp dist/redhat6/obdi_rh6.spec rpmbuild/SPECS

# Copy 
mkdir -p rpmbuild/SOURCES/obdi-$VERSION
cp -a $files rpmbuild/SOURCES/obdi-$VERSION

# Compress
pushd rpmbuild/SOURCES
tar cvzf obdi-$VERSION.tar.gz obdi-$VERSION

popd
rpmbuild --define "_topdir `pwd`" --define "BUILD_NUMBER $BUILD_NUMBER" -bb SPECS/obdi_rh6.spec

