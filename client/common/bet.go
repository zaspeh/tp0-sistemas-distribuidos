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

func (b *Bet) DNI() string {
	return b.config.DNI
}

func (b *Bet) Numero() string {
	return b.config.Numero
}
