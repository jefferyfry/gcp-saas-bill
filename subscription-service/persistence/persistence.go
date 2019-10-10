package persistence

type DatabaseHandler interface {
	UpsertAccount(*Account) error
	DeleteAccount(string) error
	GetAccount(string) (*Account, error)

	UpsertEntitlement(*Entitlement) error
	DeleteEntitlement(string) error
	GetEntitlement(string) (*Entitlement, error)

	UpsertContact(*Contact) error
	DeleteContact(string) error
	GetContact(string) (*Contact, error)

	QueryEntitlements(filters []string, order string) ([]Entitlement, error)
	QueryAccountEntitlements(accountId string,filters []string, order string) ([]Entitlement, error)
	QueryAccounts(filters []string, order string) ([]Account, error)
	QueryContacts(filters []string, order string) ([]Contact, error)

	Healthz() error
}
