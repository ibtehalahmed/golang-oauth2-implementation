# Golang-OAuth2

built on top of [Gin framework](https://github.com/gin-gonic/gin) with code arechitecture according to [Model-View-Controller (MVC)](https://en.wikipedia.org/wiki/Model–view–controller) pattern.

## Run instructions for development

### Prerequisites

    Golang latest version 1.14.2
    Git
    Postgres database
    Create your development database
        sudo su - postgres
        createuser --interactive --pwprompt; (role name oauth2 and password pwd)
        createdb -O oauth2 oauth2_dev;

    Seed database
        go run seed/seeder.go 

### Runing

    Clone the repository in $GOPATH/src/gitlab.com/coreborn
    Run the following commands
        cd $GOPATH/src/github.com/ibtehalahmed/golang-oauth2-implementation
        go run main.go
