# grpc-template-service

> copy from `https://github.com/juanjiTech/jframe` and modify somethings

### 1. overview
this is a grpc template service mixed with jinx,grpc and grpc-gateway.

you can access the route by `http://localhost:8080/v1/*any`

meanwhile, this framework is composed of mods, you can add your own mods to the framework.

### 2. Quick Start
create config
```shell
go run main.go -c config.yaml
```
create mod 
```shell
go run main.go -m modName
```
run server
```shell
go run main.go server  
```

### 3. mod
* grpc-gateway
* jinx
* jinPprof
* mysql
* pgsql
* ....
