# Frontend Service
The Frontend service provides the UI for customer signup from the marketplace. The end result is storing the customer account information for the subscription and confirming the account with Google. Auth0 and Google Identity are used to capture some of the customer profile data.

## Frontend Flow
The basic frontend flow amongst handlers and pages is the following:
![Jenkins Support SaaS - Page 4](https://user-images.githubusercontent.com/6440106/64573203-54b36280-d31f-11e9-84cb-9e0ca4e5fc67.png)

## Pages
* [signup.html](https://github.com/cloudbees/cloud-bill-saas/tree/master/frontend-service/templates/signup.html) - Initial page to direct customer to Auth0/Google sign in. The customer is sent to this page from marketplace.
* [confirmProd.html](https://github.com/cloudbees/cloud-bill-saas/tree/master/frontend-service/templates/confirmProd.html) - Auth0/Google callback page to confirm account information for VM and K8s products.
* [confirmSaas.html](https://github.com/cloudbees/cloud-bill-saas/tree/master/frontend-service/templates/confirmSaas.html) - Auth0/Google callback page to confirm account information for Saas products.
* [finish.html](https://github.com/cloudbees/cloud-bill-saas/tree/master/frontend-service/templates/finish.html) - Final page to confirm account creation and notify customer of next steps.

## Configuration
To successfully run the subscription service, configuration must be set through either environment variables, command-line options or a configuration file. You may chose an option based on on your intent (development, testing, production deployment). The following configuration is required:

* Frontend Service Endpoint - This is the listening port for the frontend service.
* Frontend Health Check Endpoint - This is the listening port for the health check service.
* Subscription Service URL - This is the URL to the subscription service.
* Google Subscription URL - This is the URL to the Google subscription service for querying VM entitlements.
* Client ID - This is the Oauth/Auth0 client ID.
* Client Secret - This is the Oauth/Auth0 client secret.
* Callback URL - This is the callback URL used by the Oauth/Auth0 service.
* Issuer - The Oauth issuer or Auth0 domain.
* Session Key - A random character sequence for session encoding.
* Cloud Commerce Procurement URL - This is the marketplace API url for querying and approving subscriptions. See [here](https://cloud.google.com/marketplace/docs/partners/commerce-procurement-api/reference/rest/).
* Partner ID - This is the unique partner ID to include in posts.
* FinishUrl - This is the url that the customer can go to after completing the signup.
* FinishUrlTitle - This is the button title of the FinishUrl.
* Test Mode - Runs the service in test mode and provides handlers /signupsaastest?acct=<acct> and /resetsaas?acct=<acct>.
* Sentry DSN - This is the key for Sentry logging.

### Configuration Precedence
command-line options > environment variables

### Environment Variables
* CLOUD_BILL_FRONTEND_CONFIG_FILE - Path to a configuration file (see below).
* CLOUD_BILL_FRONTEND_SERVICE_ENDPOINT 
* CLOUD_BILL_FRONTEND_HEALTH_CHECK_ENDPOINT
* CLOUD_BILL_FRONTEND_SUBSCRIPTION_SERVICE_URL 
* CLOUD_BILL_FRONTEND_GOOGLE_SUBSCRIPTIONS_URL
* CLOUD_BILL_FRONTEND_CLIENT_ID 
* CLOUD_BILL_FRONTEND_CLIENT_SECRET 
* CLOUD_BILL_FRONTEND_CALLBACK_URL
* CLOUD_BILL_FRONTEND_ISSUER
* CLOUD_BILL_FRONTEND_SESSION_KEY
* CLOUD_BILL_FRONTEND_CLOUD_COMMERCE_PROCUREMENT_URL
* CLOUD_BILL_FRONTEND_PARTNER_ID
* CLOUD_BILL_FRONTEND_FINISH_URL
* CLOUD_BILL_FRONTEND_FINISH_URL_TITLE
* CLOUD_BILL_FRONTEND_TEST_MODE

### Command-Line Options
* configFile - Path to a configuration file (see below).
* frontendServiceEndpoint 
* healthCheckEndpoint
* subscriptionServiceUrl 
* googleSubscriptionServiceUrl
* clientId 
* clientSecret 
* callbackUrl 
* issuer 
* sessionKey
* cloudCommerceProcurementUrl
* partnerId
* finishUrl
* finishUrlTitle 
* sentryDsn

### Configuration File
The configFile command-line option or CLOUD_BILL_FRONTEND_CONFIG_FILE environment variable requires a path to a JSON file with the configuration. Example:
```
{
  "frontendServiceEndpoint": "8086",
  "healthCheckEndpoint": "8096",
  "subscriptionServiceUrl": "http://subscription-service.default.svc.cluster.local:8085/api/v1/",
  "googleSubscriptionsUrl": "https://cloudbilling.googleapis.com/v1",
  "clientId": "clientId",
  "clientSecret": "clientSecret",
  "callbackUrl": "https://cloud-bill.35.237.116.107.beesdns.com/callback",
  "issuer": "https://cloudbees-development.auth0.com/",
  "sessionKey": "somesessionkey",
  "cloudCommerceProcurementUrl": "https://cloudcommerceprocurement.googleapis.com/v1/",
  "partnerId": "DEMO-cloud-bill-dev",
  "finishUrl": "https://grandcentral.beescloud.com/login/login?login_redirect=https://go.beescloud.com"
  "finishUrlTitle": "Login"
  "testMode": "true",
  "gcpProjectId": "cloud-bill-dev",
  "sentryDsn": "https://xxx"
}
```

### Production Configuration
For production, it is highly recommended that the service configuration be set by using the configuration file option. Set this configuration file as a kubernetes secret since there are sensitive parameters in the configuration:

```
kubectl create secret generic frontend-service-config --from-file frontend-service-config.json
```

Then mount the file and set it as an environment variable.

```
    spec:
          containers:
            - name: frontend-service
              image: gcr.io/cloud-bill-dev/frontend-service:latest
              env:
    #            - name: CLOUD_BILL_FRONTEND_SERVICE_ENDPOINT
    #              value: "8086"
    #            - name: CLOUD_BILL_FRONTEND_HEALTH_CHECK_ENDPOINT
    #              value: "8096"
    #            - name: CLOUD_BILL_FRONTEND_SUBSCRIPTION_SERVICE_URL
    #              value: "http://subscription-service.default.svc.cluster.local:8085/api/v1"
    #            -CLOUD_BILL_FRONTEND_GOOGLE_SUBSCRIPTIONS_URL
    #              value: "https://cloudbilling.googleapis.com/v1"
    #            - name: CLOUD_BILL_FRONTEND_CLIENT_ID
    #              value: "<yourauth0clientid>"
    #            - name: CLOUD_BILL_FRONTEND_CLIENT_SECRET
    #              value: "<yourauth0clientsecret>"
    #            - name: CLOUD_BILL_FRONTEND_CALLBACK_URL
    #              value: "https://<yourhost>/callback"
    #            - name: CLOUD_BILL_FRONTEND_ISSUER
    #              value: "https://<yourdomain/>"
    #            - name: CLOUD_BILL_FRONTEND_SESSION_KEY
    #              value: "<yoursessionkey>"
    #            - name: CLOUD_BILL_FRONTEND_CLOUD_COMMERCE_PROCUREMENT_URL
    #              value: "https://cloudcommerceprocurement.googleapis.com/v1/"
    #            - name: CLOUD_BILL_FRONTEND_PARTNER_ID
    #              value: "000"
    #            - name: CLOUD_BILL_FRONTEND_FINISH_URL
    #              value: "https://grandcentral.cloudbees.com/login/login?login_redirect=https://go.cloudbees.com"
    #            - name: CLOUD_BILL_FRONTEND_FINISH_URL_TITLE
    #              value: "Login"
    #            - name: CLOUD_BILL_FRONTEND_TEST_MODE
    #              value: "false"
    #            - name: CLOUD_BILL_FRONTEND_SENTRY_DSN
    #              value: "dsn"               
                - name: GOOGLE_APPLICATION_CREDENTIALS
                  value: /auth/gcp-service-account/gcp-service-account.json
                - name: CLOUD_BILL_FRONTEND_CONFIG_FILE
                  value: /auth/frontend-service-config/frontend-service-config.json
              ports:
                - containerPort: 8086
              volumeMounts:
                - name: gcp-service-account
                  mountPath: "/gcp-service-account"
                  readOnly: true
                - name: frontend-service-config
                  mountPath: "/auth/frontend-service-config"
                  readOnly: true
          volumes:
            - name: gcp-service-account
              secret:
                secretName: gcp-service-account
            - name: frontend-service-config
              secret:
                secretName: frontend-service-config
```

## GCP Service Accounts
The frontend service requires setting the environment variable **GOOGLE_APPLICATION_CREDENTIALS**. This is the path to your GCP service account credentials.

The following roles are required:
* PubSub Editor - Used to access marketplace PubSub events.
* Cloud Commerce API (assigned by GCP Marketplace team) - Allows access to the Cloud Commerce API
* Billing Account Administrator (NOT FOR PRODUCTION) - Allows the reset of test accounts.
It is recommended that the roles be used assigned to a common service account. Then the service account file can be shared and mounted for all the services.
Then create the kubernetes secret.
```
kubectl create secret generic gcp-service-account --from-file gcp-service-account.json
```

## Running Locally
The following will run the service locally.
```
go run main.go <optional command-line options>
```

## Building the docker image locally
```
docker build -t frontend-service:<tag> .

ex. 
docker build -t frontend-service:1 .
```

## Pushing to GCR
```
docker tag frontend-service:<tag> gcr.io/<path>/frontend-service:<tag>

docker push gcr.io/<path>/frontend-service:<tag>

ex.
docker tag frontend-service:1 gcr.io/cloud-bill-dev/frontend-service:1

docker push gcr.io/cloud-bill-dev/frontend-service:1
```

## Running the docker image locally with environment variables
```
docker run -it --rm -p 8086:8086 -e CLOUD_BILL_SUB_FRONTEND_SERVICE_ENDPOINT=8086 -e CLOUD_BILL_SUB_SERVICE_URL='http://localhost:8085' -e CLOUD_BILL_SUB_FRONTEND_CLIENT_ID='abcdef' -e CLOUD_BILL_SUB_FRONTEND_CLIENT_SECRET='123456' -e CLOUD_BILL_SUB_FRONTEND_CALLBACK_URL='http://localhost:8085/callback' -e CLOUD_BILL_SUB_FRONTEND_ISSUER='issuer' -e CLOUD_BILL_SUB_FRONTEND_SESSION_KEY='somekeycloudbeesjenkinssupportsessionkey1cl0udb33s1' --name my-frontend-service frontend-service-1:<tag>

```

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