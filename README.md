# Stealth Commerce

Project test simple ecommerce using golang

This project has 2 user access :

- Admin
- Buyer

#### Admin

- can access login endpoint. use username: `admin` and password: `admin123`
- can access list product endpoints
- can access detail product endpoints
- can access add product endpoints
- can access update product endpoints
- can access delete product endpoints

#### Buyer

- buyer can access endpoint using `secret_key` provided into seed data in `db.sql` file
- can view products list
- can view product detail
- can add product to cart
- can update cart
- can checkout (order and pay)

#### Run the Applications following these steps:

1. clone the project
2. copy file .env-example to .env

```bash
$ cp .env-example .env
```

3. run the sql in `db.sql` file to your database
4. install modules

```bash
$ go mod tidy

#or
$ go mod download
```

5. run the application into development mode

```bash
$ make dev
```

#### Build process:

```bash
# run build command
$ make build

# execute
$ ./target/stealth-commerce
```

TODO:

- auto cancel order if not paid after 30 minute
- migration
