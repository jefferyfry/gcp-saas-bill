# Jenkins Support Subscription Frontend Service
The Frontend service provides the UI for customer signup from the marketplace. The end result is storing the customer account information for the subscription and confirming the account with Google. Auth0 and Google Identity are used to capture some of the customer profile data.

## Frontend Flow
The basic frontend flow amongst handlers and pages is the following:
![Jenkins Support SaaS - Page 4](https://user-images.githubusercontent.com/6440106/64573203-54b36280-d31f-11e9-84cb-9e0ca4e5fc67.png)

## Handler Functions
* [Signup](https://github.com/cloudbees/jenkins-support-saas/blob/master/subscription-frontend/web/handlers.go#L63)
* [Auth0Login](https://github.com/cloudbees/jenkins-support-saas/blob/master/subscription-frontend/web/handlers.go#L162)
* [Auth0Callback](https://github.com/cloudbees/jenkins-support-saas/blob/master/subscription-frontend/web/handlers.go#L194)
* [Finish](https://github.com/cloudbees/jenkins-support-saas/blob/master/subscription-frontend/web/handlers.go#L254)

## Pages
* [signup.html](https://github.com/cloudbees/jenkins-support-saas/tree/master/subscription-frontend/templates/signup.html) - Initial page to direct customer to Auth0/Google sign in. The customer is sent to this page from marketplace.
* [confirm.html](https://github.com/cloudbees/jenkins-support-saas/tree/master/subscription-frontend/templates/confirm.html) - Auth0/Google callback page to confirm account information.
* [finish.html](https://github.com/cloudbees/jenkins-support-saas/tree/master/subscription-frontend/templates/finish.html) - Final page to confirm account creation and notify customer of next steps.

## Configuration
To successfully run the subscription service, configuration must be set through either environment variables, command-line options or a configuration file. You may chose an option based on on your intent (development, testing, production deployment). The following configuration is required:

* Frontend Service Endpoint - This is the listening port for the frontend service.
* Subscription Service URL - This is the URL to the subscription service.
* Client ID - This is the Oauth/Auth0 client ID.
* Client Secret - This is the Oauth/Auth0 client secret.
* Callback URL - This is the callback URL used by the Oauth/Auth0 service.
* Issuer - The Oauth issuer or Auth0 domain.
* Session Key - A random character sequence for session encoding.

### Configuration Precedence
command-line options > environment variables

### Environment Variables
* JENKINS_SUPPORT_FRONTEND_CONFIG_FILE - Path to a configuration file (see below).
* JENKINS_SUPPORT_FRONTEND_SERVICE_ENDPOINT 
* JENKINS_SUPPORT_SUBSCRIPTION_SERVICE_URL 
* JENKINS_SUPPORT_FRONTEND_CLIENT_ID 
* JENKINS_SUPPORT_FRONTEND_CLIENT_SECRET 
* JENKINS_SUPPORT_FRONTEND_CALLBACK_URL
* JENKINS_SUPPORT_FRONTEND_ISSUER
* JENKINS_SUPPORT_FRONTEND_SESSION_KEY

### Command-Line Options
* configFile - Path to a configuration file (see below).
* frontendServiceEndpoint 
* subscriptionServiceUrl 
* clientId 
* clientSecret 
* callbackUrl 
* issuer 
* sessionKey 

### Configuration File
The configFile command-line option or JENKINS_SUPPORT_FRONTEND_CONFIG_FILE environment variable requires a path to a JSON file with the configuration. Example:
```
{
  "frontendServiceEndpoint": "8086",
  "subscriptionServiceUrl": "http://subscription-service.default.svc.cluster.local:8085/api/v1",
  "clientId": "1234556",
  "clientSecret": "abcdefg",
  "callbackUrl": "https://cloudbees-jenkins-support.35.231.106.233.xip.io/callback",
  "issuer": "https://cloudbees-dev1.auth0.com/",
  "sessionKey": "somekeycloudbeesjen0udb33s1"
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
          image: gcr.io/cloudbees-jenkins-support-dev/frontend-service:latest
          env:
            - name: JENKINS_SUPPORT_FRONTEND_CONFIG_FILE
              value: /auth/frontend-service-config/frontend-service-config.json
          ports:
            - containerPort: 8086
          volumeMounts:
            - name: frontend-service-config
              mountPath: "/auth/frontend-service-config"
              readOnly: true
      volumes:
        - name: frontend-service-config
          secret:
            secretName: frontend-service-config
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
docker tag frontend-service:1 gcr.io/cloudbees-jenkins-support-dev/frontend-service:1

docker push gcr.io/cloudbees-jenkins-support-dev/frontend-service:1
```

## Running the docker image locally with environment variables
```
docker run -it --rm -p 8086:8086 -e JENKINS_SUPPORT_SUB_FRONTEND_SERVICE_ENDPOINT=8086 -e JENKINS_SUPPORT_SUB_SERVICE_URL='http://localhost:8085' -e JENKINS_SUPPORT_SUB_FRONTEND_CLIENT_ID='abcdef' -e JENKINS_SUPPORT_SUB_FRONTEND_CLIENT_SECRET='123456' -e JENKINS_SUPPORT_SUB_FRONTEND_CALLBACK_URL='http://localhost:8085/callback' -e JENKINS_SUPPORT_SUB_FRONTEND_ISSUER='issuer' -e JENKINS_SUPPORT_SUB_FRONTEND_SESSION_KEY='somekeycloudbeesjenkinssupportsessionkey1cl0udb33s1' --name my-frontend-service frontend-service-1:<tag>

```