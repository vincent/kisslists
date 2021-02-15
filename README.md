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

## Install with Docker compose

kisslists image size is only 20 Mb.

```
kisslists:
  image: allyouneedisgnu/kisslists
  volumes:
    - ./your/kisslists.sqlite:/kisslists.sqlite # an empty file will do
  ports:
    - 80:80
```
