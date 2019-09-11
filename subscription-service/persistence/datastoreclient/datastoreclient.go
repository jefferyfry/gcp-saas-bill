package datastoreclient

import (
	"github.com/cloudbees/jenkins-support-saas/subscription-service/persistence"
	"context"
	"cloud.google.com/go/datastore"
	"log"
)

const (
	ACCOUNT        = "Account"
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

func (datastoreClient *DatastoreClient) DeleteEntitlement(name string) error {
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, datastoreClient.ProjectId)
	if err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return err
	}

	kind := ENTITLEMENT
	key := datastore.NameKey(kind, name, nil)

	return client.Delete(ctx,key)
}

func (datastoreClient *DatastoreClient) GetEntitlement(name string) (*persistence.Entitlement, error){
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, datastoreClient.ProjectId)
	if err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return nil,err
	}

	kind := ENTITLEMENT
	key := datastore.NameKey(kind, name, nil)

	entitlement := persistence.Entitlement{}

	gtErr := client.Get(ctx,key, entitlement)

	return &entitlement, gtErr
}

