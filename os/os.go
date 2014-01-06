package os

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

	psRes := ExecuteCmd(ps, "")
	grepRes := ExecuteCmd(grep, psRes)
	if grepRes == "" {
		return 0
	}

	// pass the `grep` result to `wc` input, and run it
	wcRes := ExecuteCmd(wc, grepRes)

	result, err := strconv.Atoi(strings.TrimSpace(wcRes))
	if err != nil {
		log.Fatal("strconv.Atoi error: ", err)
	}

	// Since ```grep``` itself is also included in the results
	// And we need none nil value for pipe command ```grep``` and ```wc```
	// So we have to drop it (result -1) here
	return result - 1
}

// Pass the `pipe` as `cmd` input, and run it
// It returns the result of cmd
func ExecuteCmd(cmd *exec.Cmd, pipe string) string {
	var out bytes.Buffer

	if pipe != "" {
		cmd.Stdin = strings.NewReader(pipe)
	}
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Println("ooops:", err.Error())
		return ""
	}

	return out.String()
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
