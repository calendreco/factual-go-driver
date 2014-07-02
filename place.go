package factual

type Place struct{
	Tel string `json:"tel"`
	Name string `json:"name"`
	Email string `json:"email"`
	Website string `json:"website"`
	Hours Hours `json:"hours"`
	HoursDisplay string `json:"hours_display"`
	FactualId string `json:"factual_id"`
	Address string `json:"address"`
	Neighborhood []string `json:"neighborhood"`
	Region string `json:"region"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Postcode string `json:"postcode"`
}