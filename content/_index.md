---
---

KissLists is a very simple tool to share lists of items. 

It is a ultra-minimalistic alternative of Google Keep todo lists.

## Features

- presented as a single html page, best viewed on mobile devices
- share url to friends, pin to phone launcher, done !
- list custom color, just like Keep lists
- websockets are used to quickly sync your items across all clients
- your items are stored in a simple sqlite database, so you can use them with other tools

## Screenshots

![Screenshot](/kisslists/img/QfRBCgk.png)
![Screenshot](/kisslists/img/6I5qR5J.png)

## Install with Docker

The Docker image have tags for linux, amd64 / arm / rpi

With Docker Compose

{{< highlight go>}}
kisslists:
  image: allyouneedisgnu/kisslists
  volumes:
    - ./your/kisslists:/kisslists
  ports:
    - 80:80
{{< / highlight >}}

With Docker command line

{{< highlight go>}}
docker run -p 80:80 -v ./your/kisslists:/kisslists allyouneedisgnu/kisslists
{{< / highlight >}}

## Limitations

KissLists is so simple, some fetaures are deliberately left aside, for example

- no built-in authentication, it is your responsability to secure the access to your lists with a frontend proxy
- no user concept, it is up to you to keep url of each list

## Export lists

KissLists stores items in a an SQLite database.
This allow allow <a href="https://github.com/planetopendata/awesome-sqlite#sqlite-admin-tools">sqlite tools</a> to work.

To export them to a CSV file named kisslists.csv
{{< highlight go>}}
sqlite3 /path/to/kisslists.sqlite <<!
.headers on
.mode csv
.output kisslists.csv
select * from ListItems;
!
{{< / highlight >}}

