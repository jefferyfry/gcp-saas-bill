package persistence

type DatabaseHandler interface {
	AddAccount(*Account) error
	UpdateAccount(*Account) error
	DeleteAccount(string) error
	GetAccount(string) (*Account, error)

	AddEntitlement(*Entitlement) error
	UpdateEntitlement(*Entitlement) error
	DeleteEntitlement(string) error
	GetEntitlement(string) (*Entitlement, error)
}
