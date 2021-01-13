# user-service
user service using grpc go, postgresql and redis. The service is designed to be accessed by the internal network so that the grpc connection used is an insecure connection.

## Get Started
- git clone git@github.com:jacky-htg/user-service.git
- make init
- cp .env.example .env (and edit with your environment)
- make migrate
- make seed
- make server
- You can test the service using `go run client/main.go` and select the test case on file client/main.go

## Features
- [X] Companies
- [X] Regions
- [X] Branches
- [X] Employees
- [X] Company Features
- [X] Users
- [X] Groups
- [X] Auths
- [X] Role Base Access Control (RBAC)

### Companies
- [X] Multi companies
- [X] Company registration
- [X] Companies CRUD

### Regions
- [X] Multi Regions
- [X] One region can be assigned to many branches.
- [X] Regions CRUD

### Branches
- [X] Multi Branches
- [X] Branches CRUD

### Employees
- [X] Employees CRUD

### Auths
- [X] Login
- [X] Forgot Password
- [X] Reset Password
- [X] Change Password
- [X] Check Authorization 

### Users, Groups, Access and RBAC
- [X] Users CRUD
- [X] Group CRUD
- [X] List Access
- [X] Multi users
- [X] One role can be assigned multi access
- [X] Role Base Access Control (RBAC)

### Features
- [X] List Features
- [X] List Package Feature
- [X] View Package Feature
- [X] Company Feature Setting : The company can use the whole of features, or cherry pick part of features.

## How To Contrubute
- Give star or clone and fork the repository
- Report the bug
- Submit issue for request of enhancement
- Pull Request for fixing bug or enhancement module 

## License
[The license of application is GPL-3.0](https://github.com/jacky-htg/user-service/blob/main/LICENSE), You can use this apllication for commercial use, distribution or modification. But there is no liability and warranty. Please read the license details carefully.

## Link Repository
- [API Gateway for Inventory](https://github.com/jacky-htg/api-gateway-service) -- inprogress development
- [Simple gRPC Skeleton](https://github.com/jacky-htg/grpc-skeleton)