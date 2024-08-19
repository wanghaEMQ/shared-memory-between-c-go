package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"github.com/hslam/mmap"
	"github.com/hslam/shm"
)

func helper(cmd string) {
	fmt.Printf("example: %v write msg\n", cmd)
	fmt.Printf("         %v read size\n", cmd)
}

func main() {
	name := "test"
	size := 128

	if len(os.Args) < 3 {
		helper(os.Args[0])
		return
	}

	if os.Args[1] == "read" {
		sz, _ := strconv.Atoi(os.Args[2])

		fd, err := shm.Open(name, shm.O_RDONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer shm.Close(fd)
	
		data, err := mmap.Open(fd, 0, size, mmap.READ)
		if err != nil {
			panic(err)
		}
		defer mmap.Munmap(data)
	
		fmt.Printf("%p:%s\n", data, string(data[:sz]))
	} else if os.Args[1] == "write" {
		msg := os.Args[2]

		fd, err := shm.Open(name, shm.O_RDWR|shm.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer shm.Unlink(name)
		defer shm.Close(fd)

		shm.Ftruncate(fd, int64(size))
		data, err := mmap.Open(fd, 0, size, mmap.READ|mmap.WRITE)
		if err != nil {
			panic(err)
		}
		defer mmap.Munmap(data)

		context := []byte(msg)
		copy(data, context)
		fmt.Printf("address %p:%s\n", data, msg)


		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		// Block until a signal is received.
		for {
			<-c
			break
		}
	} else {
		helper(os.Args[0])
		return
	}
}

