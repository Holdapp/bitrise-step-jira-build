#!/bin/bash
# 
# Perform isolated build of libgit2 dedicated for bitrise ubuntu stacks.
# 
# This scripts builds specified version of libgit2 from source, and creates *.deb package.
#

# print commands & propagate errors
set -x
set -e

VERSION="1.1.0"
DOWNLOAD_URL="https://github.com/libgit2/libgit2/releases/download/v${VERSION}/libgit2-${VERSION}.zip"
PRODUCTS_DIR="/products"
PACKAGE_DIRNAME="libgit2-${VERSION}"
DEB_DIRNAME="libgit2_${VERSION}-1_amd64"

# download sources
wget "${DOWNLOAD_URL}"
unzip "${PACKAGE_DIRNAME}.zip"

# prepare directories
cd ${PACKAGE_DIRNAME}
mkdir -p "${DEB_DIRNAME}/usr/local"
mkdir build/

# build & install
cd build/
cmake ..
DESTDIR="../${DEB_DIRNAME}" make install

# create DEBIAN/control file
cd "../${DEB_DIRNAME}"
mkdir ./DEBIAN/
touch ./DEBIAN/control
echo "Package: libgit2" >>./DEBIAN/control
echo "Version: ${VERSION}"  >>./DEBIAN/control
echo "Architecture: amd64"  >>./DEBIAN/control
echo "Maintainer: Kacper Rączy <gfw.kra@gmail.com>" >>./DEBIAN/control
echo "Description: libgit2 for Bitrise Ubuntu stack" >>./DEBIAN/control

# build package
cd ..
dpkg-deb --build ${DEB_DIRNAME}
cp "${DEB_DIRNAME}.deb" "${PRODUCTS_DIR}/${DEB_DIRNAME}.deb"