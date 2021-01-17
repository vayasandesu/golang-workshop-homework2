package game

type Status struct {
	Health  int `bson:"health" json:"health"`
	Attack  int `bson:"attack" json:"attack"`
	Defense int `bson:"defense" json:"defense"`
}

func (a Status) Add(b Status) Status {
	return Status{
		Attack:  a.Attack + b.Attack,
		Defense: a.Defense + b.Defense,
		Health:  a.Health + b.Health,
	}
}
