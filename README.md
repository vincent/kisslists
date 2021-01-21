# SharedLists

SharedLists is a dead simple shared lists server.

- no user management
- no security policies
- no export
- no import
- .. just public lists

![Screenshot](https://i.imgur.com/aVSNgGE.png)

## Install with Docker compose

```
sharedlists:
  image: allyouneedisgnu/sharedlists
  volumes:
    - ./your/sharedlists.sqlite:/sharedlists.sqlite # an empty file will wo
  environment:
    THEME: "default" # original, ocean, nature, playful
  ports:
    - 80:80
```
