# Hotel reservation backend

## Project environment variables
```
HTTP_LISTEN_ADDRESS=:3000
JWT_SECRET=somethingsupersecretthatNOBODYKNOWS
MONGO_DB_NAME=hotel-reservation
MONGO_DB_URL=mongodb://localhost:27017
MONGO_DB_URL_TEST=mongodb://localhost:27017
```

## Project outline
- customers -> book rooms, get reservations
- admins -> manage rooms and reservations
- Authentication and authorization -> JWT tokens
- Rooms -> CRUD API -> JSON
- Reservations -> CRUD API -> JSON
- Scripts -> database management -> seeding, migration

## Installation

- Clone the repository: `git clone https://github.com/HMZElidrissi/hotel-reservation-system-api.git`
- Define the environment variables in a `.env` file: `cp .env.example .env`
- Run `make deps` to install dependencies
- Run `make run` to start the server
- To run tests, run `make test`

## Technologies Used
### Mongodb driver
Documentation
```
https://mongodb.com/docs/drivers/go/current/quick-start
```

Installing mongodb client
```
go get go.mongodb.org/mongo-driver/mongo
```

### Gin Web Framework
Documentation
```
https://gin-gonic.com
```

Installing gin
```
go get github.com/gin-gonic/gin
```

## To DO

- [x] Test Customer handlers
- [ ] Test Admin reservation handlers
- [ ] Test Admin room handlers
- [ ] Scripts -> database management -> seeding, migration
