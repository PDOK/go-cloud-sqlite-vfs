package cloud_sqlite_vfs

// #cgo LDFLAGS: -lpthread -ldl -lcurl -lssl -lcrypto
// #cgo !darwin LDFLAGS: -Wl,--allow-multiple-definition
// #cgo CFLAGS: -DSQLITE_ENABLE_RTREE=1
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

const (
	// Argument sets the maximum size of the cachefile in bytes. This option is only available for daemonless mode VFSs.
	sqliteBcvCacheSize = 1
	// The argument enables (non-zero) or disables (zero) verbose libcurl logging.
	sqliteBcvCurlVerbose = 4
)

var KEY = ""

// VFS represent a SQLite virtual file systems backed by (Azure/Google) cloud object storage.
type VFS struct {
	bcvfs    *C.sqlite3_bcvfs
	cacheDir string

	// true when running Cloud-Backed SQLite in daemonless mode, false when running as a daemon.
	Daemonless bool
}

//export csAuthCb
func csAuthCb(pCtx *C.void, zStorage *C.char, zAccount *C.char, zContainer *C.char, pzAuthToken **C.char) C.int {
	cKey := C.CString(KEY)
	*pzAuthToken = C.sqlite3mprintf(cKey)
	defer C.free(unsafe.Pointer(cKey))

	return C.SQLITE_OK
}

func removeCacheDir(cacheDir string) error {
	return os.RemoveAll(cacheDir)
}

func createCacheDir(cacheDir string) error {
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err = os.Mkdir(cacheDir, 0750)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewVFS build a new VFS. Note: when provided cacheSize is 0 the default CBS cache size of 1Gb will be used
func NewVFS(vfsName string, storage string, account string, key string,
	containerName string, cacheDir string, cacheSize int64, curlVerbose bool) (VFS, error) {

	KEY = key
	vfs := &VFS{}

	err := createCacheDir(cacheDir)
	if err != nil {
		return *vfs, err
	}

	var pVfs *C.sqlite3_bcvfs
	var zErr *C.char

	cCacheDir := C.CString(cacheDir)
	cVFSName := C.CString(vfsName)

	rc := C.sqlite3_bcvfs_create(cCacheDir, cVFSName, &pVfs, &zErr)

	defer C.free(unsafe.Pointer(cCacheDir))
	defer C.free(unsafe.Pointer(cVFSName))

	if rc != C.SQLITE_OK {
		_ = removeCacheDir(cacheDir)
		return *vfs, fmt.Errorf("unable to create virtual filesystem with error: %s", C.GoString(zErr))
	}

	if curlVerbose {
		C.sqlite3_bcvfs_config(pVfs, sqliteBcvCurlVerbose, 1)
	}
	if C.sqlite3_bcvfs_isdaemon(pVfs) == 1 {
		vfs.Daemonless = false
	} else {
		vfs.Daemonless = true

		if cacheSize > 0 {
			// cache only works in daemonless mode
			C.sqlite3_bcvfs_config(pVfs, sqliteBcvCacheSize, C.longlong(cacheSize))
		}
	}

	C.sqlite3_bcvfs_auth_callback(pVfs, nil, (*[0]byte)(unsafe.Pointer(C.csAuthCb)))

	cStorage := C.CString(storage)
	cAccount := C.CString(account)
	cContainerName := C.CString(containerName)

	rc = C.sqlite3_bcvfs_attach(pVfs, cStorage, cAccount, cContainerName, nil, C.SQLITE_BCV_ATTACH_IFNOT, &zErr)

	defer C.free(unsafe.Pointer(cStorage))
	defer C.free(unsafe.Pointer(cAccount))
	defer C.free(unsafe.Pointer(cContainerName))

	if rc != C.SQLITE_OK {
		_ = removeCacheDir(cacheDir)
		return *vfs, fmt.Errorf("unable to attach virtual filesystem with error: %s", C.GoString(zErr))
	}

	vfs.bcvfs = pVfs
	vfs.cacheDir = cacheDir
	return *vfs, nil
}

func (vfs VFS) Close() error {
	C.sqlite3_bcvfs_destroy(vfs.bcvfs)
	return removeCacheDir(vfs.cacheDir)
}
