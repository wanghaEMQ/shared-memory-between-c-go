package shm2u

import (
	"github.com/hslam/sem"
	"github.com/hslam/ftok"
)

func SemNew(pathname string, projectid uint8) int {
	key, err := ftok.Ftok(pathname, projectid)
	if err != nil {
		panic(err)
	}
	nsems := 1
	semid, err := sem.Get(key, nsems, 0666)
	if err != nil {
		// Create and init
		semid, err = sem.Get(key, nsems, sem.IPC_CREAT|sem.IPC_EXCL|0666)
		if err != nil {
			panic(err)
		}
		for semnum := 0; semnum < nsems; semnum++ {
			_, err := sem.SetValue(semid, semnum, 1)
			if err != nil {
				panic(err)
			}
		}
	}
	return semid
}

func SemFree(semid int) {
	sem.Remove(semid)
}

func SemValue(semid int) int {
	semnum := 0
	cnt, err := sem.GetValue(semid, semnum)
	if err != nil {
		panic(err)
	}
	return cnt
}

func SemP(semid int) (bool, error) {
	semnum := 0
	return sem.P(semid, semnum, sem.SEM_UNDO)
}

func SemV(semid int) (bool, error) {
	semnum := 0
	return sem.V(semid, semnum, sem.SEM_UNDO)
}

