package models

type Ingredient struct {
	ID   int    `gorm:"primaryKey;serial;unique;not null;autoIncrement" json:"id"`
	Name string `gorm:"unique; not null" json:"name"`
}

// type IngredientDTO struct {
// 	RecipeID     int    `json:"recipeid"`
// 	IngredientID int    `json:"ingredientid"`
// 	Amount       int    `json:"amount"`
// 	Unit         string `json:"unit"`
// }

// // ConvertToDTO converts a single m.Recipe_Ingredient to a single m.IngredientDTO
// func (i Recipe_Ingredient) ConvertToDTO() IngredientDTO {
// 	return IngredientDTO{
// 		RecipeID:     i.RecipeID,
// 		IngredientID: i.IngredientID,
// 		Amount:       i.Amount,
// 		Unit:         i.Unit,
// 	}
// }

// // ConvertAllToDTO converts a slice of m.Recipe_Ingredient to a slice of m.IngredientDTO objcets
// func (i Recipe_Ingredient) ConvertAllToDTO(ingredients []Recipe_Ingredient) []IngredientDTO {
// 	var data []IngredientDTO

// 	for _, ingredient := range ingredients {
// 		data = append(data, ingredient.ConvertToDTO())
// 	}

// 	return data
// }
