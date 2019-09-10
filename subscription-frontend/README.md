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