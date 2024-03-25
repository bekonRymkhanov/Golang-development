Name:

        the application for storing and updating the information for episodes of popular cartoon "pinguines from madagascar"

Description

        code works around app struct and it divided for 5 folders:

            ---midterm/bin directory will contain our compiled application binaries, ready for deployment
            to a production server.
            
            ---midterm/cmd/api  directory will contain the application-specific code for our Greenlight API
            application. This will include the code for running the server, reading and writing HTTP
            requests, and managing authentication.
            
            ---midterm/internal directory will contain various ancillary packages used by our API. It will
            contain the code for interacting with our database, doing data validation, sending emails
            and so on. Basically, any code which isn’t application-specific and can potentially be
            reused will live in here. Our Go code under cmd/api will import the packages in the
            internal directory (but never the other way around).
            
            ---midterm/migrations directory will contain the SQL migration files for our database.
            
            ---midterm/remote directory will contain the configuration files and setup scripts for our
            production server.(front end)
            
            
            ---midterm/go.mod file will declare our project dependencies, versions and module path.
            
            ---midterm/Makefile will contain recipes for automating common administrative tasks — like
            auditing our Go code, building binaries, and executing database migrations.
            
            ---midterm/internal/data will contain models of structs and episodes and data that used in our project
            
            ---midterm/internal/validator contains validators of strings from our project,mostly used in validation of input by user


API:

                POST /Episodes;
                GET /Episodes/:id;
                PUT /Episodes/:id;
                DELETE /Episodes/:id;
                HEALTHCHECK /healthCheck;

Database Structure :
    ![Alt text](../midterm/remote/dbschema.png?raw=true "Title")

    https://dbdiagram.io/d/65ec11b8b1f3d4062c801a4e

Badges:

        this project is overrided all go language errors by writing oun error handlers to not show many details to user and panic if some non-expected error happened,errors handler file can be found in midterm/cmd/api/errors.go 


Installation:

        go get github.com/julienschmidt/httprouter
        ...

Usage:

        libraries that where used in this project mostly:

            "context"
            "database/sql"
            "flag"
            "fmt"
            "log"
            "net/http"
            "os"

Support:
        
        telegram @illegalunicorn,  
        teams b_rymkhanov@kbtu.kz


template of message:

        im ...........
        writing about your project named .....,
        to ask you questions about ......,
        I tried to do ..... and it gives me the result of ....,
        .....................

Roadmap:

        v1.0.0 only supports routers and simple json in web page
        this is the project v2.0.0 but in endterm I will add more things like login ui/ux and e.t.c.

Contributing:

        I expect from readers to contribute in flexibility in this project by ability to change the behavior of code

        you can join my team in development of application untill the endterm exam,so we can make project more complex and detailed


Authors and acknowledgment:

        Rymkhanov Bekarys 22B030423


Project status:

        development stage stopped for 2 weeks due to lack of resourses for development (waiting for Airat agai`s lectures)