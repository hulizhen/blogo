# Overview

BloGo is a lightweight static blog engine built with Go.

The engine automatically reads and syncs the blog data from a seperated blog repo storing all your articles, like the [demo](https://github.com/hulizhen/blog/tree/test/demo-articles) repo.

# Quick Started

## Install Prerequisites

* Install [Docker Engine](https://docs.docker.com/engine/install/).
* Install [Docker Compose](https://docs.docker.com/compose/install/).

## Fork and Clone Repo

* Fork this [repo](https://github.com/hulizhen/blogo) to make some necessary changes later.
* Clone the repo and checkout to the recommended `main` branch.

## Configure the Application

* Create the `~/.blogo` directory, which will contain all the configuration files for our website.
* Create the `~/.blogo/config.toml` file. You may need to reference to the default `config/config.toml` file in the project.
* Create the `~/.blogo/docker.env` file to provide the needed environment variables. Currently needed variables are for the [mysql docker image](https://hub.docker.com/_/mysql).
* Provide the `~/.blogo/favicon.ico` and `~/.blogo/logo.svg` files.
* Add a webhook `https://${your.blog.domain}/webhook/github` in the settings page of your blog repo.
* Replace the domains, which you want to enable HTTPS, in the `script/init-letsencrypt.sh` file in the project.

## Launch the Application

Run the following commands at the root directory of this project.
```shell
# Fetche and ensure the renewal of a Let’s Encrypt certificate.
$ ./script/init-letsencrypt.sh

# Create and build the docker services to launch the application.
$ docker-compose up --build
```

The `blogo` application then will be powered up by launching some necessary docker containers.
These containers are from the docker services defined in the `docker-compose.yml` file:
* blogod: the main web application as the blog engine.
* mysql: the relational database to store blog articles, comments, website statistics data, etc.
* nginx: the reverse proxy server sitting between the browsers and the blog engine.
* certbot: the tool to manage [Let's Encrypt](https://letsencrypt.org/) certificates, which enables the HTTPS for our blog website.

# Acknowledgements

BloGo takes advantage of many brilliant open source libraries which save me a lot of works.
Check out the `go.mod` file to get a full list of them.
