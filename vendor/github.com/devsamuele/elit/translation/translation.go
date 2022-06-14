package translation

type Translation struct {
	It string `json:"it" bson:"it"`
	En string `json:"en" bson:"en"`
	De string `json:"de" bson:"de"`
	Fr string `json:"fr" bson:"fr"`
	Es string `json:"es" bson:"es"`
}

func (t Translation) GreatherThan(length int) bool {
	return !(len(t.It) > length || len(t.Fr) > length || len(t.Es) > length || len(t.En) > length || len(t.De) > length)
}

func (t Translation) IsEmpty() bool {
	return len(t.It) == 0 || len(t.De) == 0 || len(t.Fr) == 0 || len(t.Es) == 0 || len(t.En) == 0
}

// func (t *Translation) MarshalBSON() ([]byte, error) {

// 	type translationAlias Translation
// 	ta := translationAlias(*t)
// 	return bson.Marshal(&ta)
// }

// func (t *Translation) UnmarshalBSON(data []byte) error {

// 	type translationAlias Translation
// 	var ta translationAlias

// 	if err := bson.Unmarshal(data, &ta); err != nil {
// 		return err
// 	}

// 	*t = Translation(ta)

// 	return nil
// }
