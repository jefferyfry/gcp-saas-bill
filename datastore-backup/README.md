# Datastore Backup
This directory contains the code for the pubsub service which manages the 
pubsubs coming from the GCP marketplace.

## Configuration
To successfully run the , configuration must be set through either environment variables, command-line options or a configuration file. You may chose an option based on on your intent (development, testing, production deployment). The following configuration is required:

* GCP Project ID - This is your marketplace project where this service and required resources are deployed.
* GCS Bucket - This is the Google Cloud Storage bucket where you want the backup files.

### Configuration Precedence
command-line options > environment variables

### Environment Variables
* CLOUD_BILL_DATASTORE_BACKUP_CONFIG_FILE - Path to a configuration file (see below).
* CLOUD_BILL_DATASTORE_BACKUP_GCP_PROJECT_ID 
* CLOUD_BILL_DATASTORE_BACKUP_GCS_BUCKET

* **GOOGLE_APPLICATION_CREDENTIALS** - This is the path to your GCP service account credentials required to access Cloud Datastore and your GCS Bucket. This is a required environment variable for production.

### Command-Line Options
* configFile - Path to a configuration file (see below).
* gcpProjectId 
* gcsBucket

### Configuration File
The configFile command-line option or CLOUD_BILL_SAAS_CONFIG_FILE environment variable requires a path to a JSON file with the configuration. Example:
```
{
  "gcpProjectId": "cloud-bill-dev",
  "gcsBucket": "gs://cloud-bill-dev.appspot.com"
}
```

### Production Configuration
For production, it is highly recommended that the service configuration be set by using the configuration file option. Set this configuration file as a kubernetes secret since there are sensitive parameters in the configuration:

```
kubectl create secret generic datastore-backup-config --from-file datastore-backup-config.json
```

Then mount the file and set it as an environment variable.

```
        spec:
          containers:
            - name: datastore-backup
              image: gcr.io/cloud-bill-dev/datastore-backup:1.0.0
              env:
                #            - name: CLOUD_BILL_DATASTORE_BACKUP_GCS_BUCKET
                #              value: "gs://bucket"
                #            - name: CLOUD_BILL_DATASTORE_BACKUP_GCP_PROJECT_ID
                #              value: "<yourprojectid>"
                - name: GOOGLE_APPLICATION_CREDENTIALS
                  value: /auth/gcp-service-account/gcp-service-account.json
                - name: CLOUD_BILL_DATASTORE_BACKUP_CONFIG_FILE
                  value: /auth/datastore-backup-config/datastore-backup-config.json
              volumeMounts:
                - name: gcp-service-account
                  mountPath: "/auth/gcp-service-account"
                  readOnly: true
                - name: datastore-backup-config
                  mountPath: "/auth/datastore-backup-config"
                  readOnly: true
          restartPolicy: Never
          volumes:
            - name: gcp-service-account
              secret:
                secretName: gcp-service-account
            - name: datastore-backup-config
              secret:
                secretName: datastore-backup-config
```

## GCP Service Accounts
The pubsub service requires setting the environment variable **GOOGLE_APPLICATION_CREDENTIALS**. This is the path to your GCP service account credentials.

The following roles are required:
* Cloud Import Export Admin - Used to export from Cloud Datastore.
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
docker build -t datastore-backup:<tag> .

ex.
docker build -t datastore-backup:1 .
```

## Pushing to GCR
```
docker tag datastore-backup:<tag> gcr.io/<path>/datastore-backup:<tag>

docker push gcr.io/<path>/datastore-backup:<tag>

ex.
docker tag datastore-backup:1 gcr.io/cloud-bill-dev/datastore-backup:1

docker push gcr.io/cloud-bill-dev/datastore-backup:1

```

## Running the docker image locally with environment variables
```
docker run -it --rm -e CLOUD_BILL_DATASTORE_BACKUP_GCS_BUCKET=gs://bucket -e CLOUD_BILL_SAAS_GCP_PROJECT_ID='gcp-project-1' --name my-datastore-backup datastore-backup:<tag>

```
