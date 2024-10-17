package schema

import (
	"encoding/csv"
	"os"
	"strconv"

)

type SheetData struct {
	Plan            string `json:"plan"`
	Workout         string `json:"workout"`
	Exercise        string `json:"exercise"`
	Type            string `json:"type"`
	Repetitions     int    `json:"repetitions"`
	Instructions    string `json:"instructions"`
	ID              int    `json:"id"`
	WorkoutID       int    `json:"workout_id"`
	PlanID          int    `json:"plan_id"`
	AnimationMale   string `json:"animation_male"`
	AnimationFemale string `json:"animation_female"`
	Level           string `json:"level"`
	Animated        int    `json:"animated"`
}

func parseSheetData(data [][]string) (map[int]Plan, map[int]Workout, []Exercise) {
	plans := map[int]Plan{}
	workouts := map[int]Workout{}
	exercises := []Exercise{}

	for _, row := range data[1:] {
		repetitions, _ := strconv.Atoi(row[4])
		id, _ := strconv.Atoi(row[6])
		workoutID, _ := strconv.Atoi(row[7])
		planID, _ := strconv.Atoi(row[8])

		animated := 1
		if row[2] == "PLANK" {
			animated = 0
		}

		sheet := SheetData{
			Plan:            row[0],
			Workout:         row[1],
			Exercise:        row[2],
			Type:            row[3],
			Repetitions:     repetitions,
			Instructions:    row[5],
			ID:              id,
			WorkoutID:       workoutID,
			PlanID:          planID,
			AnimationMale:   row[9],
			AnimationFemale: row[10],
			Level:           row[11],
			Animated:        animated,
		}

		if _, exists := plans[sheet.PlanID]; !exists {
			plans[planID] = Plan{
				Name:  sheet.Plan,
				Level: sheet.Level,
				ID:    sheet.PlanID,
			}
		}

		if _, exists := workouts[sheet.WorkoutID]; !exists {
			workouts[workoutID] = Workout{
				Name:   sheet.Workout,
				PlanID: sheet.PlanID,
				ID:     sheet.WorkoutID,
			}
		}

		exercises = append(exercises, Exercise{
			Name:            sheet.Exercise,
			ID:              sheet.ID,
			Type:            sheet.Type,
			Repetitions:     sheet.Repetitions,
			AnimationMale:   sheet.AnimationMale,
			AnimationFemale: sheet.AnimationFemale,
			Instructions:    &sheet.Instructions,
			WorkoutID:       sheet.WorkoutID,
			Animated:        sheet.Animated,
		})
	}
	return plans, workouts, exercises
}

func ReadCSV(path string) (map[int]Plan, map[int]Workout, []Exercise, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, nil, err
	}

	plans, workouts, exercises := parseSheetData(data)

	return plans, workouts, exercises, nil
}
