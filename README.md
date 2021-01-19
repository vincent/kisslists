# SharedLists

SharedLists is a dead simple shared lists server.

- no user management
- no security policies
- no export
- no import
- .. just public lists

## Install with Docker compose

```
sharedlists:
  image: allyouneedisgnu/sharedlists
  volumes:
    - ./your/sharedlists.sqlite:/sharedlists.sqlite
  ports:
    - 80:80
```
