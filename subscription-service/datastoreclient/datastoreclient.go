package datastoreclient

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/cloudbees/cloud-bill-saas/subscription-service/persistence"
	"log"
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

	client, err := datastore.NewClient(ctx, datastoreClient.ProjectId)
	if err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return err
	}

	kind := ACCOUNT
	name := account.Name
	key := datastore.NameKey(kind, name, nil)

	_,ptErr := client.Put(ctx, key, account)

	return ptErr
}

func (datastoreClient *DatastoreClient) DeleteAccount(name string) error {
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, datastoreClient.ProjectId)
	if err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return err
	}

	kind := ACCOUNT
	key := datastore.NameKey(kind, name, nil)

	return client.Delete(ctx,key)
}

func (datastoreClient *DatastoreClient) GetAccount(name string) (*persistence.Account, error){
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, datastoreClient.ProjectId)
	if err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return nil,err
	}

	kind := ACCOUNT
	key := datastore.NameKey(kind, name, nil)

	account := persistence.Account{}

	gtErr := client.Get(ctx,key, account)

	return &account, gtErr
}

func (datastoreClient *DatastoreClient) UpsertContact(contact *persistence.Contact) error {
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, datastoreClient.ProjectId)
	if err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return err
	}
	accountKey := datastore.NameKey(ACCOUNT, contact.AccountName, nil)

	kind := CONTACT
	name := contact.AccountName
	key := datastore.NameKey(kind, name, accountKey)

	_,ptErr := client.Put(ctx, key, contact)

	return ptErr
}

func (datastoreClient *DatastoreClient) DeleteContact(accountName string) error {
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, datastoreClient.ProjectId)
	if err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return err
	}

	kind := CONTACT
	key := datastore.NameKey(kind, accountName, nil)

	return client.Delete(ctx,key)
}

func (datastoreClient *DatastoreClient) GetContact(accountName string) (*persistence.Contact, error){
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, datastoreClient.ProjectId)
	if err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return nil,err
	}

	kind := CONTACT
	key := datastore.NameKey(kind, accountName, nil)

	contact := persistence.Contact{}

	gtErr := client.Get(ctx,key, contact)

	return &contact, gtErr
}

func (datastoreClient *DatastoreClient) UpsertEntitlement(entitlement *persistence.Entitlement) error {
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, datastoreClient.ProjectId)
	if err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return err
	}
	accountKey := datastore.NameKey(ACCOUNT, entitlement.Account, nil)

	kind := ENTITLEMENT
	name := entitlement.Name
	key := datastore.NameKey(kind, name, accountKey)

	_,ptErr := client.Put(ctx, key, entitlement)

	return ptErr
}

func (datastoreClient *DatastoreClient) DeleteEntitlement(entitlementName string) error {
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, datastoreClient.ProjectId)
	if err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return err
	}

	kind := ENTITLEMENT
	key := datastore.NameKey(kind, entitlementName, nil)

	return client.Delete(ctx,key)
}

func (datastoreClient *DatastoreClient) GetEntitlement(entitlementName string) (*persistence.Entitlement, error){
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, datastoreClient.ProjectId)
	if err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return nil,err
	}

	kind := ENTITLEMENT
	key := datastore.NameKey(kind, entitlementName, nil)

	entitlement := persistence.Entitlement{}

	gtErr := client.Get(ctx,key, entitlement)

	return &entitlement, gtErr
}



