package service

import (
	"context"
	"sort"
	"twitterjokeposting/repo"
	"twitterjokeposting/util"

	"github.com/samber/lo"
)

func getJokesToPublishToday() int {
	return util.RandomNumberGenerator(10)
}

func ScheduleJokeForToday() {
	// Create background context
	ctx := context.Background()

	// Clear Schedule Job database for today
	repo.ClearAllScheduledJobs()

	// Get How many jokes we want to publish today
	jokesToPublish := getJokesToPublishToday()

	// Get that much random times today
	randomTimesToPublishJoke := []int{}
	for iterator := 0; iterator < jokesToPublish; iterator++ {
		randomTime := util.RandomNumberGenerator(24)
		randomTimesToPublishJoke = append(randomTimesToPublishJoke, randomTime)
	}
	sort.Ints(randomTimesToPublishJoke)
	randomTimesToPublishJoke = lo.Uniq[int](randomTimesToPublishJoke)

	// Insert that much documents into database today
	repo.CreateScheduledJobs(randomTimesToPublishJoke, ctx)
}

func GetAllScheduledJokes() []repo.ScheduleJobSchema {
	return repo.GetScheuledJobs()
}
