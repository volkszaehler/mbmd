go tool cgo -godefs types_windows.go | gofmt > ztypes_windows.go
mksyscall_windows syscall_windows.go | gofmt > zsyscall_windows.go
