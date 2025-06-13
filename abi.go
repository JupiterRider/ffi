//go:build ((freebsd || linux || darwin) && arm64) || (windows && amd64)

package ffi

const (
	DefaultAbi Abi = 1
)
