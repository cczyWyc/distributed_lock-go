-- 1. Check if it is the expected value(check if it is your lock)
-- 2. if it is and delete key, if not return a value
if redis.Call("get", KEYS[1]) == ARGV[1] then
    return redis.Call("del", KEYS[1])
else
    return 0
end