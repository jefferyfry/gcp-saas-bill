package persistence

type DatabaseHandler interface {
	UpsertAccount(*Account) error
	DeleteAccount(string) error
	GetAccount(string) (*Account, error)

	UpsertEntitlement(*Entitlement) error
	DeleteEntitlement(string) error
	GetEntitlement(string) (*Entitlement, error)
}
