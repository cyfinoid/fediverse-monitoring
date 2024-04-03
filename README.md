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

## Requirements:
- Step up PostgresSQL in your local environment (Change constants like username, password, Database name in DB/db.go at Line 12)
  ```
  const (
	  host     = "localhost"
	  port     = 5432
	  user     = "postgres"
	  password = "postgres"
	  dbname	 = "dailygrowth"
  )
  ``` 


## Setup 
- If running for the first time or want to update database, update Fediverse database by
  ```
  go run fediverse/main.go
  ```
  Database is being updated by cocurrently gathering up data through API calls. Semaphores are being used (a concurrency control mechanism that limits the number of threads that can 
  access a resource or a group of resources concurrently, in this case 25 goroutines are maintained) to collect API data and store in postgres.
  *Number of goroutines spawned can be customized at Line 60 in Fediverse/main.go*
  ```
  numWorkers := 25 //no. of goroutines
  ```
  *Average time taken to update the database varies from 0.5 hour to 1 hour (depending on the number of goroutines spawned). NOTE: Keep "numWorkers" up to the limit that it doesn't cause IO overhead while making API calls*
  
- Run the server and open http://localhost:8080/
  ```
  go run server.go
  ```

## Screenshots
![fedi_stats](/ss/1.png?raw=true "fedi_stats 1")
![fedi_stats](/ss/2.png?raw=true "fedi_stats 2")

## Languages and frameworks used:
- Golang
- HTML
- CSS
- Chart.js
- PostgresSQL

### Contributors

- Vansh Bulani (Internship project at Cyfinoid Research) \n
  [Checkout my portfolio](https://www.vanshbulani.info) \n
  [Checkout my blog about this project](https://www.vanshbulani.info/blogs) \n
  [LinkedIn](https://www.linkedin.com/in/vanshbulani/)
