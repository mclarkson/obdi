#!/bin/bash

# Get the version from the spec file (TODO)
if [[ -r VERSION ]]; then
    # OBDI_SEMANTIC_VERSION is set in the VERSION file
    source VERSION
else
    echo "VERSION file does not exist."
    exit 1
fi

V=$OBDI_SEMANTIC_VERSION

[[ -e obdi-${V}.orig ]] && {
    echo "Removing old build directory"
    rm -rf obdi-${V}.orig
}

# Include Golang binaries
export PATH=$PATH:/usr/sbin:/usr/local/go/bin

# Golang requirements
export GOROOT=/usr/local/go
export GOPATH=$PWD/obdi-${V}.orig

# Create the original source tarball
echo "Creating tarball"
tar czf /tmp/obdi_${V}.orig.tar.gz --transform 's#^.#./obdi-'"$V"'.orig#' .
mv /tmp/obdi_${V}.orig.tar.gz .

# Unpack into the now correctly named directory
echo "Unpacking tarball into obdi-${V}.orig"
tar xzf obdi_${V}.orig.tar.gz

# Copy debian directory into build area
echo "Copying 'debian' directory"
cp -a dist/debian obdi-${V}.orig

# Build without signing
echo "Building with 'debuild -i -us -uc'"
cd obdi-${V}.orig
debuild -e PATH -e GOROOT -e GOPATH -e VERSION=$OBDI_SEMANTIC_VERSION -us -uc

#rpmbuild --define "_topdir `pwd`" \
#         --define "BUILD_NUMBER $BUILD_NUMBER" \
#         --define "OBDI_SEMANTIC_VERSION $OBDI_SEMANTIC_VERSION" \
#         -bb SPECS/obdi_rh6.spec

