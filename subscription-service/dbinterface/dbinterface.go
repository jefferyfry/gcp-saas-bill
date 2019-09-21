package dbinterface

import (
	"github.com/cloudbees/cloud-bill-saas/subscription-service/datastoreclient"
	"github.com/cloudbees/cloud-bill-saas/subscription-service/persistence"
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
