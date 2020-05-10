# What's it?
When you use redis in your project, How to mock it is the best problem to UT. `redisgomock` is one library for mock redis.

# Why write by myself?
[redigo](https://github.com/gomodule/redigo) recommend [rafaeljusto/redigomock](https://github.com/rafaeljusto/redigomock). But I think it can be better:
- Use Interface, not struct
- Fuzzy Match can be simple
  
In fact, I want submit pull request to rafaeljusto/redigomock, but will change so much, So I write one by myself

# How to Use
```
conn, mock := NewMockClient()
// prepare 
mock.ExpectDo("SADD").WillReply(int64(1))
mock.ExpectClose()

// do
count, err := redis.Int64(conn.Do("SADD", "element"))
conn.Close()

// assert
if count != 1 {
    ....
}
```

# TODO
- [ ] UT Coverage > 70%
- [ ] design doc