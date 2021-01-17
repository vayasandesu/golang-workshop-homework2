package game

type Character struct {
	Name          string      `bson:"name" json:"name"`
	PrimaryStatus Status      `bson:"primarystatus" json:"primarystatus"`
	Equipments    []Equipment `bson:"equipments" json:"equipments"`
}

func (character Character) GetFinalStatus() Status {
	base := character.PrimaryStatus
	bonus := GetEquipmentsStatus(character.Equipments)

	return base.Add(bonus)
}

func NewCharacter(name string) Character {
	return Character{
		Name: name,
		PrimaryStatus: Status{
			Attack:  1,
			Defense: 1,
			Health:  10,
		},
		Equipments: []Equipment{},
	}
}
