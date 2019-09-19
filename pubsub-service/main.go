package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudbees/jenkins-support-saas/marketplace-agent-service/config"
	"log"
)

type PubSubMsg struct {
	EventId     	string	`json:"name"`
	EventType   	string	`json:"account"`
	Entitlement		Entitlement `json:entitlement,omitempty`
	Account			Account `json:account,omitempty`
}

type Entitlement struct {
	Name     			string	`json:"name,omitempty"`
	Account   			string	`json:"account,omitempty"`
	Provider    		string	`json:"provider,omitempty"`
	Product  			string	`json:"product,omitempty"`
	Plan     	  		string	`json:"plan,omitempty"`
	NewPendingPlan 	  	string	`json:"newPendingPlan,omitempty"`
	State    	  		int64	`json:"state,omitempty"`
	UpdateTime    	  	string	`json:"updateTime,omitempty"`
	CreateTime    	  	string	`json:"createTime,omitempty"`
	UsageReportingId    string	`json:"usageReportingId,omitempty"`
	MessageToUser    	string	`json:"messageToUser,omitempty"`
}

type Account struct {
	Name  			string     	`json:"name,omitempty"`
	UpdateTime   	string    	`json:"updateTime,omitempty"`
	CreateTime      string    	`json:"createTime,omitempty"`
	Provider     	string		`json:"provider,omitempty"`
	State 	 		string      `json:"state,omitempty"`
	Approvals    	string     	`json:"approvals,omitempty"`
}

func main() {
	fmt.Println("Starting Jenkins Support SaaS Agent Service...")
	config, err := config.GetConfiguration()

	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, config.GcpProjectId)
	if err != nil {
		log.Fatalf("Error creating pubsub client %s: %v", config.PubSubSubscription, err)
		log.Fatal(err)
	}

	subscription := client.Subscription(config.PubSubSubscription)

	exists, err := subscription.Exists(ctx)
	if err != nil {
		log.Fatalf("Error checking for subscription: %v", err)
	}

	if !exists {
		log.Fatalf("GCP marketplace pubsub subscription does not exist %s: %v", config.PubSubSubscription, err)
	}

	errRcv := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		pubSubMsg := PubSubMsg{}
		if err := json.Unmarshal(msg.Data, &pubSubMsg); err != nil {
			fmt.Println("could not decode message data: %#v", msg)
			msg.Nack()
			return
		}

		fmt.Println("Received msg %v", pubSubMsg)
		switch pubSubMsg.EventType {
			case "ACCOUNT_ACTIVE":
				err := updateAccount(&pubSubMsg.Account)
				if err != nil {
					fmt.Println("Unable to update account %v due to error %v",pubSubMsg.Entitlement,err)
				}
			case "ENTITLEMENT_CREATION_REQUESTED":
				err := postEntitlementResponseToProcurementAPI(pubSubMsg.Entitlement.Name+":approve")
				if err != nil {
					fmt.Println("Unable to approve entitlement plan change %v due to error %v",pubSubMsg.Entitlement,err)
				}
			case "ENTITLEMENT_PLAN_CHANGE_REQUESTED":
				err := postEntitlementResponseToProcurementAPI(pubSubMsg.Entitlement.Name+":approvePlanChange")
				if err != nil {
					fmt.Println("Unable to approve entitlement plan change %v due to error %v",pubSubMsg.Entitlement,err)
				}
			case "ENTITLEMENT_PLAN_CHANGED":
				err := updateEntitlement(&pubSubMsg.Entitlement)
				if err != nil {
					fmt.Println("Unable to update entitlement plan %v due to error %v",pubSubMsg.Entitlement,err)
				}
			case "ENTITLEMENT_PENDING_CANCELLATION":
				fallthrough
			case "ENTITLEMENT_CANCELLED":
				err := updateEntitlement(&pubSubMsg.Entitlement)
				if err != nil {
					fmt.Println("Unable to cancel entitlement %v due to error %v",pubSubMsg.Entitlement,err)
				}
			case "ENTITLEMENT_DELETED":
				err := deleteEntitlement(&pubSubMsg.Entitlement)
				if err != nil {
					fmt.Println("Unable to delete entitlement %v due to error %v",pubSubMsg.Entitlement,err)
				}
			case "ACCOUNT_DELETED":
				err := deleteAccount(&pubSubMsg.Account)
				if err != nil {
					fmt.Println("Unable to delete account %v due to error %v",pubSubMsg.Entitlement,err)
				}
			default:
				fmt.Println("Unknown pubsub event type %v",pubSubMsg)
		}

		/*if err := update(id); err != nil {
			log.Printf("[ID %d] could not update: %v", id, err)
			msg.Nack()
			return
		}*/

		msg.Nack()
	})
	if errRcv != nil {
		log.Fatal(errRcv)
	}
}

func postEntitlementResponseToProcurementAPI(path string) error {
	return nil;
}

func updateEntitlement(entitlement *Entitlement) error {
	return nil;
}

func deleteEntitlement(entitlement *Entitlement) error {
	return nil;
}

func updateAccount(account *Account) error {
	return nil;
}

func deleteAccount(account *Account) error {
	return nil;
}
