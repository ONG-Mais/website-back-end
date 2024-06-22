package entities

type State struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type City struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	State State  `json:"state_city"`
}

type SubjectInterest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Implements the voluntary struct
type Voluntary struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Email           string          `json:"email"`
	Phone           int             `json:"phone"`
	State           State           `json:"state"`
	City            City            `json:"city"`
	SubjectInterest SubjectInterest `json:"Subject_Interest"`
}
