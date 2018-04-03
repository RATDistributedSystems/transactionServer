# Transaction Server

This core of the system. The transaction server is in charge of recieving requests and processing it. That includes all the Buy, Sell, Quotes and Triggers. 

## How to build images

```
cd setup_transaction_image
./setup_image.sh
```

Your image should now be built as `transaction_server`. Alternatively you can pull from the docker cloud account.

```
docker pull asinha94/seng468_transation_server
```

## Executing

```
docker run asinha94/seng468_transaction_server
# or docker run transaction_server 
```

## Building the Corresponding Cassandra Instance

[Go here](/setup_cassandra_image)