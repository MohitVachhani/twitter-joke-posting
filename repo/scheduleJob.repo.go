package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ScheduleJobSchema struct {
	ID              string
	ScheduledAtHour int
	CreatedAt       string
}

var ScheduleJobDatabase = []ScheduleJobSchema{}

func ClearAllScheduledJobs() {
	ScheduleJobDatabase = []ScheduleJobSchema{}
}

func CreateScheduledJobs(scheduledAtHours []int, ctx context.Context) {

	toInsertScheduleJobs := []ScheduleJobSchema{}
	for _, scheduledAtHour := range scheduledAtHours {
		scheduledJobID := uuid.New().String()
		currentTimeStamp := time.Now().Format(time.RFC3339)

		toCreateScheduleJob := ScheduleJobSchema{
			ID:              scheduledJobID,
			ScheduledAtHour: scheduledAtHour,
			CreatedAt:       currentTimeStamp,
		}

		toInsertScheduleJobs = append(toInsertScheduleJobs, toCreateScheduleJob)
	}

	ScheduleJobDatabase = append(ScheduleJobDatabase, toInsertScheduleJobs...)
}

func GetScheuledJobs() []ScheduleJobSchema {
	return ScheduleJobDatabase
}
