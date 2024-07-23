package internal_services

type ServiceConfig struct {
	Addr     string `env:"ADDR" json:"addr"`
	Insecure bool   `env:"insecure"`
}

const (
	TeeTime = "TEE_TIME"
	Catalog = "CATALOG"
	Pricing = "PRICING"
	Payment = "PAYMENT"
	Booking = "BOOKING"
)
