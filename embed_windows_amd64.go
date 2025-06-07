//go:build !ffi_no_embed

package ffi

import (
	_ "embed"
	"os"
	"path/filepath"
)

const (
	moduleName    = "github.com/jupiterrider/ffi"
	libffiVersion = "3.5.0-rc1"
)

//go:embed assets/libffi/windows_amd64/libffi-8.dll
var embeddedLib []byte

//go:embed assets/libffi/LICENSE
var embeddedLicense []byte

func init() {
	if os.Getenv("FFI_NO_EMBED") == "1" {
		return
	}

	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return
	}

	destDir := filepath.Join(userCacheDir, moduleName, "libffi", libffiVersion)
	destLib := filepath.Join(destDir, "libffi-8.dll")
	destLicense := filepath.Join(destDir, "LICENSE")

	if os.MkdirAll(destDir, 0755) == nil {
		if fileInfo, err := os.Stat(destLib); err != nil {
			if os.IsNotExist(err) {
				if os.WriteFile(destLib, embeddedLib, 0755) == nil {
					filename = destLib
				}
			}
		} else {
			if fileInfo != nil && !fileInfo.IsDir() {
				filename = destLib
			}
		}

		if _, err := os.Stat(destLicense); err != nil {
			if os.IsNotExist(err) {
				os.WriteFile(destLicense, embeddedLicense, 0644)
			}
		}
	}
}
