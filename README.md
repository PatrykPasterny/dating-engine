# dating-engine
Service implementing four basic endpoints for dating app engine. 
It uses gRPC to enable retrieving the data using proto buffers.

### Usage
To run the service with all its dependencies clone the source code to your local machine using:
```
git clone <clonning_url>
```

Then use docker compose command to build all dependencies:
```
docker compose up --build
```

Once the service is running you can see that the logs appear inside the containers.
You can connect to the service by connecting to:
```
localhost:8080
```

### Testing
To test how the service work you can see the tests container that is running after 
docker compose call or run the tests locally once you set up the service with docker compose
command from Usage section.

To run the tests locally you need to set up the environment variables on your local machine using:
```
export BASE_URL=localhost:8080
export DATABASE_COLLECTION=matches
export DATABASE_NAME=db
export DATABASE_URI=mongodb://localhost:27017
```

and run the tests using IDE or:
```
go test ./...
```

### Decisions:
1. I decided to use MongoDB to store the data in it as the requirements are
to handle huge amounts of matches from whole span of users activity and for over 
10 mln users. It is then easier to scale the database than if we use sql database 
and the matches itself does not need to be structured. If we were to write the profiles
service as well, there I would consider using SQL DB.
2. I decided not to include Redis into the service dependencies. I instead implemented 
default pagination which size can be changed using service's configuration. It enables 
quick and performant querying of data, no matter the database size. Also the decision not 
to include Redis was based on the fact that the service enables grpc connection and won't be
used directly by frontend user, so if I were to write the Gateway service or Middleware service
then I would strongly consider adding Redis in that place.
3. I decided to not include the unit tests in the service itself as it is more or less a CRUD 
service right now that does not hold much of a logic and instead focus on integration tests in
which I treat the service as a black box. That way I am sure the outer connection to the service
works good and the endpoints correctly implement its logic.
4. I decided to handle matches from the actors side and recipient side separately. This increases
the complexity of the solution, but makes it more readable and easier to update than handling it in one object. 
This also makes the put operation a little bit more complex as it forces the changes on both actors and 
recipient entities. This forced quite wide transaction span which may be at some point a bottleneck, so it
would be worth running performance tests against the production data mirror before we decide to deploy it to production.