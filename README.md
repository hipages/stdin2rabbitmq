stdin2rabbitmq
====

[![Build Status](https://travis-ci.com/hipagesgroup/stdin2rabbitmq.svg?token=UjDkn6mLeAgqMrrNZzpp&branch=master)](https://travis-ci.com/hipagesgroup/stdin2rabbitmq)
[![codecov](https://codecov.io/gh/hipagesgroup/stdin2rabbitmq/branch/master/graph/badge.svg?token=rA7ydQy0Oy)](https://codecov.io/gh/hipagesgroup/stdin2rabbitmq)

# explaination

The idea behind this is to be a super small and super fast script to use apache CustomLog "|stdin2rabbitmq" to send apache access logs to rabbitmq to be later consumed.

# usage

## params:

* debug: show debug info (false)
* host: rabbitmq host (localhost)
* port: rabbitmq port (ampq: 5672)
* queue: rabbitmq queue name (logs)
* rabbituser: rabbitmq username (guest)
* rabbitpass: rabbitmq password (guest)

```
echo "asdf" | stdin2rabbitmq -debug -host localhost -port 5555 -queue myqueue --rabbituser myuser --rabitpass mypass
```

```
The message is: asdf
amqp://guest:guest@localhost:5672/
2018/01/19 15:06:26  [x] Sent asdf
```

the result
```
./rabbitmqadmin get queue=myqueue
+-------------+----------+---------------+---------+---------------+------------------+-------------+
| routing_key | exchange | message_count | payload | payload_bytes | payload_encoding | redelivered |
+-------------+----------+---------------+---------+---------------+------------------+-------------+
| myqueue     |          | 22            | asdf    | 4             | string           | True        |
+-------------+----------+---------------+---------+---------------+------------------+-------------+
```

job done.

The main use case I have for this tool is to use it in combination with Apache's CustomLog Piped logging feature. In your apache conf:
```
CustomLog "|$/usr/bin/stdin2rabbitmq -host localhost -port 5555 -queue myqueue --rabbituser myuser --rabitpass mypass" "combined"
```
