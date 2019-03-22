# Project details:

A REST API implmentation using golang.

Create a POST API endpoint to create new record in table, which accepts content type application-json body like below
```sh
 POST http://localhost:8585/login
   content-type : application-json
   body :
    {
      "first_name" : "String",
      "last_name" : "String",
      "organization" : "String",
      "phone_number" : "String",
      "email" : "String",
      "website" : "String"
    }
 ```

Create a PUT API endpoint to update particular row given by recordId, which accepts id(Record ID) and content type application-json body like below
```sh
 PUT http://localhost:8585/login/id
   content-type : application-json
   body :
    {
      "first_name" : "String",
      "last_name" : "String",
      "organization" : "String",
      "phone_number" : "String",
      "email" : "String",
      "website" : "String"
    }
 ```

Create a GET API endpoint to get all the records from table.
```sh
 GET http://localhost:8585/logins
```

Create a GET API endpoint to get particular record from table, which accepts id(Record ID) like below
```sh
 GET http://localhost:8585/login/id
```

The above API endpoints should return response on success, and in case of failure an error object will be expected. All the responses should have a relevant http status code, like - 200(StatusOK) and 201(StatusCreated) for success , 400(StatusBadRequest), 500(StatusInternalServerError) and 501(StatusNotImplemented) for errors.


### The Project Structure
1. main.go - The main point of the project.
2. The dsmartsys/sqlinterface consists of: 
    - sqlinterface.go - Contains logic for MySql database interface.
3. The dsmartsys/util consists of:
    - util.go - Contains logic for exceptional handling.
4. The dsmartsys/vendor consists of:
    - vendor.json - Contains all third party packages names and their revisions we use for our need. On fly we can download these packages in our machines using govendor tool.


### Instructions for build and run the REST-API application
Step 1 - Open the Terminal <br />
Step 2 - Go to the stored project source directory <br />
Step 3 - If GOBIN is not set, we need to set the GOBIN <br />
Step 4 - Set GOPATH evironment variable which determines the location of this project workspace <br />
Step 5 - Install govendor ```go get -u github.com/kardianos/govendor``` <br />
Step 6 - GOBIN(govendor) path should be added to PATH environment variable <br />
Step 7 - cd "my project in GOPATH" <br />
Step 8 - Run ```govendor sync``` command for downloading all the third party packages that we are using in this project <br />
Step 9 - Run ```go build main.go``` <br />
Step 10 - Run ```./main "MySqlUserName" "MySqlPawd" "MySqlIPAddress" "MySqlPort"``` <br />
Step 11 - Default MySql user name will be root <br />
Step 12 - Default MySql ip address will be 127.0.0.1 <br />
Step 13 - Default MySql port address will be 3306 <br />


### For more information about govendor tool please go through the below link:
[govendor](https://github.com/kardianos/govendor)