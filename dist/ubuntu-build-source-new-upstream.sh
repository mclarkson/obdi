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

[[ -e obdi-${V} ]] && {
    echo "Removing old build directory and files"
    rm -rf obdi-${V}
    rm -f obdi_*.diff.gz obdi_*_source.build obdi_*.orig.tar.gz \
          obdi_*.dsc obdi_*_source.changes 
}

# Create the original source tarball
#echo "Creating tarball"
#tar --exclude-vcs -cvzf \
  #/tmp/obdi_${V}.orig.tar.gz --transform 's#^.#./obdi-'"$V"'#' .
#mv /tmp/obdi_${V}.orig.tar.gz .

cp ../obdi_${V}.orig.tar.gz .

# Unpack into the now correctly named directory
echo "Unpacking tarball into obdi-${V}.orig"
tar xzf obdi_${V}.orig.tar.gz

# Copy debian directory into build area
echo "Copying 'debian' directory"
cp -a dist/debian obdi-${V}

# Build without signing
echo "Building with 'debuild -i -us -uc'"
cd obdi-${V}

# Build binary
#debuild -us -uc

# Build source - Include orig.tar.gz
debuild -S -sa

# Build source - Exclude orig.tar.gz
#debuild -S -sd

