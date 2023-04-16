# Whale REST API

Whale is a RESTful API written in Golang designed to manage user defined catalog data by 62teknologi.com.

## Table of Contents

1. [Introduction](#introduction)
2. [Catalog Description](#catalog-description)
3. [Features](#features)
4. [Installation](#installation)
5. [API Endpoints](#api-endpoints)
6. [Usage Examples](#usage-examples)
7. [Contributing](#contributing)
8. [License](#license)

## Introduction

Whale REST API is built with a focus on simplicity, reliability, and extensibility. With this API, users can manage catalog data for different whale species, including their names, scientific classifications, descriptions, images, and more.

## Catalog Description

A catalog is a collection of information or data that describes a set of products, services, or other items.

In the context of a whale REST API, a catalog would refer to a collection of `User Defined Data` that has certain `Behaviors` and `Associations`

###  Whale catalog behaviors
- can be created
- can be retrieved
- can be updated
- can be deleted

### Whale catalog associations
- may has one category
- may has many items
- may has many to many groups
- may has many comments
- may belong to certain user 

## Features

- Easy-to-use RESTful API
- Easy to setup
- Easy to Customizable
- Written in Golang for high performance and concurrency 
- Robust data validation and error handling
- Well-documented API endpoints

## Installation

To install and run Whale REST API on your local machine, follow these steps:

1. Clone the repository:

    git clone https://github.com/whale-rest-api.git

1. Change directory to the cloned repository:

    cd whale-rest-api

1. Build the application:

    go build

1. Run the server:

    ./whale-rest-api

The API server will start running at `http://localhost:10081`. You can now interact with the API using your preferred API client or through the command line with `curl`.


## Set Up a config (WIP)
- copy .env.example

## Generate a Catalog (WIP)
- change directory to console
- generate the catalogue
  go catalogue.go [name]

## API Endpoints

| Method | Endpoint | Description |
| - | -| - |
| GET | /api/v1/catalog/:name | Retrieve a list of all whales in the catalog |
| GET | /api/v1/catalog/:name/:id | Retrieve a specific whale by ID |
| POST | /api/v1/catalog/:name | Add a new whale to the catalog |
| PUT | /api/v1/catalog/:name/:id | Update information for a specific whale by ID |
| DELETE | /api/v1/catalog/:name/:id | Delete a specific whale from the catalog by ID |

For more detailed information about each endpoint, including request and response format, please refer to the [API documentation](./API_DOCUMENTATION.md).

### Usage Examples  (WIP)

Here are some examples of how to interact with the Whale REST API using `curl`:

1. Get a list of all whales:
2. Get a specific whale by ID:


## Contributing

If you'd like to contribute to the development of the Whale REST API, please follow these steps:

1. Fork the repository
2. Create a new branch for your feature or bugfix
3. Commit your changes to the branch
4. Create a pull request, describing the changes you've made

We appreciate your contributions and will review your pull request as soon as possible.

## License

This project is licensed under the MIT License. For more information, please see the [LICENSE](./LICENSE) file.