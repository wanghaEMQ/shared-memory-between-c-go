CC=gcc
CFLAGS=-g

all: shm2u

shm2u: shm2u.c
	$(CC) -o shm2u ./shm2u.c $(CFLAGS)

.PHONY: clean
clean: shm2u
	rm -f shm2u
