package cloud_sqlite_vfs

// #cgo LDFLAGS: -lpthread -ldl -lcurl -lssl -lcrypto
// #cgo LDFLAGS: -Wl,--allow-multiple-definition
// #include <stdlib.h>
// #include "blockcachevfs.h"
//
// int csAuthCb(void*, char*, char*, char*, char**);
//
// static char *sqlite3mprintf(char* s) {
//     return sqlite3_mprintf("%s", s);
// }
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

var KEY = ""

//export csAuthCb
func csAuthCb(pCtx *C.void, zStorage *C.char, zAccount *C.char, zContainer *C.char, pzAuthToken **C.char) C.int {
	*pzAuthToken = C.sqlite3mprintf(C.CString(KEY))
	return C.SQLITE_OK
}

func removeCacheDir(cacheDir string) error {
	return os.RemoveAll(cacheDir)
}

func createCacheDir(cacheDir string) error {
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.Mkdir(cacheDir, 0750)
		if err != nil {
			return err
		}
	}

	return nil
}

func Attach(vfsName string, storage string, account string, key string, containerName string, cacheDir string) (*C.sqlite3_bcvfs, error) {
	KEY = key

	err := createCacheDir(cacheDir)
	if err != nil {
		return nil, err
	}

	var pVfs *C.sqlite3_bcvfs
	var zErr *C.char

	rc := C.sqlite3_bcvfs_create(C.CString(cacheDir), C.CString(vfsName), &pVfs, &zErr)
	if rc == C.SQLITE_OK {
		if C.sqlite3_bcvfs_isdaemon(pVfs) == 1 {
			fmt.Println("virtual filesystem is using a daemon")
		} else {
			fmt.Println("virtual filesystem is in daemon less mode")
		}
	} else {
		_ = removeCacheDir(cacheDir)
		return nil, fmt.Errorf("unable to create virtual filesystem with error: %s", C.GoString(zErr))
	}

	if rc == C.SQLITE_OK {
		C.sqlite3_bcvfs_auth_callback(pVfs, nil, (*[0]byte)(unsafe.Pointer(C.csAuthCb)))

		rc = C.sqlite3_bcvfs_attach(pVfs, C.CString(storage), C.CString(account), C.CString(containerName), nil, C.SQLITE_BCV_ATTACH_IFNOT, &zErr)
		if rc != C.SQLITE_OK {
			_ = removeCacheDir(cacheDir)
			return nil, fmt.Errorf("unable to attach virtual filesystem with error: %s", C.GoString(zErr))
		}
	}

	return pVfs, nil
}

func Detach(pVfs *C.sqlite3_bcvfs, cacheDir string) error {
	C.sqlite3_bcvfs_destroy(pVfs)

	return removeCacheDir(cacheDir)
}
