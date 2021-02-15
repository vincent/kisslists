# SharedLists

<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-79%25-brightgreen.svg?longCache=true&style=flat)</a>

SharedLists is a dead simple shared lists server.

- no user management
- no security policies
- no export
- no import
- .. just public lists

![Screenshot](https://i.imgur.com/hhyCr3b.png)

## Install with Docker compose

Sharedlists image size is only 12 Mb.

```
sharedlists:
  image: allyouneedisgnu/sharedlists
  volumes:
    - ./your/sharedlists.sqlite:/sharedlists.sqlite # an empty file will do
  ports:
    - 80:80
```
