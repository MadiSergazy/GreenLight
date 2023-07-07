for running

go run ./cmd/api  

curl -i localhost:4000/v1/healthcheck

go run ./cmd/api  -port=3000 -env=production

see
http://localhost:4000/v1/healthcheck

Validating JSON Input
mkdir internal/validator
touch internal/validator/validator.go

BODY='{"title":"","year":1000,"runtime":"-123 mins","genres":["sci-fi","sci-fi"]}'

curl -i -d "$BODY" localhost:4000/v1/movies

for checking endpoints 
curl -i localhost:4000/v1/healthcheck

curl -X POST localhost:4000/v1/movies

curl localhost:4000/v1/movies/123


# Create a BODY variable containing the JSON data that we want to send.
BODY='{"title":"Moana","year": 2016, "runtime": 107, "genres": ["animation", "adventure"]}'
# Use the -d flag to send the contents of the BODY variable as the HTTP request body.
curl -i -d "$BODY" localhost:4000/v1/movies

deleting movie
curl -X DELETE localhost:4000/v1/movies/5

Advanced CRUD Operations

curl -X PATCH -d '{"year": 1985}' localhost:4000/v1/movies/4 &
curl -X PATCH -d '{"year": 1986}' localhost:4000/v1/movies/4 &

{           
        "movie": {
		...
        }
}
{
        "error": "unable to update the record due to an edit conflict, please try again"
}

Sorting Lists
// Sort the movies on the title field in ascending alphabetical order. curl "http://localhost:4000/v1/movies?sort=title"

// Sort the movies on the year field in descending numerical order. curl "http://localhost:4000/v1/movies?sort=-year"

Paginating Lists
// Return the 5 records on page 1 (records 1-5 in the dataset) curl "http://localhost:4000/v1/movies?page=1&page_size=5"

// Return the next 5 records on page 2 (records 6-10 in the dataset) curl "http://localhost:4000/v1/movies?page=2&page_size=5"

Structured JSON Log Entries
mkdir internal/jsonlog
touch internal/jsonlog/jsonlog.go

greenlight git:(main) ✗ go run ./cmd/api

{"level":"INFO","time":"2021-06-23T15:06:37Z","message":"database connection pool established"}
{"level":"INFO","time":"2021-06-23T15:06:37Z","message":"database migrations applied"}
{"level":"INFO","time":"2021-06-23T15:06:37Z","message":"Starting %s server on %s","properties":{"addr":":4000","env":"development"}}
{"level":"FATAL","time":"2021-06-23T15:06:37Z","message":"listen tcp :4000: bind: address already in use"}


Sending Shutdown Signals
pgrep -l api

pkill -SIGKILL api
pkill -SIGTERM api