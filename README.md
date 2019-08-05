# CodersRank task

[![Build Status](https://dev.azure.com/nyckowskijakub/codersranktask/_apis/build/status/jakule.codersranktask?branchName=master)](https://dev.azure.com/nyckowskijakub/codersranktask/_build/latest?definitionId=1&branchName=master)

This is my solution to CodersRank task. Application is written in Go and 
it's deployed to Heroku https://codersranktask.herokuapp.com/. As a backend 
I decided to use Postgres. I also decided to configure Azure DevOps pipeline 
to build and auto deploy that project to Heroku.

#### Build
Go 1.12 with modules enabled is required

```bash
export GO111MODULE=on
go build
```

#### Build image

````bash
docker build -t codersranktask .
````

#### Development/test environment

To run all dependencies locally `docker-compose` script can be used:

```bash
docker-compose up
``` 

##### TODO:

- [ ] Add system tests
- [ ] UT not covered paths
- [ ] Add encryption to secretes
- [ ] Add documentation