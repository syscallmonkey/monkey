package syscall

import (
	"fmt"
	"io"
	"text/tabwriter"
)

type SyscallCounter struct {
	Counts map[uint64]uint64
}

func NewSyscallCounter() *SyscallCounter {
	sc := SyscallCounter{}
	sc.Counts = make(map[uint64]uint64)
	return &sc
}

func (sc *SyscallCounter) Inc(code uint64) {
	sc.Counts[code] = sc.Counts[code] + 1
}

func (sc *SyscallCounter) Print(out io.Writer) {
	var total uint64
	w := tabwriter.NewWriter(out, 0, 0, 4, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, "SYSCALL (CODE)\tCOUNT\n")
	for k, v := range sc.Counts {
		fmt.Fprintf(w, "%s (%d)\t%d\n", GetSyscallName(k), k, v)
		total += v
	}
	fmt.Fprintf(w, "TOTAL\t%d\n", total)
	w.Flush()
}
