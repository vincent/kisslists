# KissLists

<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-79%25-brightgreen.svg?longCache=true&style=flat)</a>

<img align="right" width="40%" src="https://i.imgur.com/hhyCr3b.png">

KissLists is a dead simple shared lists server.

- no user management
- no security policies
- no export
- no import
- .. just public lists

<br><br><br><br>

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
