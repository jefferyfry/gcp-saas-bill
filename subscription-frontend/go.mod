module github.com/cloudbees/jenkins-support-saas/frontend-service

replace github.com/cloudbees/jenkins-support-saas/template => ../template

go 1.12

require (
	github.com/cloudbees/jenkins-support-saas/template v0.0.0-00010101000000-000000000000
	github.com/coreos/go-oidc v2.1.0+incompatible
	github.com/gorilla/mux v1.7.2
	github.com/gorilla/sessions v1.2.0
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20190829043050-9756ffdc2472 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
)
