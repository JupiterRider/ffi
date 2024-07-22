//go:build ((freebsd || linux) && arm64) || (windows && (amd64 || arm64))

package ffi

const (
	DefaultAbi Abi = 1
)
