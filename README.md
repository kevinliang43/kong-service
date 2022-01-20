# kong-service

## Setup
1. Install [PostgreSQL v14.1](https://www.postgresql.org/docs/14/index.html) and start up a PostgreSQL server.
2. Setup the Database and Tables using the SQL statements provided in https://github.com/kevinliang43/kong-service/blob/main/service_database.sql
3. In a separate Terminal tab, clone this repository and `cd` into the directory.
4. To start up `kong-service`, run `go run .`

Notes:
1. By default, the service will run on `localhost:8080` and will be expecting a PostgreSQL database to connect to at `localhost:5432`. If needed, feel free to modify the configs at https://github.com/kevinliang43/kong-service/blob/main/config/config.yaml prior to running `go run .`

## Models

1. [Service](https://github.com/kevinliang43/kong-service/blob/main/models/service.go#L11-L18)
  - Represents the latest version for a given Service, and the number of versions that this ServiceId is associated with.
2. [ServiceRecord](https://github.com/kevinliang43/kong-service/blob/main/models/service.go#L3-L9)
  - Represents a record for a given ServiceId and its Version.

## Endpoints

### Services
1. `GET /services/{serviceId}`
  - Retrieve a [Service](https://github.com/kevinliang43/kong-service/blob/main/models/service.go#L11-L18) by `ServiceId`
  - Sample Query/Response:

```
> curl -X GET -H "Content-Type: application/json" http://localhost:8080/services/1
> {
    "id": "cac11211-e3bd-488b-98df-a906306115ff",
    "serviceId": 1,
    "name": "NEW SERVICE",
    "description": "UPDATED",
    "version": 1.3,
    "versions": 4
}
```


2. `POST /services/search`
  - Provide a [ServiceSearchRequest](https://github.com/kevinliang43/kong-service/blob/main/models/service.go#L20-L26) as the POST body and retrieve a [ServiceSearchResponse](https://github.com/kevinliang43/kong-service/blob/main/models/service.go#L28-L31) that contains a paginated list of [Service](https://github.com/kevinliang43/kong-service/blob/main/models/service.go#L11-L18) and the `offset` for the next page
  - POST Body details:

```
ServicesSearchRequest {
	searchQuery string // Search Query. Response will contain services that match the provided 'searchQuery' for the given 'filterType'.
	limit       int64  // Integer representing the limit of the max number of 'Service' objects to return (0 <= limit < 100).
	offset      int64  // Integer representing the offset in which the paginated results will begin (0 <= offset)
	sortType    string // One of ['ASC', 'DESC']. Describes how the resultant set of 'Service' objects will be sorted.
	filterType  string // One of ['name', 'description']. Describes which column to search on with the provided 'searchQuery'. If not provided, defaults to 'name'.
}

ServicesSearchResponse {
	services   []Service // List of services that conform to the provided ServicesSearchRequest.
	NextOffset int64     // Offset to provide to the next 'ServicesSearchRequest::offset' for the next page of results.
}
```
  - Sample Request/Response:

```
> curl -X POST -H "Content-Type: application/json" --data '{"searchQuery": "Co", "limit": 3, "offset":0, "sortType":"DESC", "filterType": "name"}' http://localhost:8080/services/search
> {
    "services": [
        {
            "id": "1d75df8d-1639-428d-aad3-c78cd71a250f",
            "serviceId": 3,
            "name": "Contact Us",
            "description": "Service for retrieving contact us info",
            "version": 1,
            "versions": 1
        },
        {
            "id": "ad187f14-b774-4d25-a939-864aa57903a0",
            "serviceId": 2,
            "name": "Collect Monday",
            "description": "",
            "version": 1,
            "versions": 1
        }
    ],
    "nextOffset": 2
}
```
3. `POST /services`
  - Create a new [Service](https://github.com/kevinliang43/kong-service/blob/main/models/service.go#L11-L18) or a new version for an existing [Service](https://github.com/kevinliang43/kong-service/blob/main/models/service.go#L11-L18). Responds with the newly created `Service` or updated `Service` version. If the provided service version is lower than the most up to date version for an existing `Service`, then a 400 Bad Request will be returned.
  - POST body details:
```
{
	serviceId   int64   // OPTIONAL. If Provided, denotes a new Version for an existing Service. Otherwise, denotes a new Service being created (ServiceId will be generated under the hood).
	Name        string  // Name of the new Service / new Service version
	Description string  // Description of the new Service / new Service version
	Version     float64 // Version number of the new Service / new Service version. If for a new Service version, the version number must be greater than most up to date version that currently exists.
}
```
  - Sample Request/Response
```
// New Service
> curl -X POST -H "Content-Type: application/json" --data '{"name": "NEW SERVICE", "description": "NEW SERVICE", "version": 1.0}' http://localhost:8080/services
> {
    "id": "17cb6286-ad2d-46d7-ac04-9707f0f4447a",
    "serviceId": 21,
    "name": "NEW SERVICE",
    "description": "NEW SERVICE",
    "version": 1,
    "versions": 1
}


// New Service Version
> curl -X POST -H "Content-Type: application/json" --data '{"serviceId": 1, "name": "EXISTING SERVICE", "description": "NEW SERVICE VERSION", "version": 1.4}' http://localhost:8080/services
> {
    "id": "1c0fb253-99b0-44a6-8a68-43064b47d017",
    "serviceId": 1,
    "name": "EXISTING SERVICE",
    "description": "NEW SERVICE VERSION",
    "version": 1.6,
    "versions": 3
}

// New Service Version where the provided service 'version' is lower than the most up to date version for the given 'serviceId'
> curl -X POST -H "Content-Type: application/json" --data '{"serviceId": 1, "name": "EXISTING SERVICE", "description": "NEW SERVICE VERSION", "version": 0.9}' http://localhost:8080/services
> {
    "error": "new records for existing Services must have a version that is higher than the most up to date version of the existing service. Provided serviceId:'1', version:'0.900000'"
}
```

### ServiceRecords
1. `GET /service-records/{serviceId}`
  - Retrieve a list of [ServiceRecord](https://github.com/kevinliang43/kong-service/blob/main/models/service.go#L3-L9) for a given `serviceId`. If the ServiceId doesn't exist, returns an empty list.
  - Sample Request/Response:
```
// Existing Service
> curl -X GET -H "Content-Type: application/json" http://localhost:8080/service-records/1
> [
    {
        "id": "53ed660e-24be-4662-9564-f7464991d651",
        "serviceId": 1,
        "name": "Locate Us",
        "description": "Service for retrieving location info",
        "version": 1
    },
    {
        "id": "6acb5e21-dbc5-4ed7-83e3-ede75a3255d1",
        "serviceId": 1,
        "name": "Locate Us",
        "description": "Service for retrieving location info v1.1",
        "version": 1.1
    },
    {
        "id": "1c0fb253-99b0-44a6-8a68-43064b47d017",
        "serviceId": 1,
        "name": "EXISTING SERVICE",
        "description": "NEW SERVICE VERSION",
        "version": 1.6
    }
]

// Non-Existing Service
> curl -X GET -H "Content-Type: application/json" http://localhost:8080/service-records/6
> []
```

2. `GET /service-records/{serviceId}/versions/{version}`
  - Retrieve a single [ServiceRecord](https://github.com/kevinliang43/kong-service/blob/main/models/service.go#L3-L9) for a given `serviceId` and `version`. If the `serviceId` or the `version` does not match any records, a 404 response will be returned.
  - Sample Request/Response:

```
// Existing Service Id and Version
> curl -X GET -H "Content-Type: application/json" http://localhost:8080/service-records/1/versions/1.6
> {
    "id": "1c0fb253-99b0-44a6-8a68-43064b47d017",
    "serviceId": 1,
    "name": "EXISTING SERVICE",
    "description": "NEW SERVICE VERSION",
    "version": 1.6
}

// Non-Existing ServiceId/Version
> curl -X GET -H "Content-Type: application/json" http://localhost:8080/service-records/6/versions/1.7
> {
    "error": "no records found for serviceId '6' and version '1.700000'"
}
```

### Service Versions
1. `GET /service-versions/{serviceId}`
  - Retrieve a list of floats representing the `version` numbers for a given `serviceId`. If the given `serviceId` is non-existing, then an empty list will be returned.
  - Sample Request/Response:
```
// Existing serviceId
> curl -X GET -H "Content-Type: application/json" http://localhost:8080/service-versions/1
> [
    1,
    1.1,
    1.6
]

// Non-Existing serviceId
> curl -X GET -H "Content-Type: application/json" http://localhost:8080/service-versions/222
> []
```

## Design Considerations and Assumptions
1. Search Requests for a list of `Service` will be more frequent than retreiving a single `ServiceRecord`. To make retrievals for the most up to date versions of `Service`, I've created two database tables:
	- `services` will house all `ServiceRecord` records, which represent each version for a given service
	- `services_latest` will house information about the most up to date version for a given service, as well as the number of versions that it has.

This will allow for us to continue storing information about each version for a given `Service`, but also improve efficiency for search retrievals (which we presumed to be a more frequent request).

2. Returned `nextOffset` param for a given `ServicesSearchResponse` can be directly plugged into the next `ServicesSearchRequest` `offset` param. This is to allow for ease of use for any developer consuming this API (i.e a Front end team developing a GUI that makes Pagination calls).
3. `Service` logic is divided into 3 parts:
	- [service_service.go](https://github.com/kevinliang43/kong-service/blob/main/service_module/service_service.go) handles all JSON request/responses related logic.
	- [service_manager.go](https://github.com/kevinliang43/kong-service/blob/main/service_module/service_manager.go) handles core logic (validation, handling different request scenarios such as creating new version / new service, etc..)
	- [service_dao.go](https://github.com/kevinliang43/kong-service/blob/main/service_module/service_dao.go) / [service_latest_dao.go](https://github.com/kevinliang43/kong-service/blob/main/service_module/service_latest_dao.go) handles all logic related to hitting the connected DB.

Abstracting away logic into these layers allows for:

	- Dependency Injection - Allows for ease of testing with Stub/Mock objects, if unit tests were to be implemented in the future.
	- Extensibilitiy - We can add new dependencies when needed with ease, when new features are needed.


## Areas of Improvement
1. Having two database tables (that essentially store similar information), results in the following
	- Two sources of truth: Will need to keep in mind to keep these two tables consistent.
	- Two database calls upon each Create / Update / Delete request (performance inefficiency)
	- Potential solution: 
		- Make both DB calls in a single transactional call. Will prevent cases where the first update succeeds, but the second update fails that result in discrepencies between the tables
		- abstract away this logic into a Message Queue worker that have retries enabled.

2. Currently, the `filterType` is used as both the column to filter on as well as to sort on (with the provided `sortType`). We can probably add a new field for sorting as well (i.e Imagine a case where we wanted to filter on `name` of the `Service`, but sort on the `version` field).
