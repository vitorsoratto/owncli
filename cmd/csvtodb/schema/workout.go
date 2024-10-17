package schema

type Workout struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	PlanID int    `json:"plan_id"`
}

func InsertWorkouts(data map[int]Workout) error {
	db.Exec("DELETE FROM workouts")
	for _, workout := range data {
		if err := db.Create(&workout).Error; err != nil {
			return err
		}
	}

	return nil
}
