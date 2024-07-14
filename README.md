# REST API with Posgresql template

## Purpose

REST API template ! 

## features

- ✅ CORS 
- ✅ Paypal paiement 
- ✅ Gorm
- ✅ Gin 
- ✅ Google oAuth
- ✅ JWT auth

## Usage

Setup local environment with docker-compose

```make up```

Tear down local environment with docker-compose

```make down```

Run server locally

```make run```

Run migrate

```make migrate```


### Local

We use docker-compose to setup local environment

### Env variables

Use `.env.example` to generate your `.env` file


## Structure

```
api
├───handler
├───migrate
├───model
├───pkg
│   └───middleware
├───route
└───store
```

## How to Connect to a Docker Database via pgAdmin4

### Steps to Follow

1. **Check if the containers are running**:

    Use the `docker ps` command to check if the `postgres` and `pgadmin` containers are running. Look for their names or IDs in the list of active containers.

    ```shell
    docker ps
    ```

2. **Obtain the IP address of the PostgreSQL container**:

    Use the `docker inspect` command to get details of the `postgres` container. You can specify the name or ID of the `postgres` container.

    ```shell
    docker inspect <container_ID>
    ```

    In the command output, find the `NetworkSettings` section. The IP address of the `postgres` container is located under `IPAddress` in the corresponding network interface.

3. **Open pgAdmin4**:

    Go to the URL `http://localhost:5050` (or the URL you configured for pgAdmin4 in your `docker-compose.yml` file).

4. **Log in to pgAdmin4**:

    Once you access the pgAdmin4 interface, log in using the credentials configured in your `docker-compose.yml` file (variables `PGADMIN_DEFAULT_EMAIL` and `PGADMIN_DEFAULT_PASSWORD`).

5. **Add a new server**:

    - In pgAdmin4, in the left navigation panel, right-click on "Servers" and select "Create" > "Server...".
    - Give your server a name.
    - In the "Connection" tab, fill in the required fields to connect to your PostgreSQL database:
        - **Host name/address**: Enter the IP address of the `postgres` container obtained earlier.
        - **Port**: Use port `5432` or the port you configured for your PostgreSQL container.
        - **Maintenance database**: Enter the name of the database you want to use.
        - **Username** and **Password**: Use the PostgreSQL credentials configured in your `.env` file.

6. **Test the connection**:

    Once you've filled in all the connection information, click "Save" to save the settings. Navigate within your newly created server to ensure the connection is successful.

If you follow these steps and configure your services and settings correctly, you should be connected to your PostgreSQL database hosted in a Docker container via pgAdmin4.
