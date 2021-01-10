# user-service
user service using grpc go, postgresql and redis. The service is designed to be accessed by the internal network so that the grpc connection used is an insecure connection.

## Get Started
- git clone git@github.com:jacky-htg/user-service.git
- make init
- cp .env.example .env (and edit with your environment)
- make migrate
- make seed
- make server

## Features
- [ ] Companies
- [ ] Regions
- [ ] Branches
- [ ] Employees
- [ ] Company Features
- [X] Users
- [ ] Groups
- [X] Auths
- [ ] Role Base Access Control (RBAC)

### Companies
- [X] Multi companies
- [X] Company registration
- [ ] Companies CRUD

### Regions
- [ ] Multi Regions
- [ ] One region can be assigned to many branches.
- [ ] Regions CRUD

### Branches
- [ ] Multi Branches
- [ ] Branches CRUD

### Employees
- [ ] Employees CRUD

### Auths
- [X] Login
- [X] Forgot Password
- [X] Reset Password
- [X] Change Password
- [X] Check Authorization 

### Users, Groups, Access and RBAC
- [X] Users CRUD
- [ ] Group CRUD
- [ ] Access CRUD
- [ ] Multi users, multi roles and multi access
- [ ] One role can be assigned multi access
- [ ] Role Base Access Control (RBAC)

### Features
- [ ] CRUD Features
- [ ] CRUD Pacakage Feature
- [X] Company Feature Setting : The company can use the whole of features, or cherry pick part of features.

## License
[The license of application is GPL-3.0](https://github.com/jacky-htg/user-service/blob/main/LICENSE), You can use this apllication for commercial use, distribution or modification. But there is no liability and warranty. Please read the license details carefully.