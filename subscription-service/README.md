# Subscription Service
This directory contains the code for the subscription service which manages the 
subscriptions coming from the GCP marketplace.

## Configuration
To successfully run the subscription service, configuration must be set through either environment variables, command-line options or a configuration file. You may chose an option based on on your intent (development, testing, production deployment). The following configuration is required:

* Subscription Service Endpoint - This is the listening port for the service.
* Subscription Service Health Check Endpoint - Listening port for Kubernetes health checks (readiness and liveness).
* GCP Project ID - This is your marketplace project where this service and required resources are deployed.
* Sentry DSN - This is the key for Sentry logging.

### Configuration Precedence
command-line options > environment variables

### Environment Variables
* CLOUD_BILL_SUBSCRIPTION_CONFIG_FILE - Path to a configuration file (see below).
* CLOUD_BILL_SUBSCRIPTION_SERVICE_ENDPOINT 
* CLOUD_BILL_SUBSCRIPTION_HEALTH_CHECK_ENDPOINT
* CLOUD_BILL_SUBSCRIPTION_GCP_PROJECT_ID 
* CLOUD_BILL_DATASTORE_BACKUP_SENTRY_DSN

* **GOOGLE_APPLICATION_CREDENTIALS** - This is the path to your GCP service account credentials required to access GCP resources like Datastore. This is a required environment variable for production.

### Command-Line Options
* configFile - Path to a configuration file (see below).
* subscriptionServiceEndpoint
* healthCheckEndpoint
* gcpProjectId
* sentryDsn

### Configuration File
The configFile command-line option or CLOUD_BILL_SAAS_CONFIG_FILE environment variable requires a path to a JSON file with the configuration. Example:
```
{
  "subscriptionServiceEndpoint": ":8085",
  "healthCheckEndpoint": "8095",
  "gcpProjectId": "cloud-billing",
  "sentryDsn": "https://xxx"
}
```

### Production Configuration
For production, it is highly recommended that the service configuration be set by using the configuration file option. Set this configuration file as a kubernetes secret since there are sensitive parameters in the configuration:

```
kubectl create secret generic subscription-service-config --from-file subscription-service-config.json
```

Then mount the file and set it as an environment variable.

```
    spec:
      containers:
        - name: subscription-service
          image: gcr.io/cloud-bill-dev/subscription-service:latest
          env:
#            - name: CLOUD_BILL_SUBSCRIPTION_SERVICE_ENDPOINT
#              value: "8085"
#            - name: CLOUD_BILL_SUBSCRIPTION_HEALTH_CHECK_ENDPOINT
#              value: "8095"
#            - name: CLOUD_BILL_SUBSCRIPTION_GCP_PROJECT_ID
#              value: "<yourprojectid>"
#            - name: CLOUD_BILL_SUBSCRIPTION_SENTRY_DSN
#              value: "dsn"               
            - name: GOOGLE_APPLICATION_CREDENTIALS
              value: /auth/gcp-service-account/gcp-service-account.json
            - name: CLOUD_BILL_SUBSCRIPTION_CONFIG_FILE
              value: /auth/subscription-service-config/subscription-service-config.json
          ports:
            - containerPort: 8085
          volumeMounts:
            - name: gcp-service-account
              mountPath: "/auth/gcp-service-account"
              readOnly: true
            - name: subscription-service-config
              mountPath: "/auth/subscription-service-config"
              readOnly: true
      volumes:
        - name: gcp-service-account
          secret:
            secretName: gcp-service-account
        - name: subscription-service-config
          secret:
            secretName: subscription-service-config
```

## GCP Service Accounts
The subscription service requires setting the environment variable **GOOGLE_APPLICATION_CREDENTIALS**. This is the path to your GCP service account credentials. Also ensure that you have set the correct GCP Project ID. This should be the same as where you created your Datastore database. 

The following roles are required:
* Cloud Datastore Owner - Used for the Cloud Datastore subscription DB.
It is recommended that the roles be used assigned to a common service account. Then the service account file can be shared and mounted for all the services.

Then create the kubernetes secret.
```
kubectl create secret generic gcp-service-account --from-file gcp-service-account.json
```

