# Jenkins Support Subscription Frontend Service
The Frontend service provides the UI for customer signup from the marketplace. The end result is storing the customer account information for the subscription and confirming the account with Google. Auth0 and Google Identity are used to capture some of the customer profile data.

## Frontend Flow
The basic frontend flow amongst handlers and pages is the following:
![Jenkins Support SaaS - Page 4](https://user-images.githubusercontent.com/6440106/64573203-54b36280-d31f-11e9-84cb-9e0ca4e5fc67.png)

## Pages
* [signup.html](https://github.com/cloudbees/jenkins-support-saas/tree/master/subscription-frontend/templates/signup.html) - Initial page to direct customer to Auth0/Google sign in. The customer is sent to this page from marketplace.
* [confirm.html](https://github.com/cloudbees/jenkins-support-saas/tree/master/subscription-frontend/templates/confirm.html) - Auth0/Google callback page to confirm account information.
* [finish.html](https://github.com/cloudbees/jenkins-support-saas/tree/master/subscription-frontend/templates/finish.html) - Final page to confirm account creation and notify customer of next steps.

