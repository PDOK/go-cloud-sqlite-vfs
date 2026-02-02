#!/bin/sh

CBS_VERSION=866fc737fa
SQLITE_VERSION=3510200
SQLITE_YEAR=2026

# Download and install cloud-backed-sqlite
wget -qO ./tmp.zip https://sqlite.org/cloudsqlite/zip/${CBS_VERSION}/cloudsqlite-${CBS_VERSION}.zip
unzip ./tmp.zip
rm tmp.zip

rm -f *.c
rm -f *.h
cp -f ./cloudsqlite-${CBS_VERSION}/src/bcv_int.h ./
cp -f ./cloudsqlite-${CBS_VERSION}/src/bcvutil.c ./
cp -f ./cloudsqlite-${CBS_VERSION}/src/bcvutil.h ./
cp -f ./cloudsqlite-${CBS_VERSION}/src/bcvmodule.c ./
cp -f ./cloudsqlite-${CBS_VERSION}/src/bcvmodule.h ./
cp -f ./cloudsqlite-${CBS_VERSION}/src/blockcachevfs.c ./
cp -f ./cloudsqlite-${CBS_VERSION}/src/blockcachevfs.h ./
cp -f ./cloudsqlite-${CBS_VERSION}/src/simplexml.c ./
cp -f ./cloudsqlite-${CBS_VERSION}/src/simplexml.h ./
cp -f ./cloudsqlite-${CBS_VERSION}/src/bcvencrypt.c ./
cp -f ./cloudsqlite-${CBS_VERSION}/src/bcvencrypt.h ./
cp -f ./cloudsqlite-${CBS_VERSION}/src/bcvlog.c ./
cp -f ./cloudsqlite-${CBS_VERSION}/src/sqlite3.c ./
cp -f ./cloudsqlite-${CBS_VERSION}/src/sqlite3.h ./

rm -rf ./cloudsqlite-${CBS_VERSION}

# Download sqlite and upgrade the default sqlite version shipped with cloud-backed-sqlite
wget -qO ./tmp.zip https://sqlite.org/${SQLITE_YEAR}/sqlite-amalgamation-${SQLITE_VERSION}.zip
unzip ./tmp.zip
rm tmp.zip

cp -f ./sqlite-amalgamation-${SQLITE_VERSION}/sqlite3.c ./
cp -f ./sqlite-amalgamation-${SQLITE_VERSION}/sqlite3.h ./

rm -rf ./sqlite-amalgamation-${SQLITE_VERSION}