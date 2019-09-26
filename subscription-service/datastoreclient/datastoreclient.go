package datastoreclient

import (
	"cloud.google.com/go/datastore"
	"context"
	"errors"
	"github.com/cloudbees/cloud-bill-saas/subscription-service/persistence"
	"google.golang.org/api/iterator"
	"log"
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

func NewDatastore(projectId string) (persistence.DatabaseHandler) {
	return &DatastoreClient{
		projectId,
	}
}

func (datastoreClient *DatastoreClient) UpsertAccount(account *persistence.Account) error {
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return err
	} else {
		kind := ACCOUNT
		name := account.Name
		key := datastore.NameKey(kind, name, nil)
		_,ptErr := client.Put(ctx, key, account)
		return ptErr
	}
}

func (datastoreClient *DatastoreClient) DeleteAccount(name string) error {
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return err
	} else {
		kind := ACCOUNT
		key := datastore.NameKey(kind, name, nil)

		return client.Delete(ctx, key)
	}
}

func (datastoreClient *DatastoreClient) GetAccount(name string) (*persistence.Account, error){
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return nil,err
	} else {
		kind := ACCOUNT
		key := datastore.NameKey(kind, name, nil)
		account := persistence.Account{}
		gtErr := client.Get(ctx,key, account)
		return &account, gtErr
	}
}

func (datastoreClient *DatastoreClient) UpsertContact(contact *persistence.Contact) error {
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return err
	} else {
		accountKey := datastore.NameKey(ACCOUNT, contact.AccountName, nil)
		kind := CONTACT
		name := contact.AccountName
		key := datastore.NameKey(kind, name, accountKey)
		_, ptErr := client.Put(ctx, key, contact)
		return ptErr
	}
}

func (datastoreClient *DatastoreClient) DeleteContact(accountName string) error {
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return err
	} else {
		kind := CONTACT
		key := datastore.NameKey(kind, accountName, nil)
		return client.Delete(ctx, key)
	}
}

func (datastoreClient *DatastoreClient) GetContact(accountName string) (*persistence.Contact, error){
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return nil,err
	} else {
		kind := CONTACT
		key := datastore.NameKey(kind, accountName, nil)
		contact := persistence.Contact{}
		gtErr := client.Get(ctx, key, contact)
		return &contact, gtErr
	}
}

func (datastoreClient *DatastoreClient) UpsertEntitlement(entitlement *persistence.Entitlement) error {
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return err
	} else {
		accountKey := datastore.NameKey(ACCOUNT, entitlement.Account, nil)
		kind := ENTITLEMENT
		name := entitlement.Name
		key := datastore.NameKey(kind, name, accountKey)
		_, ptErr := client.Put(ctx, key, entitlement)
		return ptErr
	}
}

func (datastoreClient *DatastoreClient) DeleteEntitlement(entitlementName string) error {
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return err
	} else {
		kind := ENTITLEMENT
		key := datastore.NameKey(kind, entitlementName, nil)
		return client.Delete(ctx, key)
	}
}

func (datastoreClient *DatastoreClient) GetEntitlement(entitlementName string) (*persistence.Entitlement, error){
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return nil,err
	} else {
		kind := ENTITLEMENT
		key := datastore.NameKey(kind, entitlementName, nil)
		entitlement := persistence.Entitlement{}
		gtErr := client.Get(ctx, key, entitlement)
		return &entitlement, gtErr
	}
}

func (datastoreClient *DatastoreClient) QueryEntitlements(filters []string, order string) ([]persistence.Entitlement, error){
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		log.Printf("Failed to create datastore client: %v", err)
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

func (datastoreClient *DatastoreClient) QueryAccountEntitlements(accountName string,filters []string, order string) ([]persistence.Entitlement, error){
	if accountName == "" {
		return nil,errors.New("Must specify account name.")
	}
	ctx := context.Background()

	if client, err := datastore.NewClient(ctx, datastoreClient.ProjectId); err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return nil,err
	} else {
		q := datastore.NewQuery(ENTITLEMENT)
		q.Filter("Account=",accountName)
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
		log.Printf("Failed to create datastore client: %v", err)
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




