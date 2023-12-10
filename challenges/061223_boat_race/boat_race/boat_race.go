package boat_race

import (
	"context"
	"strconv"
	"strings"

	"github.com/olivermonaco/2023_advent_of_code/kit"
)

type boatRace struct {
	alottedTime    int
	recordDistance int
}

type raceApproachResult struct {
	buttonHoldTime int
	distanceSailed int
}

func CalculatePartOne(ctx context.Context, input []string) int {
	boatRaces := createBoatRaces(input)
	var productWinningPossibilities *int
	for _, boatRace := range boatRaces {
		distanceToRaceApproaches := boatRace.possibleButtonHoldTimes()
		winningPossiblities := calcWinningPossibilities(boatRace, distanceToRaceApproaches)
		if productWinningPossibilities == nil {
			productWinningPossibilities = kit.Ptr(winningPossiblities)
			continue
		}
		productWinningPossibilities = kit.Ptr(
			kit.Deref(productWinningPossibilities) * winningPossiblities,
		)
	}
	return *productWinningPossibilities
}

func CalculatePartTwo(ctx context.Context, input []string) int {
	return 0
}

func createBoatRaces(lines []string) []boatRace {

	var times, distances []string
	for _, line := range lines {
		_, timesStr, foundTimes := strings.Cut(line, "Time:")
		if foundTimes {
			times = strings.Fields(timesStr)
		}
		_, distancesStr, foundDistances := strings.Cut(line, "Distance:")
		if foundDistances {
			distances = strings.Fields(distancesStr)
		}
	}
	if len(times) != len(distances) {
		panic(append(times, distances...))
	}

	var boatRaces []boatRace
	for i := 0; i < len(times); i++ {
		timeInt, err := strconv.Atoi(times[i])
		if err != nil {
			panic(timeInt)
		}
		distanceInt, err := strconv.Atoi(distances[i])
		if err != nil {
			panic(timeInt)
		}
		boatRaces = append(boatRaces, boatRace{
			recordDistance: distanceInt,
			alottedTime:    timeInt,
		})
	}
	return boatRaces

}

func (bR boatRace) possibleButtonHoldTimes() map[int][]raceApproachResult {

	halfTime := bR.alottedTime

	raceApproaches := make(map[int][]raceApproachResult)

	for i := 0; i < halfTime+1; i++ {
		buttonHoldTime := i
		runTime := bR.alottedTime - buttonHoldTime
		distanceTraveled := buttonHoldTime * runTime
		distanceRaceApproaches := []raceApproachResult{
			{
				buttonHoldTime: buttonHoldTime,
				distanceSailed: distanceTraveled,
			},
		}
		if existingRaceApproaches, ok := raceApproaches[distanceTraveled]; ok {
			distanceRaceApproaches = append(
				existingRaceApproaches,
				distanceRaceApproaches...,
			)
		}
		raceApproaches[distanceTraveled] = distanceRaceApproaches
	}
	return raceApproaches
}

func calcWinningPossibilities(bR boatRace, distanceToRaceApproaches map[int][]raceApproachResult) int {

	var winningPossibilities int
	for distanceTraveled, raceApproaches := range distanceToRaceApproaches {
		if distanceTraveled <= bR.recordDistance {
			continue
		}
		winningPossibilities += len(raceApproaches)
	}
	return winningPossibilities
}
