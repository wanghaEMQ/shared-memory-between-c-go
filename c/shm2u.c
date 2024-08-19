#include <stdio.h>
#include <stdlib.h>
#include <sys/mman.h>
#include <errno.h>
#include <fcntl.h>
#include <string.h>
#include <unistd.h>
#include <signal.h>

static int running = 1;
void
inthandler(int sig)
{
	running = 0;
}

int
nsh_open_shmem(char *name)
{
	int fd = shm_open(name, O_CREAT | O_RDWR, 0644);
	if (fd == -1) {
		printf("open %d failed errno%d\n", fd, errno);
		return -1;
	}

	return fd;
}

int
nsh_close_shmem(char *name)
{
	int rv = shm_unlink(name);
	if (rv != 0) {
		printf("close shmem %s failed errno%d\n", name, errno);
		return -1;
	}

	return 0;
}

int
nsh_grow_shmem(int fd, size_t size)
{
	if (ftruncate(fd, size) == -1) {
		printf("ftruncate %d failed errno%d\n", fd, errno);
		return -1;
	}

	return 0;
}

void *
nsh_open_mmap(int fd, size_t size)
{
	int protection = PROT_READ | PROT_WRITE;

	// The buffer will be shared (meaning other processes can access it),
	// but anonymous (meaning third-party processes cannot obtain an
	// address for it), so only this process and its children will be able
	// to use it:
	// int visibility = MAP_SHARED | MAP_ANONYMOUS;
	int visibility = MAP_SHARED;

	return mmap(NULL, size, protection, visibility, fd, 0);
}

int
nsh_close_mmap(void *mem, size_t size)
{
	return munmap(mem, size);
}

void
helper(char *cmd)
{
	printf("example: %s write msg\n", cmd);
	printf("         %s read size\n", cmd);
}

int
main(int argc, char *argv[])
{
	char *bin = argv[0];
	if (argc < 3) {
		helper(bin);
		return 0;
	}
	char *cmd = argv[1];
	size_t size = 128;
	char *name = "test";

	if (0 == strcmp(cmd, "write")) {
		signal(SIGINT, inthandler);
		char * msg = argv[2];

		int fd = nsh_open_shmem(name);
		if (fd == -1) {
			printf("failed to open shmem errno%d\n", errno);
			return 0;
		}

		int rv = nsh_grow_shmem(fd, size);
		if (rv != 0) {
			printf("failed to grow shmem errno%d\n", errno);
			return 0;
		}

		void *shmem = nsh_open_mmap(fd, size);
		if (shmem == MAP_FAILED) {
			printf("failed to get mmap errno%d\n", errno);
			return 0;
		}

		memcpy(shmem, msg, strlen(msg) + 1);
		printf("address %p:%s\n", shmem, msg);
		while (running) {
			sleep(1);
		}

		rv = nsh_close_mmap(shmem, size);
		if (rv != 0) {
			printf("failed to close mmap errno%d\n", errno);
			return 0;
		}

		rv = nsh_close_shmem(name);
		if (rv != 0) {
			printf("failed to unlink shmem%d errno%d\n", rv, errno);
			return 0;
		}
	} else if (0 == strcmp(cmd, "read")) {
		int size = atoi(argv[2]);

		int fd = nsh_open_shmem(name);
		if (fd == -1) {
			printf("failed to open shmem errno%d\n", errno);
			return 0;
		}

		void *shmem = nsh_open_mmap(fd, size);
		if (shmem == MAP_FAILED) {
			printf("failed to get mmap errno%d\n", errno);
			return 0;
		}

		printf("%p:%.*s\n", shmem, size, (char *)shmem);
	} else {
		helper(bin);
	}
}
