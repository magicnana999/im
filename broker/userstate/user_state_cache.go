package userstate

const (
	NEW_CONNECT_SCRIPT string = `
			local key1 = KEYS[1]
			local key2 = KEYS[2]
			local filed = ARGV[1]
			local val = ARGV[2]
			local ov = redis.call("get", key1)
			redis.call("set", key1, val)
			redis.call("hset", key2, filed, val)
			return ov
		`
)
