# go-cloud-sqlite-vfs

# Description

This project wraps the [Cloud Backed SQLite](https://sqlite.org/cloudsqlite/doc/trunk/www/index.wiki) (CSB)
solution into a golang package. The project uses [The SQLite OS Interface or "VFS"](https://www.sqlite.org/vfs.html)
concept to create a VFS which is backed by either Azure Blob Storage or Google Cloud Storage. The VFS can be used 
with every SQLite golang package as long as it supports setting a custom VFS name.

Below are some examples about how to use this package, for further information about the workings of this package
please read the [documentation](https://sqlite.org/cloudsqlite/doc/trunk/www/index.wiki) of the CBS project.

# Installation

This package can be installed with the `go get` command:

```bash
go get github.com/PDOK/go-cloud-sqlite-vfs 
```

**go-cloud-sqlite-vfs is cgo package**. If you want to build your app using go-cloud-sqlite-vfs , you need gcc.

# Usage

```go
package main

import (
	"fmt"
	"github.com/PDOK/go-cloud-sqlite-vfs"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	STORAGE         = "azure?emulator=127.0.0.1:10000"
	ACCOUNT         = "<ACCOUNT>"
	KEY             = "<KEY>"
	CONTAINER_NAME  = "<CONTAINER_NAME>"
	VFS_NAME        = "myvfs"
	CACHE_DIR       = "./tmp"
	SQLITE_FILENAME = "<SQLITE_FILENAME>"
)

func main() {
	vfs, err := cloud_sqlite_vfs.Attach(VFS_NAME, STORAGE, ACCOUNT, KEY, CONTAINER_NAME, CACHE_DIR)
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := sqlx.Open("sqlite3", "/"+CONTAINER_NAME+"/"+SQLITE_FILENAME+"?vfs="+VFS_NAME)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(db)

	err = cloud_sqlite_vfs.Detach(vfs, CACHE_DIR)
	if err != nil {
		fmt.Println(err)
	}
}
```

# Example project

1. Build `blockcachevfsd` cli. For instuctions see the [CBS website](https://sqlite.org/cloudsqlite/doc/trunk/www/index.wiki)
2. Start Azurite
    ```bash
    docker run -p 10000:10000 mcr.microsoft.com/azure-storage/azurite azurite-blob --blobHost 0.0.0.0
    ```
3. Create an `account.txt` file with the following content (don't forget the linebreak):
    ```
   -module azure?emulator=127.0.0.1:10000&sas=0
   
   ```
4. Create a container in Azurite
    ```bash
    ./blockcachevfsd create -f account.txt example
    ```
5. Upload the SQLite database to the newly created container on Azurite
    ```bash
   ./blockcachevfsd upload -f account.txt -container example ./chinook.db chinook.db
   ```
6. Run the example application
   ```bash
   go run ./
   ```
   
# Dev

Because the C code of the [CBS project](https://sqlite.org/cloudsqlite/dir?ci=tip) needs to be included in 
the package it can be updated with the `download-c-code.sh` script located in the root of this project.
