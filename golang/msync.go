package shm2u

type MSync struct {
	sz      uint
	semid   int
	sempath string

	start  uint
	end    uint
}

type MChunk struct {
	data []byte
	sz   uint
}

func New(sempath string, semprojectid uint8, sz uint) *MSync {
	ms := MSync{}
	ms.sz = sz
	ms.sempath = sempath

	ms.semid = SemNew(sempath, semprojectid)

	return &ms
}

func (ms *MSync) Free() {
	SemFree(ms.semid)
}

func (ms *MSync) GetChunk(sz uint) *MChunk {
	data := make([]byte, sz)
	// TODO get it from mmap
	return &MChunk{
		data: data,
		sz:   sz,
	}
}

func (ms *MSync) Push(mc *MChunk) {
}

func (ms *MSync) Pop() *MChunk {
	sz := 2
	data := make([]byte, sz)
	// TODO get it from mmap
	return &MChunk{
		data: data,
		sz:   uint(sz),
	}
}

func (mc *MChunk) Data() *[]byte {
	return &mc.data
}

func (mc *MChunk) Rele() {
}