## Using the Datastore Emulator
For development and testing, GCP provides a [Datastore emulator](https://cloud.google.com/datastore/docs/tools/datastore-emulator). Follow the [instructions](https://cloud.google.com/datastore/docs/tools/datastore-emulator#installing_the_emulator) to install the emulator. Then start the datastore emulator:

```
gcloud beta emulators datastore start
```
When running the subscription service locally, you may need to set environment variables for the service to connect to the emulator. Take note of the emulator output to get the correct emulator port. Here is an example of setting these:

```
export DATASTORE_DATASET=cloud-bill
export DATASTORE_EMULATOR_HOST=::1:8039
export DATASTORE_EMULATOR_HOST_PATH=::1:8039/datastore
export DATASTORE_HOST=http://::1:8039
export DATASTORE_PROJECT_ID=cloud-bill
```

### Importing Cloud Datastore DB to the Emulator for Testing
1. Follow these [instructions] to create a GCS bucket.
2. Export the database to the GCS bucket. Ensure you are authenticated, have the correct permissions, and have the correct project set.
```
gcloud datastore export gs://${BUCKET} --async
```
3. Use gsutil to copy the bucket to your local directory.
```
gsutil cp -r gs://cloud-bill-dev.appspot.com .
```
4. With the emulator running, use curl to import the data.
```
curl -X POST localhost:8081/v1/projects/[PROJECT_ID]:import \
-H 'Content-Type: application/json' \
-d '{"input_url":"[ENTITY_EXPORT_FILES]"}'

ex.
curl -X POST localhost:8116/v1/projects/cloud-bill-dev:import \
-H 'Content-Type: application/json' \
-d '{"input_url":"/Users/jefferyfry/tmp/emulator_db/cloud-bill-dev.appspot.com/2019-09-26T15:54:16_87320/2019-09-26T15:54:16_87320.overall_export_metadata"}'
```
Modify localhost:8081 if the emulator uses a different port.

### Google Sheets to View the Cloud Datastore DB
The [google-sheets directory](google-sheets/datastore-read-only.gs) contains a Google Apps Script that you can use to pull data from the Datastore DB and into a Google Sheet. Follow these [instructions](https://developers.google.com/apps-script/guides/sheets) to execute the script for a Google Sheet. [Time-driven triggers](https://developers.google.com/apps-script/guides/triggers/installable#time-driven_triggers) can be used to automatically update the Google Sheet on a schedule. [Spreadsheet actions](https://developers.google.com/apps-script/guides/triggers/installable#g_suite_application_triggers) can also trigger updates. Running the script requires having a service account that has the Datastore Viewer IAM permission. Then place the Service Account json file in the same directory as the Google Sheet and name it _datastore-viewer-service-account.json_. 

## Running Locally
The following will run the service locally.
```
go run main.go <optional command-line options>
```

## Building the docker image locally
```
docker build -t subscription-service:<tag> .

ex.
docker build -t subscription-service:1 .
```

## Pushing to GCR
```
docker tag subscription-service:<tag> gcr.io/<path>/subscription-service:<tag>

docker push gcr.io/<path>/subscription-service:<tag>

ex.
docker tag subscription-service:1 gcr.io/cloud-bill-dev/subscription-service:1

docker push gcr.io/cloud-bill-dev/subscription-service:1

```

## Running the docker image locally with environment variables
```
docker run -it --rm -p 8085:8085 -e CLOUD_BILL_SAAS_SUBSCRIPTION_SERVICE_ENDPOINT=8085 -e CLOUD_BILL_SAAS_CLOUD_COMMERCE_PROCUREMENT_URL='https://cloudcommerceprocurement.googleapis.com/' -e CLOUD_BILL_SAAS_PARTNER_ID='123456' -e CLOUD_BILL_SAAS_GCP_PROJECT_ID='gcp-project-1' --name my-subscription-service subscription-service-1:<tag>

```

## Upgrades
It is recommended that you use K8s rolling update for service images. This can be executed in a single command:

```
kubectl set image deployments/<deployment-name> <container>=image

eg.
kubectl set image deployments/subscription-service subscription-service=gcr.io/cloud-bill-dev/subscription-service:2
```
If the upgrade includes configuration changes, apply those configuration changes first.

You can monitor the rolling update using:

```
kubectl rollout status deployments/<deployment-name>
```

## Swagger
A Swagger UI is embedded in the service. It is located at http(s)://host:port/swagger/index.html

![Swagger](https://user-images.githubusercontent.com/6440106/63872211-430eaa00-c972-11e9-93b9-fd417ae02eb8.png)

### Updating Swagger
[Swaggo](https://github.com/swaggo/swag) is used to generate the Swagger UI. Swaggo uses annotations on main.go and rest/handlers.go to generate the Swagger UI. To update after changes, execute:

```
swag init
```

in the subscription-service directory.

