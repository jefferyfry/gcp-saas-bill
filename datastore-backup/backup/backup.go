package backup

import (
	"bytes"
	"errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http/httputil"
)

type DatastoreBackupHandler struct {
	ProjectId   string
	GcsBucket	string
}

func GetDatastoreBackupHandler(projectId string, gcsBucket string) *DatastoreBackupHandler {
	return &DatastoreBackupHandler{
		projectId,
		gcsBucket,
	}
}

func (hdlr *DatastoreBackupHandler) Run() error {
	reqBody := []byte(`{ "entityFilter": { "kinds": [], "namespaceIds": [] },
				"outputUrlPrefix": "`+hdlr.GcsBucket+`" }`)

	datastoreUrl := "https://datastore.googleapis.com/v1/projects/"+hdlr.ProjectId+":export"

	client, clientErr := google.DefaultClient(oauth2.NoContext,"https://www.googleapis.com/auth/cloud-platform https://www.googleapis.com/auth/datastore")

	if clientErr != nil {
		log.Printf("Failed to create oath2 client for the datastore backup %#v \n", clientErr)
		return clientErr
	}

	log.Printf("Requesting datastore backup with request body: %s %s \n", datastoreUrl, reqBody)
	resp, err := client.Post(datastoreUrl,"",bytes.NewBuffer(reqBody))
	if nil != err {
		log.Printf("Failed sending datastore backup request %s %s \n",datastoreUrl, err.Error())
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("Datastore backup request received error response: ",resp.StatusCode)
		responseDump, _ := httputil.DumpResponse(resp, true)
		log.Println(string(responseDump))
		return errors.New("Datastore backup request received error response: "+resp.Status)
	} else {
		log.Printf("Completed datastore backup %s %s %s",datastoreUrl,resp.Status,reqBody)
	}
	return nil
}

