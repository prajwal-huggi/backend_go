package main

// =============================================================================
// REDIS CACHING — the Cache-Aside pattern in Go
// =============================================================================
// This shows exactly how caching slots into your existing
// handler -> service -> repository -> database flow.
//
// Library: github.com/redis/go-redis/v9
// Install:  go get github.com/redis/go-redis/v9
// =============================================================================

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// ---------------------------------------------------------------------------
// The data model — same User struct your CRUD app already has.
// ---------------------------------------------------------------------------
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ---------------------------------------------------------------------------
// UserRepository — this is your existing repository layer, now with a
// Redis client added alongside the database. In your real app, `db` would
// be your actual database connection (Postgres, MySQL, etc.). Here it's
// faked with a function so the example runs conceptually.
// ---------------------------------------------------------------------------
type UserRepository struct {
	redis *redis.Client
	// db  *sql.DB   <- in your real app, your DB handle goes here
}

// cacheKey builds a consistent Redis key for a user.
// Pattern "user:<id>" keeps keys organized and predictable — the same
// naming discipline you saw with Kafka topics.
func cacheKey(id string) string {
	return fmt.Sprintf("user:%s", id)
}

// ===========================================================================
// READ — the Cache-Aside pattern
// ===========================================================================
func (r *UserRepository) GetUser(ctx context.Context, id string) (*User, error) {
	key := cacheKey(id)

	// ----- STEP 1: check Redis first -----
	cached, err := r.redis.Get(ctx, key).Result()

	if err == nil {
		// CACHE HIT — found it in Redis. Deserialize and return immediately.
		log.Printf("✅ CACHE HIT for %s", key)
		var user User
		if jsonErr := json.Unmarshal([]byte(cached), &user); jsonErr == nil {
			return &user, nil
		}
		// if unmarshal fails, fall through and treat it like a miss
	}

	if err != nil && !errors.Is(err, redis.Nil) {
		// A real Redis error (network down, etc.) — NOT just "key missing".
		// Important design choice: we log it but DON'T fail the request.
		// The database is still there, so we fall through to it.
		// Caching should never take your app down if Redis has a hiccup.
		log.Printf("⚠️ Redis error (falling back to DB): %v", err)
	}

	// ----- STEP 2: CACHE MISS — go to the database -----
	log.Printf("❌ CACHE MISS for %s — querying database", key)
	user, dbErr := r.getUserFromDB(id)
	if dbErr != nil {
		return nil, dbErr
	}

	// ----- STEP 3: store a copy in Redis for next time, with a TTL -----
	payload, _ := json.Marshal(user)
	// EX 300 seconds = this entry auto-expires in 5 minutes.
	if setErr := r.redis.Set(ctx, key, payload, 5*time.Minute).Err(); setErr != nil {
		// Again — if caching the result fails, we don't fail the request.
		log.Printf("⚠️ failed to cache %s: %v", key, setErr)
	}

	return user, nil
}

// ===========================================================================
// UPDATE — write to DB, then INVALIDATE the cache
// ===========================================================================
func (r *UserRepository) UpdateUser(ctx context.Context, user *User) error {
	// ----- STEP 1: update the database (the source of truth) FIRST -----
	if err := r.updateUserInDB(user); err != nil {
		return err
	}

	// ----- STEP 2: invalidate the cache so stale data isn't served -----
	// We DELETE the key rather than updating it. The next read will be a
	// miss, repopulating the cache with guaranteed-fresh data from the DB.
	// (Deleting is simpler and less error-prone than trying to keep the
	// cached copy perfectly in sync on every write.)
	key := cacheKey(user.ID)
	if err := r.redis.Del(ctx, key).Err(); err != nil {
		log.Printf("⚠️ failed to invalidate cache for %s: %v", key, err)
		// Not fatal — the TTL would eventually expire it anyway.
	} else {
		log.Printf("🗑️  invalidated cache for %s after update", key)
	}

	return nil
}

// ===========================================================================
// Fake DB functions — stand-ins for your real database calls.
// ===========================================================================
func (r *UserRepository) getUserFromDB(id string) (*User, error) {
	// Simulate a slow DB read
	time.Sleep(100 * time.Millisecond)
	return &User{ID: id, Name: "Prajwal", Email: "prajwal@example.com"}, nil
}

func (r *UserRepository) updateUserInDB(user *User) error {
	time.Sleep(100 * time.Millisecond)
	log.Printf("💾 database updated for user %s", user.ID)
	return nil
}

// ===========================================================================
// main — demonstrates the hit/miss behavior
// ===========================================================================
func main() {
	ctx := context.Background()

	// Connect to Redis (default local Redis runs on port 6379)
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0, // Redis has 16 numbered DBs (0-15); 0 is the default
	})
	defer rdb.Close()

	// Quick connectivity check
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("cannot reach Redis: %v", err)
	}
	log.Println("connected to Redis")

	repo := &UserRepository{redis: rdb}

	// First read — CACHE MISS (goes to DB, ~100ms, then caches)
	log.Println("\n--- First read ---")
	repo.GetUser(ctx, "123")

	// Second read — CACHE HIT (from Redis, <1ms)
	log.Println("\n--- Second read (should be a hit) ---")
	repo.GetUser(ctx, "123")

	// Update — invalidates the cache
	log.Println("\n--- Update user ---")
	repo.UpdateUser(ctx, &User{ID: "123", Name: "Prajwal Kumar", Email: "new@example.com"})

	// Third read — CACHE MISS again (because we invalidated), then re-cached
	log.Println("\n--- Read after update (miss again, then fresh) ---")
	repo.GetUser(ctx, "123")
}
