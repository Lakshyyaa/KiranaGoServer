# Kirana Club Backend assignment documentation


## Description:
- This server processes jobs and images collected, calculating simulated dimensions with and GPU processing times.
- It provides RESTful endpoints for job submission and status checking with request validation and checks like (storeId exists, jobId exists, request body format is correct, etc)
- To handle multiple requests simultaneously, goroutines are utilized to handle concurrency and used locks to handle race conditions(mutex and wait groups)
- Each job status is committed and updated in real-time to have latest updates


## Assumptions:
- Since no CRUD/database operations were requested, I opted to use in-memory data structures, which work well for this assignment but should never be used in real-life
- As the JobID was not specified for each job, I chose to implement a global counter for generating unique IDs for each
- Downloading the images would cause latency(depending on image size) so I opted to simulate image
downloading and chose random parameters
- When creating a job in submit endpoint, I decided to send 201 Created status even for requests
whose store_id did not exists in the csv file, which though logically wrong, was only mentioned to be done using the second endpoint(status). I just have to check it in the storeMaster map for it.
- Since it was not implicit if we should send the job_id and error fields as blank or with the actual
data in the http responses, I decided to share the actual data. Sending blank will only require
changing struct definitions.


## Development Environment:
- **Operating System**: Windows 11 Home 64-bit
- **Hardware**: AMD Ryzen 5 5500U
- **IDE**: Visual Studio Code
- **Language**: Go 1.23.2
- **Libraries**:
- `gorilla/mux` for routing
- `air-verse/air` for live-reloading
- `net/http` for HTTP server
- `encoding/json` for JSON handling
- `encoding/csv` for CSV processing
- `sync` for concurrency control


## Workflow:
- main.go file initializes the server and loads the csv file into the storeMaster map
- router.go file creates the router with 2 endpoints and their handler functions
- controllers.fo has 2 handlers functions, one for each endpoint and all data structures and
  variable being used
- utils.go has functions for job and image processing/simulation
- the 1st function for submission validates the request, assigns it a jobId, stores its status and calls
  a goroutine for processing the job
- the job processing request calls the goroutine to do image simulation(gpu) and does validation
  and updates the job status
- the 2nd function for status checks job status and sends back a suitable response
- all code that accesses critical section uses locks to handle race conditions and for goroutines wait groups are used


## Improvements:
- Currently using in-memory data structures/variables. One things that would help is have a persistent database to store data in
- Currently gpu is simulated which simply logs to the console and in the end the job-status is committed as it is. If actually working with image data, to store it we could have a channel which reads data of each image from its goroutine(each image processing is done in a separate goroutine) and read from it and store image data where needed. This would be better than waiting for all images' processing to finish and then store them in the end like it would happen currently(wg.Wait() blocks committing the final job-status after image gpu simulation is done for all images of a request which could be 1000s in real-life)


## Installation:
- 