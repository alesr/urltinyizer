# URLTINYIZER

URLTINYIZER is a REST API written in Go that provides the following functionality:

### Endpoint for creating short url
A POST request to /shorten endpoint with a JSON payload containing the long_url as a string returns a shortened url.

### Endpoint for redirecting users
A GET request to /{shortURL} redirects the user to the original long url and increments the number of hits.

### Stats endpoint
A GET request to /{shortURL}/stats returns the number of times a short url has been used.


The application runs on two Docker containers: one for the PostgreSQL database and the other for the application itself. To run the application, simply run make run.
