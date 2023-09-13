package main

import (
	"fmt"
	"github.com/PDOK/go-cloud-sqlite-vfs"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	STORAGE         = "azure?emulator=127.0.0.1:10000"
	ACCOUNT         = "devstoreaccount1"
	KEY             = "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw=="
	CONTAINER_NAME  = "example"
	VFS_NAME        = "myvfs"
	CACHE_DIR       = "./tmp"
	SQLITE_FILENAME = "chinook.db"
)

type Genre struct {
	Id   string `db:"GenreId"`
	Name string `db:"Name"`
}

func main() {
	vfs, err := cloud_sqlite_vfs.NewVFS(VFS_NAME, STORAGE, ACCOUNT, KEY, CONTAINER_NAME, CACHE_DIR)
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := sqlx.Open("sqlite3", "/"+CONTAINER_NAME+"/"+SQLITE_FILENAME+"?vfs="+VFS_NAME)
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := db.Queryx("SELECT * FROM genres")
	if err != nil {
		fmt.Println(err)
	} else {
		genre := Genre{}
		for rows.Next() {
			err := rows.StructScan(&genre)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%#v\n", genre)
		}
	}

	err = vfs.Close()
	if err != nil {
		fmt.Println(err)
	}
}
