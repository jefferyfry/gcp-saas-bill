package datastoreclient

import (
	"cloud.google.com/go/datastore"
	"context"
	"errors"
	"github.com/cloudbees/cloud-bill-saas/subscription-service/persistence"
	"github.com/jefferyfry/funclog"
	"google.golang.org/api/iterator"
	"strings"
)

const (
	ACCOUNT        = "Account"
	CONTACT    		= "Contact"
	ENTITLEMENT    = "Entitlement"
)

type DatastoreClient struct {
	ProjectId string
}

var (
	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

func NewDatastore(projectId string) (persistence.DatabaseHandler) {
	return &DatastoreClient{
		projectId,
	}
}

func (datastoreClient *DatastoreClient) UpsertAccount(account *persistence.Account) error {
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		LogE.Printf("Failed to create datastore client: %v", err)
		return err
	} else {
		kind := ACCOUNT
		id := account.Id
		key := datastore.NameKey(kind, id, nil)
		_,ptErr := client.Put(ctx, key, account)
		return ptErr
	}
}

func (datastoreClient *DatastoreClient) DeleteAccount(accountId string) error {
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		LogE.Printf("Failed to create datastore client: %v", err)
		return err
	} else {
		kind := ACCOUNT
		key := datastore.NameKey(kind, accountId, nil)
		return client.Delete(ctx, key)
	}
}

func (datastoreClient *DatastoreClient) GetAccount(accountId string) (*persistence.Account, error){
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		LogE.Printf("Failed to create datastore client: %v", err)
		return nil,err
	} else {
		kind := ACCOUNT
		key := datastore.NameKey(kind, accountId, nil)
		account := persistence.Account{}
		gtErr := client.Get(ctx,key, &account)
		return &account, gtErr
	}
}

func (datastoreClient *DatastoreClient) UpsertContact(contact *persistence.Contact) error {
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		LogE.Printf("Failed to create datastore client: %v", err)
		return err
	} else {
		kind := CONTACT
		id := contact.AccountId
		key := datastore.NameKey(kind, id, nil)
		_, ptErr := client.Put(ctx, key, contact)
		return ptErr
	}
}

func (datastoreClient *DatastoreClient) DeleteContact(accountId string) error {
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		LogE.Printf("Failed to create datastore client: %v", err)
		return err
	} else {
		kind := CONTACT
		key := datastore.NameKey(kind, accountId, nil)
		return client.Delete(ctx, key)
	}
}

func (datastoreClient *DatastoreClient) GetContact(accountId string) (*persistence.Contact, error){
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		LogE.Printf("Failed to create datastore client: %v", err)
		return nil,err
	} else {
		kind := CONTACT
		key := datastore.NameKey(kind, accountId, nil)
		contact := persistence.Contact{}
		gtErr := client.Get(ctx, key, &contact)
		return &contact, gtErr
	}
}

func (datastoreClient *DatastoreClient) UpsertEntitlement(entitlement *persistence.Entitlement) error {
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		LogE.Printf("Failed to create datastore client: %v", err)
		return err
	} else {
		kind := ENTITLEMENT
		id := entitlement.Id
		key := datastore.NameKey(kind, id, nil)
		_, ptErr := client.Put(ctx, key, entitlement)
		return ptErr
	}
}

func (datastoreClient *DatastoreClient) DeleteEntitlement(entitlementId string) error {
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		LogE.Printf("Failed to create datastore client: %v", err)
		return err
	} else {
		kind := ENTITLEMENT
		key := datastore.NameKey(kind, entitlementId, nil)
		return client.Delete(ctx, key)
	}
}

func (datastoreClient *DatastoreClient) GetEntitlement(entitlementId string) (*persistence.Entitlement, error){
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		LogE.Printf("Failed to create datastore client: %v", err)
		return nil,err
	} else {
		kind := ENTITLEMENT
		key := datastore.NameKey(kind, entitlementId, nil)
		entitlement := persistence.Entitlement{}
		gtErr := client.Get(ctx, key, &entitlement)
		return &entitlement, gtErr
	}
}

