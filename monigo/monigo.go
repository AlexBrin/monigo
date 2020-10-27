package monigo

import "time"

func Start() {
	go CPUWorker()
	go RAMWorker()
	go NetworkWorker()
	go CleanOldData()

	StartWebInterface()
}

func CleanOldData() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			old := time.Now().Unix() - (Config.Resource.OldDataLifeTime * 60 * 60 * 24)
			DB.Where("created_at < ?", old).Delete(&History{})
		}
	}
}

func CPUWorker() {
	ticker := time.NewTicker(time.Second * time.Duration(Config.Resource.CPU.Frequency))
	for {
		select {
		case <-ticker.C:
			cpu := GetCPUUsage()

			DB.Save(&History{
				Type: "cpu",
				Val:  cpu.Usage,
			})
		}
	}
}

func RAMWorker() {
	ticker := time.NewTicker(time.Second * time.Duration(Config.Resource.CPU.Frequency))
	for {
		select {
		case <-ticker.C:
			ram := GetRamUsage()

			DB.Save(&History{
				Type: "ram",
				Val:  float64(ram.Total),
			})
		}
	}
}

func NetworkWorker() {
	ticker := time.NewTicker(time.Second * time.Duration(Config.Resource.CPU.Frequency))
	for {
		select {
		case <-ticker.C:
			network := GetNetworkUsage()

			DB.Save(&History{
				Type: "network",
				Val:  float64(network),
			})
		}
	}
}


