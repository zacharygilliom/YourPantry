package recipe

type Next struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

type Link struct {
	Det Next `json:"next"`
}

type Ingredient struct {
	Text         string  `json:"text"`
	Weight       float64 `json:"weight"`
	FoodCategory string  `json:"foodCategory"`
	FoodID       string  `json:"foodId"`
	Image        string  `json:"image"`
}

type RecipeData struct {
	Uri             string       `json:"uri"`
	Label           string       `json:"label"`
	Image           string       `json:"image"`
	Source          string       `json:"source"`
	Url             string       `json:"url"`
	ShareAs         string       `json:"shareAs"`
	Yield           float64      `json:"yield"`
	DietLabels      []string     `json:"dietLabels"`
	HealthLabels    []string     `json:"healthLabels"`
	Cautions        []string     `json:"cautions"`
	IngredientLines []string     `json:"ingredientLines"`
	Ingredients     []Ingredient `json:"ingredients"`
	Calories        float64      `json:"calories"`
	TotalWeight     float64      `json:"totalWeight"`
	TotalTime       float64      `json:"totalTime"`
	CuisineType     []string     `json:"cuisineType"`
	MealType        []string     `json:"mealType"`
	DishType        []string     `json:"dishType"`
}
type Hit struct {
	Recipe RecipeData `json:"recipe"`
}
type Recipes struct {
	From  int  `json:"from"`
	To    int  `json:"to"`
	Count int  `json:"count"`
	Links Link `json:"_links"`
	Hits  Hit  `json:"hits"`
}
