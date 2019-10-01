# PubSub Service
This directory contains the code for the pubsub service which manages the 
pubsubs coming from the GCP marketplace.

## Configuration
To successfully run the pubsub service, configuration must be set through either environment variables, command-line options or a configuration file. You may chose an option based on on your intent (development, testing, production deployment). The following configuration is required:

* PubSub Subscription -
* PubSub Topic -
* Subscription Service URL -
* Cloud Commerce Procurement URL
* Partner ID -
* GCP Project ID - This is your marketplace project where this service and required resources are deployed.

### Configuration Precedence
command-line options > environment variables

### Environment Variables
* CLOUD_BILL_PUBSUB_CONFIG_FILE - Path to a configuration file (see below).
* CLOUD_BILL_PUBSUB_SUBSCRIPTION - _Subscription Service Endpoint_ from above.
* CLOUD_BILL_PUBSUB_TOPIC_PREFIX
* CLOUD_BILL_SUBSCRIPTION_SERVICE_URL
* CLOUD_BILL_PUBSUB_CLOUD_COMMERCE_PROCUREMENT_URL
* CLOUD_BILL_PUBSUB_PARTNER_ID
* CLOUD_BILL_PUBSUB_GCP_PROJECT_ID - _GCP Project ID_ from above.

* **GOOGLE_APPLICATION_CREDENTIALS** - This is the path to your GCP service account credentials required to access GCP PubSub and Cloud Commerce Procurement API. This is a required environment variable for production.

### Command-Line Options
* configFile - Path to a configuration file (see below).
* pubSubSubscription - _Subscription Service Endpoint_ from above.
* pubSubTopicPrefix
* subscriptionServiceUrl
* cloudCommerceProcurementUrl
* partnerId
* gcpProjectId - _GCP Project ID_ from above.

### Configuration File
The configFile command-line option or CLOUD_BILL_SAAS_CONFIG_FILE environment variable requires a path to a JSON file with the configuration. Example:
```
{
  "pubSubSubscription": "codelab",
  "subscriptionServiceUrl": "http://subscription-service.default.svc.cluster.local:8085/api/v1/",
  "cloudCommerceProcurementUrl": "https://cloudcommerceprocurement.googleapis.com/v1/",
  "partnerId": "DEMO-codelab-project",
  "gcpProjectId": "cloud-bill-dev"
}
```

### Production Configuration
For production, it is highly recommended that the service configuration be set by using the configuration file option. Set this configuration file as a kubernetes secret since there are sensitive parameters in the configuration:

```
kubectl create secret generic pubsub-service-config --from-file pubsub-service-config.json
```

Then mount the file and set it as an environment variable.

```
    spec:
      containers:
        - name: pubsub-service
          image: gcr.io/cloud-bill-dev/pubsub-service:latest
          env:
            #            - name: CLOUD_BILL_PUBSUB_SUBSCRIPTION
            #              value: "codelab"
            #            - name: CLOUD_BILL_SUBSCRIPTION_SERVICE_URL
            #              value: "https://subscription-service.cloudbees-jenkins-support.svc.cluster.local"
            #            - name: CLOUD_BILL_PUBSUB_CLOUD_COMMERCE_PROCUREMENT_URL
            #              value: "https://cloudcommerceprocurement.googleapis.com/"
            #            - name: CLOUD_BILL_PUBSUB_PARTNER_ID
            #              value: "<yourpartnerid>"
            #            - name: CLOUD_BILL_PUBSUB_GCP_PROJECT_ID
            #              value: "<yourprojectid>"
            - name: GOOGLE_APPLICATION_CREDENTIALS
              value: /auth/gcp-service-account/gcp-service-account.json
            - name: CLOUD_BILL_PUBSUB_CONFIG_FILE
              value: /auth/pubsub-service-config/pubsub-service-config.json
          volumeMounts:
            - name: gcp-service-account
              mountPath: "/auth/gcp-service-account"
              readOnly: true
            - name: pubsub-service-config
              mountPath: "/auth/pubsub-service-config"
              readOnly: true
      volumes:
        - name: gcp-service-account
          secret:
            secretName: gcp-service-account
        - name: pubsub-service-config
          secret:
            secretName: pubsub-service-config
```

## GCP Service Accounts
The pubsub service requires setting the environment variable **GOOGLE_APPLICATION_CREDENTIALS**. This is the path to your GCP service account credentials.

The following roles are required:
* PubSub Editor - Used to access marketplace PubSub events.
* Cloud Commerce API (assigned by GCP Marketplace team) - Allows access to the Cloud Commerce API
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
docker build -t pubsub-service:<tag> .

ex.
docker build -t pubsub-service:1 .
```

## Pushing to GCR
```
docker tag pubsub-service:<tag> gcr.io/<path>/pubsub-service:<tag>

docker push gcr.io/<path>/pubsub-service:<tag>

ex.
docker tag pubsub-service:1 gcr.io/cloud-bill-dev/pubsub-service:1

docker push gcr.io/cloud-bill-dev/pubsub-service:1

```

## Running the docker image locally with environment variables
```
docker run -it --rm -p 8085:8085 -e CLOUD_BILL_SAAS_PUBSUB_SERVICE_ENDPOINT=8085 -e CLOUD_BILL_SAAS_CLOUD_COMMERCE_PROCUREMENT_URL='https://cloudcommerceprocurement.googleapis.com/' -e CLOUD_BILL_SAAS_PARTNER_ID='123456' -e CLOUD_BILL_SAAS_GCP_PROJECT_ID='gcp-project-1' --name my-pubsub-service pubsub-service-1:<tag>

```

## Testing
TBD
