package common

// Información de la apuesta
type Bet struct {
	config BetConfig
}

type BetConfig struct {
	Nombre     string
	Apellido   string
	DNI        string
	Nacimiento string
	Numero     string
}

func NewBet(config BetConfig) *Bet {
	bet := &Bet{
		config: config,
	}
	return bet
}
