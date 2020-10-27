package monigo

import (
	"os/exec"
	"regexp"
)

type CPU struct {
	User int
	Nice int
	Sys  int
	Idle int

	Busy int
	Work int

	Usage float64
}

var (
	lastCPU = CPU{}
)

func GetCPUUsage() *CPU {
	output, err := exec.Command("cat", "/proc/stat").Output()
	if err != nil {
		LogError("Err `cat /proc/stat` execution: %s", err)
		return nil
	}

	re, _ := regexp.Compile("cpu\\s*([\\d]+) ([\\d]+) ([\\d]+) ([\\d]+) ([\\d]+) ([\\d]+) ([\\d]+) ([\\d]+) ([\\d]+) ([\\d]+)")
	result := re.FindStringSubmatch(string(output))

	var cpu = CPU{
		User: str2int(result[1]),
		Nice: str2int(result[2]),
		Sys:  str2int(result[3]),
		Idle: str2int(result[4]),
	}

	cpu.Busy = cpu.User + cpu.Nice + cpu.Sys
	cpu.Work = cpu.Busy + cpu.Idle

	cpu.Usage = round(100.0*float64(cpu.Busy-lastCPU.Busy)/float64(cpu.Work), 2)

	return &cpu
}
