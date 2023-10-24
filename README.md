# KissLists

<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-79%25-brightgreen.svg?longCache=true&style=flat)</a>

KissLists is a very simple shared lists server.
<br><br>

<img align="right" width="40%" src="https://i.imgur.com/EXNsN7s.png">

with

- mobile optimised design
- basic theme support
- websockets messages
- sqlite database

but

- no built-in authentication 
- no user management
- no admin panel
- no import / export

<br>

## Install

### with Docker compose

> available tags for linux amd64 and arm <br>

```
services:
  kisslists:
    image: allyouneedisgnu/kisslists
    volumes:
      - ./your/kisslists:/kisslists
    ports:
      - 80:80
```
