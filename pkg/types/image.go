package types

type Image struct {
	Name      string  `json:"name"`
	Location  string  `json:"location"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func (i Image) Fields() []interface{} {
	return []interface{}{
		i.Name,
		i.Location,
		i.Longitude,
		i.Latitude,
	}
}
