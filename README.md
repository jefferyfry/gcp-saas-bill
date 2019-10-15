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
* datastore backup cron job [README](/datastore-backup/README.md)

## Architecture
![Architecture](https://user-images.githubusercontent.com/6440106/66532740-5a8a9800-eac5-11e9-8501-7a008e8d6f9c.png)

### Components
#### Google/Auth0 Components
* GCP Marketplace - This is the Google Cloud Platform marketplace where the listing resides and a customer initiates a subscription.
* GCP Marketplace Pub/Sub - A GCP marketplace pub/sub topics notifies the Agent Cloud Function of new subscriptions, cancellations, upgrades and renewals.
* Procurement API - The Procurement API is required to determine the subscription entitlement (User Tier) and status for an account.
* Google Kubernetes Engine with Istio - CloudBees components are run on GKE with Istio. Istio is used for some routing and security policies.
* Auth0 Application/API - Auth0 is used to authenticate and gather user profile data.

#### CloudBees Components
* PubSub Service (CloudBees Developed) - The PubSub Service is triggered by the GCP Marketplace Pub/Sub topics to process new accounts, subscriptions, updates, cancellations and renewals. 
* Subscription Service (CloudBees Developed) - This web app serves the signup page and then approves new accounts and entitlements after receiving account information.
* Frontend Service(CloudBees Developed) - Lightweight web interface that provides the signup page.
* Subscription DB (CloudBees Developed) - This is a backup database that stores the current account and subscription data.
* Support Systems - Support systems are the current backend systems such as Zendesk and Salesforce that must be provisioned to enable Jenkins Support services for a customer. The provisioning of these systems is TBD.
* Datastore Backup Cron Job - Daily executing datastore backup.

## Customer Workflow
![Customer Workflow](https://user-images.githubusercontent.com/6440106/66532891-e00e4800-eac5-11e9-8db3-4a2656066d51.png)

## Data Workflow
![Data Workflow](https://user-images.githubusercontent.com/6440106/66532860-c1a84c80-eac5-11e9-8559-f055a89e66c8.png)

1 - Customer subscribes (or makes changes) to Jenkins Support in the marketplace.

2a - PubSub Service receives a notification of a change to a Jenkins Support subscription. This can be a new subscription, cancellation, update or renewal. 

2b - For new accounts, the marketplace directs the customer to a signup page served by the Frontend Service. Auth0 is used to authenticate and gather account information. This information is stored via the Subscription Service.

3 - The PubSub Service queries the procurement API for the entitlement or account.

4 - The PubSub Service provisions/makes updates via the Subscription Service via REST API.

5 - Subscription Service stores account, subscription to subscription database.

6 - Subscription Service triggers provisioning of backend systems (TBD).

7 - Frontend Service and PubSub Service sends final approval for account and/or entitlement to GCP Procurement API.

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

#### External Access
Access to the application is only allowed to the Frontend-Service and is controlled by Istio. See the [Istio Gateway manifest](/manifests/istio-gateway.yaml) that is applied to configure this. Four pages are hosted by the frontend-service:

* [signup.html](https://github.com/cloudbees/cloud-bill-saas/tree/master/frontend-service/templates/signup.html) - Initial page to direct customer to Auth0/Google sign in. The customer is sent to this page from marketplace.
* [confirmProd.html](https://github.com/cloudbees/cloud-bill-saas/tree/master/frontend-service/templates/confirmProd.html) - Auth0/Google callback page to confirm account information for VM and K8s products.
* [confirmSaas.html](https://github.com/cloudbees/cloud-bill-saas/tree/master/frontend-service/templates/confirmSaas.html) - Auth0/Google callback page to confirm account information for Saas products.
* [finish.html](https://github.com/cloudbees/cloud-bill-saas/tree/master/frontend-service/templates/finish.html) - Final page to confirm account creation and notify customer of next steps.

Additionally, the frontend-service redirects to the CloudBees Auth0 service for account creation and authentication.

#### Firewall Rules for External Access
(All other ports are blocked for external access)
| Port  | Source | Description |
|-------|--------|-------------|
| 80    | 0.0.0.0/0 | Redirects to 443.       |
| 443   | 0.0.0.0/0 | HTTPS for serving pages.|

#### Cloud Management Access
The development version of this application is hosted in the GCP project cloud-bill-dev/cloud-bill-dev. IAM membership and management console access for this project can bee seen [here](https://console.cloud.google.com/iam-admin/iam?project=cloud-bill-dev).
The production version of this application is hosted in the GCP project gcp-marketplace-solutions/cje-marketplace-dev. IAM membership and management console access for this project can bee seen [here](https://console.cloud.google.com/iam-admin/iam?project=cje-marketplace-dev).

#### Secrets and Service Accounts
A common GCP service account is used across all services with the following roles:
* PubSub Editor - Used to access marketplace PubSub events.
* Cloud Datastore Owner - Used for the Cloud Datastore subscription DB.
* Cloud Import Export Admin - Used to export from Cloud Datastore.
* Cloud Commerce API (assigned by GCP Marketplace team) - Allows access to the Cloud Commerce API
The service account JSON file is installed as a Kubernetes secret for access by all the services:

```
kubectl create secret generic gcp-service-account --from-file gcp-service-account.json
```

Additionally, due to sensitive CloudBees partner metadata stored in the services configuration files, these are also stores as Kubernetes secrets. See each of the services' READMEs for more.

#### Images and Scanning
Four images are used in the application. Images are hosted in GCR and automatically scanned. These are the GCR locations for development and production:

* pubsub-service - The PubSub Service is triggered by the GCP Marketplace Pub/Sub topics to process new accounts, subscriptions, updates, cancellations and renewals. 
* subscription-service - This web app serves the signup page and then approves new accounts and entitlements after receiving account information.
* frontend-service - Lightweight web interface that provides the signup page.
* datastore-backup - Daily executing datastore backup.

Development (cloud-bill-dev/cloud-bill-dev): [Dev GCR Repo](https://console.cloud.google.com/gcr/images/cloud-bill-dev?project=cloud-bill-dev)
Production (gcp-marketplace-solutions/cje-marketplace-dev): [Prod GCR Repo](https://console.cloud.google.com/gcr/images/cloud-bill-dev?project=cje-marketplace-dev)

#### Logging and Audits
Application logging and security audits are provided by Google Stackdriver. 

Development (cloud-bill-dev/cloud-bill-dev): [Dev Stackdriver](https://console.cloud.google.com/logs/viewer?project=cloud-bill-dev&organizationId=41792434410&minLogLevel=0&expandAll=false&timestamp=2019-10-14T20:33:35.889000000Z&customFacets=&limitCustomFacetWidth=true&dateRangeStart=2019-10-14T19:33:36.144Z&dateRangeEnd=2019-10-14T20:33:36.144Z&interval=PT1H&resource=audited_resource&scrollTimestamp=2019-10-14T20:33:18.273999553Z)
Production (gcp-marketplace-solutions/cje-marketplace-dev): [Prod Stackdriver](https://console.cloud.google.com/logs/viewer?organizationId=41792434410&project=cje-marketplace-dev&minLogLevel=0&expandAll=false&timestamp=2019-10-14T20:33:02.130000000Z&customFacets=&limitCustomFacetWidth=true&dateRangeStart=2019-10-14T19:33:02.381Z&dateRangeEnd=2019-10-14T20:33:02.381Z&interval=PT1H&resource=audited_resource&scrollTimestamp=2019-10-14T20:27:52.458843159Z)
