# Jenkins Support Subscription Subscription Service
This directory contains the code for the subscription service which manages the 
Jenkins Support subscriptions coming from the GCP marketplace.

## Configuration
To successfully run the subscription service, configuration must be set through either environment variables, command-line options or a configuration file. You may chose an option based on on your intent (development, testing, production deployment). The following configuration is required:

* Subscription Service Endpoint - This is the listening port for the service.
* Cloud Commerce Procurement URL - This is the marketplace API url for querying and approving subscriptions. See [here](https://cloud.google.com/marketplace/docs/partners/commerce-procurement-api/reference/rest/).
* Partner ID - This is the unique partner ID to include in posts.
* GCP Project ID - This is your marketplace project where this service and required resources are deployed.

### Precedence
command-line options > environment variables

### Environment Variables
JENKINS_SUPPORT_SAAS_CONFIG_FILE
JENKINS_SUPPORT_SAAS_SUBSCRIPTION_SERVICE_ENDPOINT
JENKINS_SUPPORT_SAAS_CLOUD_COMMERCE_PROCUREMENT_URL
JENKINS_SUPPORT_SAAS_PARTNER_ID
JENKINS_SUPPORT_SAAS_GCP_PROJECT_ID

GOOGLE_APPLICATION_CREDENTIALS

### Command-Line Options
configFile
subscriptionServiceEndpoint
cloudCommerceProcurementUrl
partnerId
gcpProjectId

### Configuration File
```
{
  "subscriptionServiceEndpoint": ":8085",
  "cloudCommerceProcurementUrl": "https://cloudcommerceprocurement.googleapis.com/",
  "partnerId": "0000",
  "gcpProjectId": "jenkins-support-saas"
}
```

## Persistence

## Running Locally
The following will be and run the service locally.
```
go run main.go <optional command-line options>
```

## Swagger

## Testing