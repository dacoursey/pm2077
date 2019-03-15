# pm2077
Vulnerable go application to demo Go specific web security issues

# Application

pm2077 is a project management application. Currently available data types are Customers, Contacts, Projects.  More will be added as we expand.

It utilizes the Go stdlib for most functionality including go html templates. https://github.com/gorilla/mux is used for routing.

# Setup 

First we need a the PostgreSQL server, you can go with the Docker image or the local install.

*Docker Image*:

From the command line run:

```
docker run -d --name gonv -p 5432:5432 -v /Users:/host/users -e "POSTGRES_USER=nvuser" -e "POSTGRES_DB=gonv" postgres:latest
docker exec -it gonv /bin/bash
psql -d gonv -U nvuser
\i /host/users/path/to/schema.sql
```

*Local Install*:

From the command line run:

```
brew install postgres
pg_ctl -D /usr/local/var/postgres start && brew services start postgresql
psql postgres
```

Once in the psql cli set up the blank database:

```
CREATE ROLE nvuser WITH LOGIN PASSWORD 'nvuser';
ALTER ROLE nvuser CREATEDB;
CREATE DATABASE gonv;
GRANT ALL PRIVILEGES ON DATABASE gonv TO nvuser;
\q
```

Now log back in as nvuser and import the schema:

```
psql -d gonv -U nvuser
\i /go/work/path/src/github.com/user/pm2077/db/schema.sql
\list
\dt
\q
```

Now we can get started with pm2077 code.

```
cd $GOPATH
go get github.com/dacoursey/pm2077
go install github.com/dacoursey/pm2077
$GOPATH/bin/pm2077
```

You should have the server running now. Open a browser and load http://localhost:5000 by default. This can be changed in `main.go` if needed.

If you get panics loading the site, you may need to create symlinks in the Go bin directory to the `/public` and the `/templates` directories in main pm2077 path.
