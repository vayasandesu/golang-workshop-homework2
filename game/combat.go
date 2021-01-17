package game

func Attack(attacker Character, defender Character) {
	atk := attacker.GetFinalStatus().Attack
	def := defender.GetFinalStatus().Defense

	dmg := atk - def

	defender.PrimaryStatus.Health -= dmg
}
