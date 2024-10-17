package schema

type Plan struct {
	Name  string `json:"name"`
	Level string `json:"level"`
	ID    int    `json:"id"`
}

func InsertPlans(data map[int]Plan) error {
	db.Exec("DELETE FROM plans")
	for _, plan := range data {
		if err := db.Create(&plan).Error; err != nil {
			return err
		}
	}

	return nil
}
