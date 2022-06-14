package model

type Translation struct {
	It string `json:"it" bson:"it"`
	En string `json:"en" bson:"en"`
	De string `json:"de" bson:"de"`
	Fr string `json:"fr" bson:"fr"`
	Es string `json:"es" bson:"es"`
}

func (t Translation) IsValidLength(length int) bool {
	return !(len(t.It) > length || len(t.Fr) > length || len(t.Es) > length || len(t.En) > length || len(t.De) > length)
}

func (t Translation) IsEmpty() bool {
	return len(t.It) == 0 || len(t.De) == 0 || len(t.Fr) == 0 || len(t.Es) == 0 || len(t.En) == 0
}
