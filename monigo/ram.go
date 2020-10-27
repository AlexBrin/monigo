package monigo

import (
	"os/exec"
	"regexp"
)

type RAMValue int64

func (r RAMValue) ToKilobytes() float64 {
	return float64(r) / 1024
}

func (r RAMValue) ToMegabytes() float64 {
	return r.ToKilobytes() / 1024
}

func (r RAMValue) ToGigabytes() float64 {
	return r.ToMegabytes() / 1024
}

type RAM struct {
	Total RAMValue
	Free  RAMValue
	Cached RAMValue
}

func GetRamUsage() *RAM {
	output, err := exec.Command("cat", "/proc/meminfo").Output()
	if err != nil {
		LogError("Err executing `cat /proc/meminfo`: %s", err)
		return &RAM{}
	}

	reTotal, _ := regexp.Compile("MemTotal:\\s*([\\d]+)")
	reFree, _ := regexp.Compile("MemFree:\\s*([\\d]+)")
	reCached, _ := regexp.Compile("Cached:\\s*([\\d]+)")

	total := str2int64(reTotal.FindStringSubmatch(string(output))[1])
	free := str2int64(reFree.FindStringSubmatch(string(output))[1])
	cached := str2int64(reCached.FindStringSubmatch(string(output))[1])

	return &RAM{
		Total: RAMValue(total * 1024),
		Free:  RAMValue((free + cached) * 1024),
		Cached: RAMValue(cached * 1024),
	}
}
