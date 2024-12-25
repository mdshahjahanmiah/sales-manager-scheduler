# Sales Manager Scheduler
A tool for efficiently managing and scheduling sales team activities and appointments.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Setup and Run](#setup-and-run)
- [Setup and Tests](#setup-and-tests)
- [Features](#features)
- [Further Improvements](#further-improvements)

## Prerequisites

Ensure the following are installed on your system:
1. **Docker**: You can download and install Docker from [Docker's official website](https://www.docker.com/products/docker-desktop).
2. **Docker Compose**: Docker Compose usually comes bundled with Docker Desktop. Verify your installation by running:
```sh
   docker-compose --version
```

## Setup and Run 
Navigate to the project root where the docker-compose.yml file is located.
   ```sh
   cd sales-manager-scheduler
   ```
#### Run the following commands:
   ```sh
   docker-compose down
   docker-compose up --build
   ```
If the build is successful, you will see logs similar to:
   ```sh
[+] Running 2/0
 ✔ Container db      Created                                                                                                                                                                                                                            0.0s 
 ✔ Container app     Created                                                                                                                                                                                                                          0.1s 
Attaching to app, db
db   | 2024-09-24 15:15:33.501 UTC [1] LOG:  database system is ready to accept connections
app  | 2024/09/24 15:15:33 INFO configuration is loaded successfully
app  | 2024/09/24 15:15:33 INFO logger is initialized successfully
app  | 2024/09/24 15:15:33 INFO di container is starting up
app  | 2024/09/24 15:15:33 INFO http server is started successfully addr=0.0.0.0:3000

   ```
## Setup and Tests
To test the implementation, a test application as provided is included in the project directory. Ensure Node.js is installed.
Navigate to the project root, locate the tests folder, and enter it -
   ```sh
cd sales-manager-scheduler/tests
   ```

#### Install dependencies:
   ```sh
npm install
   ```
#### Run the tests:
   ```sh
npm run test
   ```
Upon successful test execution, you should see output similar to:
   ```sh
 PASS  ./test.js
  Coding challenge calendar tests
    ✓ Monday 2024-05-03, Solar Panels and Heatpumps, German and Gold customer. Only Seller 2 is selectable. (49 ms)
    ✓ Monday 2024-05-03, Heatpumps, English and Silver customer. Both Seller 2 and Seller 3 are selectable. (2 ms)
    ✓ Monday 2024-05-03, SolarPanels, German and Bronze customer. All Seller 1 and 2 are selectable, but Seller 1 does not have available slots. (2 ms)
    ✓ Tuesday 2024-05-04, Solar Panels and Heatpumps, German and Gold customer. Only Seller 2 is selectable, but it is fully booked (1 ms)
    ✓ Tuesday 2024-05-04, Heatpumps, English and Silver customer. Both Seller 2 and Seller 3 are selectable, but Seller 2 is fully booked. (1 ms)
    ✓ Monday 2024-05-03, SolarPanels, German and Bronze customer. Seller 1 and 2 are selectable, but Seller 2 is fully booked (2 ms)

Test Suites: 1 passed, 1 total
Tests:       6 passed, 6 total
Snapshots:   0 total
Time:        0.201 s
Ran all test suites.
   ```
## Features
- **Slot Availability**: The system focuses on returning available slots based on matching sales managers, ensuring efficient scheduling for the sales team.
## Further Improvements
- **Database Query Optimization**: We can reduce the number of database calls, ensure that the database tables have appropriate indexing to speed up queries, especially on columns used in where, order by or join clauses. If the performance bottleneck persists, an alternative and efficient approach could be using views and/or materialized view.
- **Object-Relational Mapping (ORM) Solution**: We can also consider object-relational mapping solution instead database first approach. However, we might need to consider performance overhead specially when need to deal with complex queries.
- **Dockerfile for Multiple Environments**: Multi-Stage Builds, Environment-Specific configuration. 
- **Logging and Monitoring**: Distributed Tracing, Metrics and Alerts.
- **Robust E2E Testing**: As provided test client, I haven't focused on the testing scenarios extensively. However, it's essential to consider potential flaky test cases and implement advanced testing scenarios to ensure the system is production-ready.
- **Connection Pooling**: Consider using database connection pooling instead of opening a new database connection unnecessarily.