func (datastoreClient *DatastoreClient) QueryEntitlements(filters []string, order string) ([]persistence.Entitlement, error){
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		LogE.Printf("Failed to create datastore client: %v", err)
		return nil,err
	} else {
		q := datastore.NewQuery(ENTITLEMENT)
		if order != "" {
			q = q.Order(order)
		}

		if filters != nil && len(filters)>0 {
			for _, s := range filters {
				if filter := strings.SplitAfter(s, "="); len(filter) > 0 {
					filterStr := filter[0]
					value := filter[1]
					q = q.Filter(filterStr, value)
				}
			}
		}

		t := client.Run(ctx, q)
		var entitlements []persistence.Entitlement
		for {
			entitlement := persistence.Entitlement{}
			_, err := t.Next(&entitlement)
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}
			entitlements = append(entitlements, entitlement)
		}
		return entitlements, nil
	}
}

func (datastoreClient *DatastoreClient) QueryAccountEntitlements(accountId string,filters []string, order string) ([]persistence.Entitlement, error){
	if accountId == "" {
		return nil,errors.New("Must specify account name.")
	}
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		LogE.Printf("Failed to create datastore client: %v", err)
		return nil,err
	} else {
		q := datastore.NewQuery(ENTITLEMENT)
		q.Filter("Account=",accountId)
		if order != "" {
			q = q.Order(order)
		}

		if filters != nil && len(filters)>0 {
			for _, s := range filters {
				if filter := strings.SplitAfter(s, "="); len(filter) > 0 {
					filterStr := filter[0]
					value := filter[1]
					q = q.Filter(filterStr, value)
				}
			}
		}

		t := client.Run(ctx, q)
		var entitlements []persistence.Entitlement
		for {
			entitlement := persistence.Entitlement{}
			_, err := t.Next(&entitlement)
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}
			entitlements = append(entitlements, entitlement)
		}
		return entitlements, nil
	}
}

func (datastoreClient *DatastoreClient) QueryAccounts(filters []string, order string) ([]persistence.Account, error){
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		LogE.Printf("Failed to create datastore client: %v", err)
		return nil,err
	} else {
		q := datastore.NewQuery(ACCOUNT)
		if order != "" {
			q = q.Order(order)
		}

		if filters != nil && len(filters)>0 {
			for _, s := range filters {
				if filter := strings.SplitAfter(s, "="); len(filter) > 0 {
					filterStr := filter[0]
					value := filter[1]
					q = q.Filter(filterStr, value)
				}
			}
		}

		t := client.Run(ctx, q)
		var accounts []persistence.Account
		for {
			account := persistence.Account{}
			_, err := t.Next(&account)
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}
			accounts = append(accounts, account)
		}
		return accounts, nil
	}
}

func (datastoreClient *DatastoreClient) QueryContacts(filters []string, order string) ([]persistence.Contact, error){
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		LogE.Printf("Failed to create datastore client: %v", err)
		return nil,err
	} else {
		q := datastore.NewQuery(CONTACT)
		if order != "" {
			q = q.Order(order)
		}

		if filters != nil && len(filters)>0 {
			for _, s := range filters {
				if filter := strings.SplitAfter(s, "="); len(filter) > 0 {
					filterStr := filter[0]
					value := filter[1]
					q = q.Filter(filterStr, value)
				}
			}
		}

		t := client.Run(ctx, q)
		var contacts []persistence.Contact
		for {
			contact := persistence.Contact{}
			_, err := t.Next(&contact)
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}
			contacts = append(contacts, contact)
		}
		return contacts, nil
	}
}

func (datastoreClient *DatastoreClient) Healthz() error{
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		LogE.Printf("Failed to create datastore client: %v", err)
		return err
	} else {
		q := datastore.NewQuery(ACCOUNT).Limit(1)

		t := client.Run(ctx, q)
		for {
			account := persistence.Account{}
			_, err := t.Next(&account)
			if err == iterator.Done {
				return nil
			}
			if err != nil {
				return err
			}
		}
		return nil
	}
}




