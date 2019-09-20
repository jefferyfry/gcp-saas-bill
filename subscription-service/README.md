# Subscription Service
This directory contains the code for the subscription service which manages the 
subscriptions coming from the GCP marketplace.

## Configuration
To successfully run the subscription service, configuration must be set through either environment variables, command-line options or a configuration file. You may chose an option based on on your intent (development, testing, production deployment). The following configuration is required:

* Subscription Service Endpoint - This is the listening port for the service.
* GCP Project ID - This is your marketplace project where this service and required resources are deployed.

### Configuration Precedence
command-line options > environment variables

### Environment Variables
* CLOUD_BILLING_SUBSCRIPTION_CONFIG_FILE - Path to a configuration file (see below).
* CLOUD_BILLING_SUBSCRIPTION_SERVICE_ENDPOINT - _Subscription Service Endpoint_ from above.
* CLOUD_BILLING_SUBSCRIPTION_GCP_PROJECT_ID - _GCP Project ID_ from above.

* **GOOGLE_APPLICATION_CREDENTIALS** - This is the path to your GCP service account credentials required to access GCP resources like Datastore. This is a required environment variable for production.

### Command-Line Options
* configFile - Path to a configuration file (see below).
* subscriptionServiceEndpoint - _Subscription Service Endpoint_ from above.
* gcpProjectId - _GCP Project ID_ from above.

### Configuration File
The configFile command-line option or CLOUD_BILLING_SAAS_CONFIG_FILE environment variable requires a path to a JSON file with the configuration. Example:
```
{
  "subscriptionServiceEndpoint": ":8085",
  "gcpProjectId": "cloud-billing"
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
            - name: GOOGLE_APPLICATION_CREDENTIALS
              value: /auth/datastore-service-account/service-account.json
            - name: CLOUD_BILLING_SUBSCRIPTION_CONFIG_FILE
              value: /auth/subscription-service-config/subscription-service-config.json
          ports:
            - containerPort: 8085
          volumeMounts:
            - name: datastore-service-account
              mountPath: "/auth/datastore-service-account"
              readOnly: true
            - name: subscription-service-config
              mountPath: "/auth/subscription-service-config"
              readOnly: true
      volumes:
        - name: datastore-service-account
          secret:
            secretName: datastore-service-account
        - name: subscription-service-config
          secret:
            secretName: subscription-service-config
```

## Persistence
The subscription service uses GCP Datastore/Firestore NoSQL. Connecting the service requires setting the environment variable **GOOGLE_APPLICATION_CREDENTIALS**. This is the path to your GCP service account credentials required to access GCP resources like Datastore. Also ensure that you have set the correct GCP Project ID. This should be the same as where you created your Datastore database. 

### Creating the GCP Service Account to Access Datastore
Follow the instructions [here](https://cloud.google.com/datastore/docs/activate#other-platforms) to create a service account with permission Cloud Datastore Owner and download the key.

Then create the kubernetes secret.
```
kubectl create secret generic datastore-service-account --from-file datastore-service-account.json
```

### Using the Datastore Emulator
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
docker run -it --rm -p 8085:8085 -e CLOUD_BILLING_SAAS_SUBSCRIPTION_SERVICE_ENDPOINT=8085 -e CLOUD_BILLING_SAAS_CLOUD_COMMERCE_PROCUREMENT_URL='https://cloudcommerceprocurement.googleapis.com/' -e CLOUD_BILLING_SAAS_PARTNER_ID='123456' -e CLOUD_BILLING_SAAS_GCP_PROJECT_ID='gcp-project-1' --name my-subscription-service subscription-service-1:<tag>

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

## Testing
TBD