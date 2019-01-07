# go-db-data

This project created to make query MySQl (Using <bold> GoLang and GoLang template </bold>)
    - If there are multiple schema and data is distributed around the instances.
    
Below you will find some information on how to perform common tasks.<br>

## Steps to proceed -
# Install GoLang 
   https://golang.org/doc/install

# Used GoLang Feature 
  #Channel
    https://tour.golang.org/concurrency/2
  # Template 
    https://golang.org/pkg/text/template/



#Steps to use
 1. Configure MySql Parent MySql details in conf.json
            {
            "host": "127.0.0.1",
            "port": "3306",
            "userName": "root",
            "password": "root",
            "db": "schema",
            "dbTable": "ip_table"
        }

    1A. host - MySql parent instance - which has information about all the schema IP's.
    1B . db - Default schema to connect to 
    1C= ip_table -> IP and schema information 

2. To run 
  go run main_http.go

3. To ceate binary 
        go build main_http.go
   
4. To create binary for different OS's 
    Linux example
    go build -o db_data.linux main_http.go

    GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o db_data_log.linux main_http.go