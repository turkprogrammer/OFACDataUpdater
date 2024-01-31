# OFACDataUpdater

**OFACDataUpdater** is a Go (Golang) application designed to manage and update data from the Specially Designated Nationals List (SDNList) of the U.S. Department of the Treasury's Office of Foreign Assets Control (OFAC). The project provides the following key features:

## Core Features

### Database Initialization

The application initializes a PostgreSQL database and creates a "people" table to store information about individuals.

### SDNList Data Loading

SDNList data about specially designated individuals is loaded into the application.

### Importing Data from OFAC

Data from SDNList is imported into the PostgreSQL database for further use.

### HTTP Server with Request Handlers

The application offers an HTTP server with various request handlers for operations such as data updates, retrieving name lists, checking data status, etc.

### Graceful Shutdown Handling

The application correctly handles termination signals (SIGINT or SIGTERM), ensuring the proper shutdown of the server upon receiving such signals.

# Usage

The application will be accessible at [http://localhost:8082](http://localhost:8082).

Use the provided API endpoints for data updates, name extraction, and data status checks.

## API Endpoints

### Data Update

- **Endpoint:** `/api/update`
- **Method:** `POST`
- **Description:** Initiates the OFAC data update process.

### Check Data Status

- **Endpoint:** `/api/state`
- **Method:** `GET`
- **Description:** Returns the current data status.

### Get Names with Strong Match

- **Endpoint:** `/api/get_names_strong`
- **Method:** `GET`
- **Query Parameter:** `name` (e.g., `/api/get_names_strong?name=JohnDoe`)
- **Description:** Returns a list of names with a strong match.

### Get Names with Weak Match

- **Endpoint:** `/api/get_names_weak`
- **Method:** `GET`
- **Query Parameter:** `name` (e.g., `/api/get_names_weak?name=JohnDoe`)
- **Description:** Returns a list of names with a weak match.

### Get Names

- **Endpoint:** `/api/get_names`
- **Method:** `GET`
- **Query Parameters:** `name`, `type` (e.g., `/api/get_names?name=JohnDoe&type=strong`)
- **Description:** Returns a list of names with a customizable match type.


To use **OFACDataUpdater**, follow these steps:

1. Set up the PostgreSQL database and provide the necessary configuration data in the "database" package.

2. Build the Docker image using the provided Dockerfile:

   ```bash
   docker build -t ofac-data-updater-image .

# Run

To run the application and the PostgreSQL container, use the following commands:

```bash
   docker-compose up