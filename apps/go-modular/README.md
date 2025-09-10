# Go Application

## Scaffold New Module
```sh
mkdir -p apps/go-modular/modules/dummy
mkdir -p apps/go-modular/modules/dummy/handler
mkdir -p apps/go-modular/modules/dummy/models
mkdir -p apps/go-modular/modules/dummy/repository
mkdir -p apps/go-modular/modules/dummy/services
echo 'package handler' > apps/go-modular/modules/dummy/handler/handler.go
echo 'package models' > apps/go-modular/modules/dummy/models/model.go
echo 'package models' > apps/go-modular/modules/dummy/models/schema.go
echo 'package repository' > apps/go-modular/modules/dummy/repository/repository.go
echo 'package services' > apps/go-modular/modules/dummy/services/services.go
echo 'package dummy' > apps/go-modular/modules/dummy/module.go
```
