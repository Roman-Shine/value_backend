# value_backend
Experimental api on golang

Create .env file in root directory:
```dotenv
APP_ENV=local

MONGO_URI=mongodb://mongodb:27017
MONGO_USER=admin
MONGO_PASS=qwerty

PASSWORD_SALT=
JWT_SIGNING_KEY=

HTTP_HOST=localhost
```

Use `make run` to build&run project


If you are using Linux, uncomment line 7 in the Makefile and comment out line 6