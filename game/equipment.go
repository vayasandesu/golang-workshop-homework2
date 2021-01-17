package game

type Equipment struct {
	BonusStatus Status `bson:"bonusstatus" json:"bonusstatus"`
}

func GetEquipmentsStatus(equipments []Equipment) Status {
	bonus := Status{}

	for _, value := range equipments {
		bonus = bonus.Add(value.BonusStatus)
	}

	return bonus
}
