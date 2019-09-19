package dbinterface

import (
	"github.com/cloudbees/cloud-bill-saas/subscription-service/persistence"
	"github.com/cloudbees/cloud-bill-saas/subscription-service/persistence/datastoreclient"
)

type DBTYPE string

const (
	DATASTOREDB DBTYPE = "datastoredb"
)

func NewPersistenceLayer(options DBTYPE, connection string) persistence.DatabaseHandler {
	switch options {
		case DATASTOREDB:
			return datastoreclient.NewDatastore(connection)
	}
	return nil
}
