package sessions

import (
	"time"

	"github.com/go-redis/redis"
)

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	//initialize and return a new RedisStore struct
	return &RedisStore{
		Client: client,
	}
}

//Store implementation

//Save saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
/*func (s *SessionState) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}*/
func (rs *RedisStore) Save(token string, exp time.Duration) error {
	err := rs.Client.Set(token, 0, exp).Err()
	if err != nil {
		return err
	}
	return nil
}
