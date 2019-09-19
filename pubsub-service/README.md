#Subscription Service
This directory contains the code for the pubsub service which manages the 
pubsubs coming from the GCP marketplace.

## Configuration
To successfully run the pubsub service, configuration must be set through either environment variables, command-line options or a configuration file. You may chose an option based on on your intent (development, testing, production deployment). The following configuration is required:

* Subscription Service Endpoint - This is the listening port for the service.
* GCP Project ID - This is your marketplace project where this service and required resources are deployed.

### Configuration Precedence
command-line options > environment variables

### Environment Variables
* CLOUD_BILLING_PUBSUB_CONFIG_FILE - Path to a configuration file (see below).
* CLOUD_BILLING_PUBSUB_SERVICE_ENDPOINT - _Subscription Service Endpoint_ from above.
* CLOUD_BILLING_PUBSUB_GCP_PROJECT_ID - _GCP Project ID_ from above.

* **GOOGLE_APPLICATION_CREDENTIALS** - This is the path to your GCP service account credentials required to access GCP resources like PubSub. This is a required environment variable for production.

### Command-Line Options
* configFile - Path to a configuration file (see below).
* pubsubServiceEndpoint - _Subscription Service Endpoint_ from above.
* gcpProjectId - _GCP Project ID_ from above.

### Configuration File
The configFile command-line option or CLOUD_BILLING_SAAS_CONFIG_FILE environment variable requires a path to a JSON file with the configuration. Example:
```
{
  "pubSubSubscription": "codelab",
  "pubSubTopicPrefix": "DEMO-",
  "cloudCommerceProcurementUrl": "https://cloudcommerceprocurement.googleapis.com/",
  "partnerId": "DEMO-codelab-project",
  "gcpProjectId": "cloudbees-jenkins-support-dev"
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
          image: gcr.io/cloudbees-jenkins-support-dev/pubsub-service:latest
          env:
            #            - name: CLOUD_BILLING_AGENT_PUBSUB_SUBSCRIPTION
            #              value: "codelab"
            #            - name: CLOUD_BILLING_AGENT_PUBSUB_TOPIC_PREFIX
            #              value: "DEMO-"
            #            - name: CLOUD_BILLING_AGENT_CLOUD_COMMERCE_PROCUREMENT_URL
            #              value: "https://cloudcommerceprocurement.googleapis.com/"
            #            - name: CLOUD_BILLING_AGENT_PARTNER_ID
            #              value: "<yourpartnerid>"
            #            - name: CLOUD_BILLING_AGENT_GCP_PROJECT_ID
            #              value: "<yourprojectid>"
            - name: GOOGLE_APPLICATION_CREDENTIALS
              value: /auth/pubsub-service-account/service-account.json
            - name: CLOUD_BILLING_SUBSCRIPTION_CONFIG_FILE
              value: /auth/pubsub-service-config/pubsub-service-config.json
          ports:
            - containerPort: 8085
          volumeMounts:
            - name: pubsub-service-account
              mountPath: "/auth/pubsub-service-account"
              readOnly: true
            - name: pubsub-service-config
              mountPath: "/auth/pubsub-service-config"
              readOnly: true
      volumes:
        - name: pubsub-service-account
          secret:
            secretName: pubsub-service-account
        - name: pubsub-service-config
          secret:
            secretName: pubsub-service-config
```

## PubSub
The pubsub service uses GCP PubSub. Connecting the service requires setting the environment variable **GOOGLE_APPLICATION_CREDENTIALS**. This is the path to your GCP service account credentials required to access GCP resources like PubSub. Also ensure that you have set the correct GCP Project ID. 

### Creating the GCP Service Account to Access PubSub
Follow the instructions [here](https://cloud.google.com/pubsub/docs/reference/libraries#setting_up_authentication) to create a service account and download the key.

Then create the kubernetes secret.
```
kubectl create secret generic pubsub-service-account --from-file pubsub-service-account.json
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
docker tag pubsub-service:1 gcr.io/cloudbees-jenkins-support-dev/pubsub-service:1

docker push gcr.io/cloudbees-jenkins-support-dev/pubsub-service:1

```

## Running the docker image locally with environment variables
```
docker run -it --rm -p 8085:8085 -e CLOUD_BILLING_SAAS_PUBSUB_SERVICE_ENDPOINT=8085 -e CLOUD_BILLING_SAAS_CLOUD_COMMERCE_PROCUREMENT_URL='https://cloudcommerceprocurement.googleapis.com/' -e CLOUD_BILLING_SAAS_PARTNER_ID='123456' -e CLOUD_BILLING_SAAS_GCP_PROJECT_ID='gcp-project-1' --name my-pubsub-service pubsub-service-1:<tag>

```

## Testing
TBD
