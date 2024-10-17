package schema

type Exercise struct {
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Repetitions     int     `json:"repetitions"`
	AnimationMale   string  `json:"animation_male"`
	AnimationFemale string  `json:"animation_female"`
	Instructions    *string `json:"instructions"`
	ID              int     `json:"id"`
	WorkoutID       int     `json:"workout_id"`
	Animated        int     `json:"animated"`
}

func InsertExercises(data []Exercise) error {
	db.Exec("DELETE FROM exercises")
	for _, exercise := range data {
		if err := db.Create(&exercise).Error; err != nil {
			return err
		}
	}

	return nil
}
