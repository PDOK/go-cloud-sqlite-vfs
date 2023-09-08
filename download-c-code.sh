#!/bin/sh

VERSION=8535bc08c4

wget -qO ./tmp.zip https://sqlite.org/cloudsqlite/zip/${VERSION}/cloudsqlite-${VERSION}.zip
unzip ./tmp.zip
rm tmp.zip

rm -f *.c
rm -f *.h
cp -f ./cloudsqlite-${VERSION}/src/bcv_int.h ./
cp -f ./cloudsqlite-${VERSION}/src/bcvutil.c ./
cp -f ./cloudsqlite-${VERSION}/src/bcvutil.h ./
cp -f ./cloudsqlite-${VERSION}/src/bcvmodule.c ./
cp -f ./cloudsqlite-${VERSION}/src/bcvmodule.h ./
cp -f ./cloudsqlite-${VERSION}/src/blockcachevfs.c ./
cp -f ./cloudsqlite-${VERSION}/src/blockcachevfs.h ./
cp -f ./cloudsqlite-${VERSION}/src/simplexml.c ./
cp -f ./cloudsqlite-${VERSION}/src/simplexml.h ./
cp -f ./cloudsqlite-${VERSION}/src/bcvencrypt.c ./
cp -f ./cloudsqlite-${VERSION}/src/bcvencrypt.h ./
cp -f ./cloudsqlite-${VERSION}/src/bcvlog.c ./
cp -f ./cloudsqlite-${VERSION}/src/sqlite3.c ./
cp -f ./cloudsqlite-${VERSION}/src/sqlite3.h ./

rm -rf ./cloudsqlite-${VERSION}