#!/bin/bash

# Get the version from the spec file (TODO)
if [[ -r VERSION ]]; then
    # OBDI_SEMANTIC_VERSION is set in the VERSION file
    source VERSION
else
    echo "VERSION file does not exist."
    exit 1
fi

# Get the list of top level files before we make changes
files=`echo * | sed 's/rpmbuild//'`

mkdir -p rpmbuild/BUILD/
mkdir -p rpmbuild/RPMS/
mkdir -p rpmbuild/SOURCES/
mkdir -p rpmbuild/SPECS/
mkdir -p rpmbuild/SRPMS/

cp dist/redhat6/obdi_rh6.spec rpmbuild/SPECS

# Copy 
mkdir -p rpmbuild/SOURCES/obdi-$OBDI_SEMANTIC_VERSION
cp -a $files rpmbuild/SOURCES/obdi-$OBDI_SEMANTIC_VERSION

# Compress
cd rpmbuild/SOURCES
tar cvzf obdi-$OBDI_SEMANTIC_VERSION obdi-$OBDI_SEMANTIC_VERSION

cd ..
rpmbuild --define "_topdir `pwd`" \
         --define "BUILD_NUMBER $BUILD_NUMBER" \
         --define "OBDI_SEMANTIC_VERSION $OBDI_SEMANTIC_VERSION" \
         -bb SPECS/obdi_rh6.spec

