# fediverse-monitoring

The fediverse is an ensemble of social networks which can communicate with each other, while remaining independent platforms. Users on different social networks and websites can send and receive updates from others across the network. Nearly all fediverse platforms are free and open-source software.
Due to the open nature of protocols a lot of details are available if you check the right api.

This project gathers up the following statistics of fediverse:

- Top 10 server software
- Top 10 Server Software + Version
- Top 10 TLD (Top level domain)
- Daily server count
- Daily users count
- Daily posts count
- Daily comments count

## Setup 
As of now this project runs only on local environment.
- Step up PostgresSQL in your local environment
- Run the server
  ```
  go run server.go
  ```
- If running for the first time or want to update database, update Fediverse database by
  ```
  go run fediverse/main.go
  ```
  Database is being updated by cocurrently gathering up data through API calls. Semaphores are being used (a concurrency control mechanism that limits the number of threads that can 
  access a resource or a group of resources concurrently, in this case 25 goroutines are maintained) to collect API data and store in postgres.

### Languages and frameworks used:
- Golang
- HTML
- CSS
- Chart.js
- PostgresSQL
