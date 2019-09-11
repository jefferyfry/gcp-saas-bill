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

## Additional READMEs
* subscription service [README](/subscription-service/README.md)
* front-end service [README](/subscription-frontend/README.md)
* agent cloud function README

## Architecture
![Architecture](https://user-images.githubusercontent.com/6440106/63956792-00190900-ca3c-11e9-98ab-b84d1fc2f660.png)

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
![Data Workflow](https://user-images.githubusercontent.com/6440106/63956757-e972b200-ca3b-11e9-82d9-51f4b3ab8556.png)

1 - Customer subscribes (or makes changes) to Jenkins Support in the marketplace.

2a - Agent Cloud Function receives a notification of a change to a Jenkins Support subscription. This can be a new subscription, cancellation, update or renewal. 

2b - For new accounts, the marketplace directs the customer to a signup page served by the Subscription Service. 

3 - The Agent Cloud Function queries the procurement API for the entitlement.

4 - The Agent Cloud Function provisions/makes updates via the Subscription Service via REST API.

5 - Subscription Service stores account, subscription to subscription database.

6 - Subscription Service triggers provisioning of backend systems (TBD).

7 - Subscription Service sends successful web page response to customer.

8 - Subscription Service sends final approval for account and/or entitlement to GCP Procurement API.

## Deploying

### Istio with GKE
This application was deployed and tested on GKE clusters version 1.13.7-gke.24 with Istio 1.1.13-gke.0. Kubernetes manifest files are includes for deployment on a GKE cluster with Istio-enabled. For simplicity, set up Istio sidecar auto-injection. Additionally, Istio strict mTLS should be configured.

```
kubectl label namespace NAMESPACE istio-injection=enabled
```

The manifest for the Istio ingress gateway is configured for HTTPS and references certificates. Before applying the manifests, create the cert, key and add the secret.

#### Creating the Cert, Key for the Istio Ingress Gateway HTTPS

```
1. openssl genrsa -des3 -passout pass:x -out server.pass.key 2048

2. openssl rsa -passin pass:x -in server.pass.key -out server.key

3. rm server.pass.key

4. openssl req -new -key server.key -out server.csr
(answer questions)

5. openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.crt

6. kubectl create -n istio-system secret tls istio-ingressgateway-certs --cert=server.crt --key=server.key

```

#### Applying the Istio Manifest
Before applying the manifest update the hosts value to your domain.
```
kubectl apply -f manifests/istio-gateway.yaml
```

#### Applying the Application Manfiest
Before applying the manifest update the environment variables.
```
kubectl apply -f manifests/jenkins-support-saas.yaml
```