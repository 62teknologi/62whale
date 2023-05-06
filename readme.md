# 62Whale

62Whale is a REST API written in Golang designed to manage user defined catalog data and can run independly as a stateless service.

The main goal of 62Whale is to reduce repetition of creating and managing catalog data.

Created by 62teknologi.com, perfected by Community.

## Catalog
This introduction will help You explain the concept and characteristic of catalog.

### Concept

Catalog is a collection of data, for example; products, articles, galleries.

In the context of a 62Whale, a catalog would refer to a collection of `User Defined Data` that has certain `information`, `Behaviors`, `Associations` and `characteristic`.

You will learn how to defined data on later section.

### Information
- Must Have ID
- Must Have Slug
- Must Have Created At
- Must Have Updated At
- Must Have Deleted At

### Behaviors
- can be created
- can be retrieved
- can be updated
- can be deleted

### Associations
- may has one category
- may has many items
- may has many to many groups
- may has many comments
- may belong to certain user 


## Running 62Whale

Follow the instruction below to running 62Whale on Your local machine.

### Prerequisites
Make sure have preinstalled this prerequisites app before you continue to installation manual. we don't include how to install these app below Most of this prerequisites is a free app which you can find the "How to" installation tutorial anywhere in web and different machine os have different way to install.
- MySql

### Installation manual
This installation manual will guide You to running the binary on Your ubuntu or mac terminal.

1. Clone the repository
```
git clone https://github.com/62teknologi/62Whale
```

1. Change directory to the cloned repository
```
cd 62Whale
```

1. Create .env base on .env.example
```
cp .env.example .env
```

1. change DB variable on .env using Your mysql configuration or the staging database on cloud server eg
```
HTTP_SERVER_ADDRESS=0.0.0.0:10081
DB_DRIVER=mysql
DB_SOURCE_1=root@tcp(127.0.0.1:3306)/whale_local
```

1. Run the server
```
./62whale
```

The API server will start running at `http://localhost:10081`. You can now interact with the API using your preferred API client or through the command line with `curl`.

## API Endpoints

### Retrieve Catalog by id

#### Endpoint
```
GET /api/v1/catalog/:name/:id
```

### Retrieve Catalog List

#### Endpoint
```
GET /api/v1/catalog/:name
```

#### Parameter
| Name | Def | Description |
| - | - | - |
| page | 1 | return response in pagination format, eg: ```page=1``` will return first page of the response    |
| per_page | 30 | set how many data per pagination response |
| search | null | filter response by string |
| :field | null | filter specific column you want, eg: if your catalog have ```user_id``` field then you can add ```user_id=1``` to params for searching all catalog where user_id is 1. it support multi value by sending ```user_id[]``` instead ```user_id``` |

### Create Catalog

#### Endpoint
```
POST /api/v1/catalog/:name
```

#### Parameter
| Name | Def | Description |
| - | - | - |
| :field | null | field you want to insert |
| groups[] | null | lorem |
| items[] | null | lorem |


### Update Catalog

#### Endpoint
```
PUT /api/v1/catalog/:name/:id
```

#### Parameter
| Name | Def | Description |
| - | - | - |
| :field | null | field you want to insert |
| groups[] | null | lorem |

### Delete Catalog

#### Endpoint
```
DEL /api/v1/catalog/:name/:id
```
# Set Up a Catalog
- WIP   

## Generate Catalog
- WIP

## Set Information
- WIP

## Set Validation
- WIP

## Set Associations
- WIP

# Contributing

If you'd like to contribute to the development of the Whale REST API, please follow these steps:

1. Fork the repository
2. Create a new branch for your feature or bugfix
3. Commit your changes to the branch
4. Create a pull request, describing the changes you've made

We appreciate your contributions and will review your pull request as soon as possible.

## Must Preserve Characteristic 
- Reduce repetition
- Easy to use REST API
- Easy to setup
- Easy to Customizable
- high performance
- Robust data validation and error handling
- Well documented API endpoints

## License

This project is licensed under the MIT License. For more information, please see the [LICENSE](./LICENSE) file.

# About 62
**E.nam\Du.a**

Indonesian language; spelling: A-num\Due-wa

Origin: Enam Dua means ‘six-two’ or sixty two. It is Indonesia’s international country code (+62), that was also used as a meme word for “Indonesia” by “Indonesian internet citizen” (netizen) in social media.