## id1: all-in-one backend for everything
API, identity and access, storage, messaging, scheduling, integration

Single executable, 0.5 Ghz CPU, 15 Mb RAM, sub 1 ms latency

#### How to build executable (requires go 1.23.3+)
    $ go build
    $ ./id1

#### How to build docker image
    $ docker build -t id1:latest .
    $ docker run --rm -d -p 8080:8080 --mount type=volume,src=id1db,dst=/mnt/id1db --name id1 id1:latest

#### optional .env file
    PORT=8080
    DBPATH=id1db

## Web Socket commands

#### Create or update
    set:/path/to/key
    <data>

#### Append to existing value
    add:/path/to/key
    <data>

#### Rename key
    mov:/path/to/source
    <path/to/target>

#### Read value
    get:/path/to/key

#### Delete value
    del:/path/to/key

#### List
    get:/path/to/*

#### List up to 100 items under 1Kb recursively
    get:/path/to/*?recursive=true&limit=100&size-limit=1024

#### List keys only
    get:/path/to/*?keys=true

#### List children
    get:/path/to/*?children=true

#### Create and delete after 60 seconds
    set:/path/to/key?ttl=60
    <data>

#### Schedule create, archive and delete
    set:/path/.after.1745000011111
    set:/path/to/file
    <data>
  
    set:/path/.after.1745000022222
    mov:/path/to/file
    path/archive/file?ttl=86400

## HTTP equivalent
- POST = set
- GET = get, list
- DELETE = del
- PATCH = add, move
    
Examples:

    GET https://id1.au/max/pub/*
    GET https://id1.au/max/pub/name
    POST https://id1.au/max/pub/name
    Max
    PATCH https://id1.au/max/pub/name
    X-Move-To: max/arch/pub/name
    DELETE https://id1.au/max/pub/name

