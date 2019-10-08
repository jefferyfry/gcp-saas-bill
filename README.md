# Cloud Bill SaaS
This document describes the technical architecture and implementation for a prototype SaaS service that enables a marketplace transactions in the Google Cloud Platform (GCP) marketplace. 

The Cloud Bill SaaS service will act as an agent between the GCP marketplace enablement APIs and the CloudBees support systems. The service will be responsible for processing notifications from the GCP marketplace. These notifications include:

* New subscriptions 
* Cancellations
* Upgrades
* Renewals

The service will manage the subscriptions through a centrally stored customer subscription database which would include:

* Customer contact info
* Subscription period
* Subscription tier

## Additional READMEs
* front-end service [README](/frontend-service/README.md)
* subscription service [README](/subscription-service/README.md)
* pubsub service [README](/pubsub-service/README.md)

## Architecture
![Architecture](https://user-images.githubusercontent.com/6440106/64708575-ae27a880-d469-11e9-8006-e947c950cc91.png)

### Components
#### Google/Auth0 Components
* GCP Marketplace - This is the Google Cloud Platform marketplace where the listing resides and a customer initiates a subscription.
* GCP Marketplace Pub/Sub - A GCP marketplace pub/sub topics notifies the Agent Cloud Function of new subscriptions, cancellations, upgrades and renewals.
* Procurement API - The Procurement API is required to determine the subscription entitlement (User Tier) and status for an account.
* Google Kubernetes Engine with Istio - CloudBees components are run on GKE with Istio. Istio is used for some routing and security policies.
* Auth0 Application/API - Auth0 is used to authenticate and gather user profile data.

#### CloudBees Components
* Agent Cloud Function (CloudBees Developed) - The Agent Cloud Function is triggered by the GCP Marketplace Pub/Sub topics to process new accounts, subscriptions, updates, cancellations and renewals. 
* Subscription Service (CloudBees Developed) - This web app serves the signup page and then approves new accounts and entitlements after receiving account information.
* Subscription Front-end (CloudBees Developed) - Lightweight web interface that provides the signup page.
* Subscription DB (CloudBees Developed) - This is a backup database that stores the current account and subscription data.
* Support Systems - Support systems are the current backend systems such as Zendesk and Salesforce that must be provisioned to enable Jenkins Support services for a customer. The provisioning of these systems is TBD.

## Customer Workflow
![Customer Workflow](https://user-images.githubusercontent.com/6440106/63820521-6435b300-c8fe-11e9-86aa-dfdef195d2e1.png)

## Data Workflow
![Data Workflow](https://user-images.githubusercontent.com/6440106/64708366-5d17b480-d469-11e9-8137-2977472a1515.png)

1 - Customer subscribes (or makes changes) to the listing in the marketplace.

2a - Agent Cloud Function receives a notification of a change to a subscription. This can be a new subscription, cancellation, update or renewal. 

2b - For new accounts, the marketplace directs the customer to a signup page served by the Subscription Service. 

3 - The Agent Cloud Function queries the procurement API for the entitlement.

4 - The Agent Cloud Function provisions/makes updates via the Subscription Service via REST API.

5 - Subscription Service stores account, subscription to subscription database.

6 - Subscription Service triggers provisioning of backend systems (TBD).

7 - Subscription Service sends successful web page response to customer.

8 - Subscription Service sends final approval for account and/or entitlement to GCP Procurement API.

## Operations

### Deploying 

#### Istio with GKE
This application was deployed and tested on GKE clusters version 1.13.7-gke.24 with Istio 1.1.13-gke.0. Kubernetes manifest files are includes for deployment on a GKE cluster with Istio-enabled. For simplicity, set up Istio sidecar auto-injection. Additionally, Istio strict mTLS should be configured.

```
kubectl label namespace NAMESPACE istio-injection=enabled
```

The manifest for the Istio ingress gateway is configured for HTTPS and references certificates. Before applying the manifests, create the cert, key and add the secret.

##### Creating the Cert, Key for the Istio Ingress Gateway HTTPS

```
1. openssl genrsa -des3 -passout pass:x -out server.pass.key 2048

2. openssl rsa -passin pass:x -in server.pass.key -out server.key

3. rm server.pass.key

4. openssl req -new -key server.key -out server.csr
(answer questions)

5. openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.crt

6. kubectl create -n istio-system secret tls istio-ingressgateway-certs --cert=server.crt --key=server.key

```

##### Applying the Istio Manifest
Before applying the manifest update the hosts value to your domain.
```
kubectl apply -f manifests/istio-gateway.yaml
```

##### Applying the Application Manfiest
Before applying the manifest update the environment variables or provide configuration files.
```
kubectl apply -f manifests/cloud-bill-saas.yaml
```

##### GCP Service Accounts
The following roles are required:
* PubSub Editor - Used to access marketplace PubSub events.
* Cloud Datastore Owner - Used for the Cloud Datastore subscription DB.
* Cloud Import Export Admin - Used to export from Cloud Datastore.
* Cloud Commerce API (assigned by GCP Marketplace team) - Allows access to the Cloud Commerce API
* Billing Account Administrator (NOT FOR PRODUCTION) - Allows the reset of test accounts.
It is recommended that the roles be used assigned to a common service account. Then the service account file can be shared and mounted for all the services.

### Monitoring

#### Datadog
DataDog can be configured using the Kubernetes agent. First configure RBAC permissions for DataDog.

```
kubectl create -f "https://raw.githubusercontent.com/DataDog/datadog-agent/master/Dockerfiles/manifests/rbac/clusterrole.yaml"
kubectl create -f "https://raw.githubusercontent.com/DataDog/datadog-agent/master/Dockerfiles/manifests/rbac/serviceaccount.yaml"
kubectl create -f "https://raw.githubusercontent.com/DataDog/datadog-agent/master/Dockerfiles/manifests/rbac/clusterrolebinding.yaml"
```
Then create the Kubernetes secret for your API key.

```
kubectl create secret generic datadog-secret --from-literal api-key="<api-key>"
```

Then deploy the datadog-agent using the manifest.

```
kubectl apply -f manifest/datadog-agent.yaml
```

#### Sentry
Sentry is configured for all services. 

### Upgrades
It is recommended that you use K8s rolling update for service images. This can be executed in a single command:

```
kubectl set image deployments/<deployment-name> <container>=image

eg.
kubectl set image deployments/frontend-service frontend-service=gcr.io/cloud-bill-dev/frontend-service:2
```
If the upgrade includes configuration changes, apply those configuration changes first.

You can monitor the rolling update using:

```
kubectl rollout status deployments/<deployment-name>
```

### Security

### Troubleshooting