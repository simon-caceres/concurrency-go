# concurrency-go

## install mailhog for testing
```
go get github.com/mailhog/mhsendmail
```
## run mailhog local
```
mhsendmail -smtp-addr=localhost:1025

or


mailhog
```

## install redis
```
brew install redis
```

## run redis local
```
redis-server
```

## install postgresql
```
brew install postgresql
```

## run postgresql local
```
pg_ctl -D /usr/local/var/postgres start
```
