# What is this?
This is the REST API component of a three-parts project where an API, 
an [arduino](https://github.com/sqrrrrl/pc-power-arduino) 
and a [mobile app](https://github.com/sqrrrrl/pc-power-app) communicate together
in order to enable users to remotely control the power state of their desktop computer(s).

# Installation
1. Download and install [golang](https://go.dev/doc/install)
2. Clone the repository:
```
git clone https://github.com/sqrrrrl/pc-power-api.git
```
3. Navigate to the project directory and download the dependencies:
```
go mod download
```
4. Build the project:
```
go build -v -o ./outputDirectory/appName ./src/app
```

# Usage
## Environment variables
Some environment variables need to be set for the API to work properly:\
```PORT```: The port which the API should listen to\
```JWT_SECRET```: A random 32 characters string used to generate the authentication tokens\
```DBTYPE```: Either sqlite or mysql\
```DBNAME```: The name of the database\
\
The environment variables required to use mysql:\
```DBUSER```: The username of the database.\
```DBPASS```: The password of the database.\
```DBHOST```: The ip address or domain of the database.\
```DBPORT```: The port of the database.

## Starting the API
Once the environment is set the API can be started by executing the compiled project: 
```
./outputDirectory/appName
```

# License
Distributed under the [GPL-3.0 license](https://github.com/sqrrrrl/pc-power-app#GPL-3.0-1-ov-file)
