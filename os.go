package goutil

import (
	"bytes"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// GetPsAuxCount using linux command "ps aux | grep process_name | wc -l" to get a process count from system environment
// It returns the number of requested process.
func GetPsAuxCount(pn string) int {
	ps := exec.Command("ps", "aux")
	grep := exec.Command("grep", pn)
	wc := exec.Command("wc", "-l")
	var out bytes.Buffer
	grep.Stdin, _ = ps.StdoutPipe()
	wc.Stdin, _ = grep.StdoutPipe()
	wc.Stdout = &out

	if err := wc.Start(); err != nil {
		log.Fatal("wc start", err)
	}
	if err := grep.Start(); err != nil {
		log.Fatal("grep start", err)
	}
	if err := ps.Run(); err != nil {
		log.Fatal("ps run", err)
	}
	if err := grep.Wait(); err != nil {
		log.Fatal("grep wait:", err)
	}
	if err := wc.Wait(); err != nil {
		log.Fatal("wc wait", err)
	}

	result, err := strconv.Atoi(strings.TrimSpace(out.String()))
	if err != nil {
		log.Fatal(err)
	}

	// Since ```grep``` itself is also included in the results
	// And we need none nil value for pipe command ```grep``` and ```wc```
	// So we have to drop it (result -1) here
	return result - 1
}

//Check whether/proc/meminfo MemFree/MemTotal is less than $mb MB
func IsFreeMemoryLessThanMB(mb int) bool {
	if stat, err := linuxproc.ReadMemInfo("/proc/meminfo"); err != nil {
		log.Fatal("meminfo read fail")
		return false
	} else {
		mem := ByteSize(stat["MemFree"]) * KB
		return mem < (ByteSize(mb))
	}
}
