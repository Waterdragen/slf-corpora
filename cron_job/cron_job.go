package cron_job

import "time"

type CronJob struct {
	Task func()
	Dur  time.Duration
}

func NewDailyCronJob(task func()) CronJob {
	return CronJob{Task: task, Dur: 86400 * time.Second}
}

func (cj CronJob) Start() {
	go func() {
		ticker := time.NewTicker(cj.Dur)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cj.Task()
			}
		}
	}()
}
