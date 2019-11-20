# Entitlement Check
This directory contains the code for the entitlement check cron job which checks the status of entitlements for specified products. This is required
for Google VM solution offerings which are not integrated into the marketplace pubsub events service.

## Configuration
To successfully run the entitlement check cron job, configuration must be set through either environment variables, command-line options or a configuration file. You may chose an option based on on your intent (development, testing, production deployment). The following configuration is required:

* GCP Project ID - This is your marketplace project where this service and required resources are deployed.
* Products - These are the VM products to check.
* Subscription Service URL - This is the URL to the subscription service.
* Google Subscription URL - This is the URL to the Google subscription service for querying entitlements.
* Sentry DSN - This is the key for Sentry logging.

### Configuration Precedence
command-line options > environment variables

### Environment Variables
* CLOUD_BILL_ENTITLEMENT_CHECK_CONFIG_FILE - Path to a configuration file (see below).
* CLOUD_BILL_ENTITLEMENT_CHECK_GCP_PROJECT_ID 
* CLOUD_BILL_ENTITLEMENT_CHECK_GCP_PRODUCTS
* CLOUD_BILL_ENTITLEMENT_CHECK_SUBSCRIPTION_SERVICE_URL
* CLOUD_BILL_ENTITLEMENT_CHECK_GOOGLE_SUBSCRIPTIONS_URL
* CLOUD_BILL_ENTITLEMENT_CHECK_SENTRY_DSN

* **GOOGLE_APPLICATION_CREDENTIALS** - This is the path to your GCP service account credentials required to access Cloud Datastore and your GCS Bucket. This is a required environment variable for production.

### Command-Line Options
* configFile - Path to a configuration file (see below).
* gcpProjectId 
* products
* subscriptionServiceUrl 
* googleSubscriptionServiceUrl 
* sentryDsn

### Configuration File
The configFile command-line option or CLOUD_BILL_SAAS_CONFIG_FILE environment variable requires a path to a JSON file with the configuration. Example:
```
{
  "gcpProjectId": "cloud-bill-dev",
  "products": "cloudbees-accelerator",
  "subscriptionServiceUrl": "http://subscription-service.default.svc.cluster.local:8085/api/v1/",
  "googleSubscriptionsUrl": "https://cloudbilling.googleapis.com/v1",
  "sentryDsn": "https://xxx"
}
```

### Production Configuration
For production, it is highly recommended that the service configuration be set by using the configuration file option. Set this configuration file as a kubernetes secret since there are sensitive parameters in the configuration:

```
kubectl create secret generic entitlement-check-config --from-file entitlement-check-config.json
```

Then mount the file and set it as an environment variable.

```
        spec:
          schedule: "0 1 * * *"
          jobTemplate:
            spec:
              template:
                spec:
                  containers:
                    - name: entitlement-check
                      image: gcr.io/cje-marketplace-dev/cloud-bill-saas/entitlement-check:latest
                      env:
                        - name: GOOGLE_APPLICATION_CREDENTIALS
                          value: /auth/gcp-service-account/gcp-service-account.json
                        - name: CLOUD_BILL_ENTITLEMENT_CHECK_CONFIG_FILE
                          value: /auth/entitlement-check-config/entitlement-check-config.json
                      volumeMounts:
                        - name: gcp-service-account
                          mountPath: "/auth/gcp-service-account"
                          readOnly: true
                        - name: entitlement-check-config
                          mountPath: "/auth/entitlement-check-config"
                          readOnly: true
                  restartPolicy: Never
                  volumes:
                    - name: gcp-service-account
                      secret:
                        secretName: gcp-service-account
                    - name: entitlement-check-config
                      secret:
                        secretName: entitlement-check-config
```

## GCP Service Accounts
The service requires setting the environment variable **GOOGLE_APPLICATION_CREDENTIALS**. This is the path to your GCP service account credentials.

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
docker build -t entitlement-check:<tag> .

ex.
docker build -t entitlement-check:1 .
```

## Pushing to GCR
```
docker tag entitlement-check:<tag> gcr.io/<path>/entitlement-check:<tag>

docker push gcr.io/<path>/entitlement-check:<tag>

ex.
docker tag entitlement-check:1 gcr.io/cloud-bill-dev/entitlement-check:1

docker push gcr.io/cloud-bill-dev/entitlement-check:1

```

## Running the docker image locally with environment variables
```
docker run -it --rm -e CLOUD_BILL_ENTITLEMENT_CHECK_GCS_BUCKET=gs://bucket -e CLOUD_BILL_SAAS_GCP_PROJECT_ID='gcp-project-1' --name my-entitlement-check entitlement-check:<tag>

```

### Upgrades
This can be executed in a single command:

```
kubectl set image cronjob/<cronjob-name> <container>=image

eg.
kubectl set image cronjob/entitlement-check entitlement-check=gcr.io/cloud-bill-dev/entitlement-check:2
```
If the upgrade includes configuration changes, apply those configuration changes first.

