# Jenkins Support Subscription Service
This directory contains the code for the subscription service which manages the 
Jenkins Support subscriptions coming from the GCP marketplace.

## Configuration
To successfully run the subscription service, configuration must be set through either environment variables, command-line options or a configuration file. You may chose an option based on on your intent (development, testing, production deployment). The following configuration is required:

* Subscription Service Endpoint - This is the listening port for the service.
* Cloud Commerce Procurement URL - This is the marketplace API url for querying and approving subscriptions. See [here](https://cloud.google.com/marketplace/docs/partners/commerce-procurement-api/reference/rest/).
* Partner ID - This is the unique partner ID to include in posts.
* GCP Project ID - This is your marketplace project where this service and required resources are deployed.

### Configuration Precedence
command-line options > environment variables

### Environment Variables
* JENKINS_SUPPORT_SAAS_CONFIG_FILE - Path to a configuration file (see below).
* JENKINS_SUPPORT_SAAS_SUBSCRIPTION_SERVICE_ENDPOINT - _Subscription Service Endpoint_ from above.
* JENKINS_SUPPORT_SAAS_CLOUD_COMMERCE_PROCUREMENT_URL - _Cloud Commerce Procurement URL_ from above.
* JENKINS_SUPPORT_SAAS_PARTNER_ID - _Partner ID_ from above.
* JENKINS_SUPPORT_SAAS_GCP_PROJECT_ID - _GCP Project ID_ from above.

* **GOOGLE_APPLICATION_CREDENTIALS** This is the path to your GCP service account credentials required to access GCP resources like Datastore. This is a required environment variable.

### Command-Line Options
* configFile - Path to a configuration file (see below).
* subscriptionServiceEndpoint - _Subscription Service Endpoint_ from above.
* cloudCommerceProcurementUrl - _Cloud Commerce Procurement URL_ from above.
* partnerId - _Partner ID_ from above.
* gcpProjectId - _GCP Project ID_ from above.

### Configuration File
The configFile command-line option or JENKINS_SUPPORT_SAAS_CONFIG_FILE environment variable requires a path to a JSON file with the configuration. Example:
```
{
  "subscriptionServiceEndpoint": ":8085",
  "cloudCommerceProcurementUrl": "https://cloudcommerceprocurement.googleapis.com/",
  "partnerId": "0000",
  "gcpProjectId": "cloudbees-jenkins-support"
}
```

## Persistence
The subscription service uses GCP Datastore/Firestore NoSQL. Connecting the service requires setting the environment variable **GOOGLE_APPLICATION_CREDENTIALS**. This is the path to your GCP service account credentials required to access GCP resources like Datastore. Also ensure that you have set the correct GCP Project ID. This should be the same as where you created your Datastore database. 

### Using the Datastore Emulator
For development and testing, GCP provides a [Datastore emulator](https://cloud.google.com/datastore/docs/tools/datastore-emulator). Follow the [instructions](https://cloud.google.com/datastore/docs/tools/datastore-emulator#installing_the_emulator) to install the emulator.

When running the subscription service locally, you may need to set environment variables for the service to connect to the emulator. Take note of the emulator output to get the correct emulator port. Here is an example of setting these:

```
export DATASTORE_DATASET=cloudbees-jenkins-support
export DATASTORE_EMULATOR_HOST=::1:8039
export DATASTORE_EMULATOR_HOST_PATH=::1:8039/datastore
export DATASTORE_HOST=http://::1:8039
export DATASTORE_PROJECT_ID=cloudbees-jenkins-support
```

## Running Locally
The following will be and run the service locally.
```
go run main.go <optional command-line options>
```

## Swagger
A Swagger UI is embedded in the service. It is located at http(s)://host:port/swagger/index.html

### Updating Swagger
[Swaggo](https://github.com/swaggo/swag) is used to generate the Swagger UI. Swaggo uses annotations on main.go and rest/handlers.go to generate the Swagger UI. To update after changes, execute:

```
swag init
```

in the subscription-service directory.

## Testing
TBD