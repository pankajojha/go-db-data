# go-db-data

This project created to make query MySQl (Using <bold> GoLang and GoLang template </bold>)
    - If there are multiple schema and data is distributed around the instances.
    - Below you will find some information on how to perform common tasks.<br>

## Steps to proceed -
 1. Configure MySql Parent MySql details in conf.json
            {
            "host": "127.0.0.1",
            "port": "3306",
            "userName": "root",
            "password": "root",
            "db": "schema",
            "dbTable": "ip_table"
        }

 2. Run on Mac
   ./build/db-tool.mac

 3. On Linux
   ./build/db-tool.linux

 4. Access to run query on multiple wf's
    http://localhost:3000/query



For development contribution 
# Install GoLang    https://golang.org/doc/install