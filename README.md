# M800 Backend Homework

## Requirements
- RESTful (JSON format) API
- mongoDB basic operation
- code push to github with proper commit

use below golang libs
- HTTP framework: https://github.com/gin-gonic/gin
- Config: https://github.com/spf13/viper
- mongo driver: https://github.com/mongodb/mongo-go-driver
- cobra: command line tools https://github.com/spf13/cobra

    - commit #1 setup project
    - commit #2 Makefile or a script for local setup and run MongoDB docker (version: 4.4)
    - commit #3 setup necessary config of LINE, MongoDB 
        - Line official account message integration (use go line sdk),
        - Create a test line dev official account
    - commit #4 Create a Go package connect to mongoDB, create a model/DTO to save/query user message to MongoDB
    - commit #5 Create a Gin API
        - receive message from line webhook, save the user info and message in MongoDB
        - (Hint: using ngrok for local test to generate a https endpoint)
    - commit #6 Create a API send message back to line
    - commit #7 Create a API query message list of the user from MongoDB
    - provide a demo video or steps of test (or postman)