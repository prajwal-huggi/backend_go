# Redis Caching in Go — Learning Project

## Run Redis locally
```bash
docker-compose up -d
```
- Redis runs on localhost:6379
- RedisInsight UI runs on http://localhost:5540 (connect it to host `redis`, port 6379)

## Run the Go example
```bash
go mod init redis-demo
go get github.com/redis/go-redis/v9
go run main.go
```

## What to watch for in the output
- First read  → CACHE MISS (slow, ~100ms, hits the fake DB)
- Second read → CACHE HIT (fast, from Redis)
- After update → cache is invalidated (deleted)
- Next read   → MISS again, then re-cached fresh

## Useful redis-cli commands (run inside the container)
```bash
docker exec -it redis redis-cli

# then inside:
KEYS *              # list all keys (fine for learning; avoid in prod)
GET user:123        # see a cached value
TTL user:123        # seconds left before it expires
DEL user:123        # manually delete a key
FLUSHALL            # wipe everything (learning only!)
```

## Experiments to build intuition
1. Run the program twice quickly — watch hit vs miss.
2. In redis-cli, run `TTL user:123` right after a read — watch the countdown.
3. Wait 5 minutes, read again — it'll be a MISS because the TTL expired.
4. Change the TTL in code from 5*time.Minute to 10*time.Second and watch how much faster it expires.
