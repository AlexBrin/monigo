package monigo

import (
	"os/exec"
	"regexp"
	"time"
)

type Speed int64

func (s Speed) ToKilobytesPerSecond() float64 {
	return float64(s) / 1024.0
}

func (s Speed) ToMegabytesPerSecond() float64 {
	return s.ToKilobytesPerSecond() / 1024.0
}

var (
	devices   = map[string]int64{}
	last      = int64(0)
	lastTotal = int64(0)
)

func GetNetworkUsage() Speed {
	output, err := exec.Command("cat", "/proc/net/dev").Output()
	if err != nil {
		LogError("Error executing `cat /proc/net/dev`: %s", err)
		return Speed(0)
	}

	re, _ := regexp.Compile("([a-z0-9]+):\\s*([\\d]+)\\s*([\\d]+)\\s*([\\d]+)\\s*([\\d]+)\\s*([\\d]+)\\s*([\\d]+)\\s*([\\d]+)\\s*([\\d]+)\\s*([\\d]+)")
	result := re.FindAllStringSubmatch(string(output), 10)

	total := int64(0)
	for _, device := range result {
		devices[device[1]] = str2int64(device[2]) + str2int64(device[10])
		total += devices[device[1]]
	}

	seconds := time.Now().Unix() - last
	bytesPerSeconds := (total - lastTotal) / seconds

	last = time.Now().Unix()
	lastTotal = total

	return Speed(bytesPerSeconds)
}
