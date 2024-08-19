# SHM2U

This project show how to shared memory IPC between C and Golang (without cgo) program.

## Start

### C Program

```
cd c
make
./shm2u
```

### Golang Program

```
cd golang
go build
./shm2u-go
```

### Test

```
./shm2u-go write abc
./shm2u read 3

or

./shm2u write bcd
./shm2u-go read 3
```
