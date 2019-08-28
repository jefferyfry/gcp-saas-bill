# Jenkins Support SaaS
This document describes the technical architecture and implementation for a prototype SaaS service that enables a Jenkins Support listing in the Google Cloud Platform (GCP) marketplace. This solution will allow CloudBees to list Jenkins Support as a separate and independent offering in the GCP marketplace so that customers can purchase and transact via the marketplace. The listing would target open source Jenkins users and CloudBees Jenkins Distribution users.

The Jenkins Support SaaS service will act as an agent between the GCP marketplace enablement APIs and the CloudBees support systems. The service will be responsible for processing notifications from the GCP marketplace. These notifications include:

* New subscriptions to Jenkins Support 
* Cancellations
* Upgrades
* Renewals

The service will manage the Jenkins Support subscriptions through a centrally stored customer subscription database which would include:

* Customer contact info
* Subscription period
* Subscription tier

## Architecture
![Architecture](https://user-images.githubusercontent.com/6440106/63820624-c8587700-c8fe-11e9-8493-87e2c761efda.png)

### Components
#### Google Components
* GCP Marketplace - This is the Google Cloud Platform marketplace where the Jenkins Support listing resides and a customer initiates a subscription.
* GCP Marketplace Pub/Sub - A GCP marketplace pub/sub topics notifies the Agent Cloud Function of new subscriptions, cancellations, upgrades and renewals.
* Procurement API - The Procurement API is required to determine the subscription entitlement (User Tier) and status for an account.

#### CloudBees Components
* Agent Cloud Function (CloudBees Developed) - The Agent Cloud Function is triggered by the GCP Marketplace Pub/Sub topics to process new accounts, subscriptions, updates, cancellations and renewals. 
* Subscription Service (CloudBees Developed) - This web app serves the signup page and then approves new accounts and entitlements after receiving account information.
* Subscription Front-end (CloudBees Developed) - Lightweight web interface that provides the signup page.
* Subscription DB (CloudBees Developed) - This is a backup database that stores the current account and subscription data.
* Support Systems - Support systems are the current backend systems such as Zendesk and Salesforce that must be provisioned to enable Jenkins Support services for a customer. The provisioning of these systems is TBD.

## Customer Workflow
![Customer Workflow](https://user-images.githubusercontent.com/6440106/63820521-6435b300-c8fe-11e9-86aa-dfdef195d2e1.png)

## Data Workflow
![Data Workflow](https://user-images.githubusercontent.com/6440106/63820546-7adc0a00-c8fe-11e9-99b1-dcf6eb50192e.png)

1 - Customer subscribes (or makes changes) to Jenkins Support in the marketplace.

2a - Agent Cloud Function receives a notification of a change to a Jenkins Support subscription. This can be a new subscription, cancellation, update or renewal. 

2b - For new accounts, the marketplace directs the customer to a signup page served by the Subscription Manager Web App. 

3 - The Agent Cloud Function queries the procurement API for the entitlement.

4 - The Agent Cloud Function provisions/makes updates via the Subscription Manager Web App via REST API.

5 - Subscription Manager Web App stores account, subscription to subscription database.

6 - Subscription Manager Web App triggers provisioning of support systems (TBD).

7 - Subscription Manager Web App sends successful web page response to customer.

8 - Subscription Manager Web App sends final approval for account and/or entitlement to GCP Procurement API.