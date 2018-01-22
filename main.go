//

package main

import (
  "github.com/streadway/amqp"
  "bufio"
  "flag"
  "fmt"
  "log"
  "os"
  "strconv"
)

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}

func read_in_stdin(debug bool) *bufio.Scanner {
  // massive copy and paste
  // https://flaviocopes.com/go-shell-pipes/

  info, err := os.Stdin.Stat()
  failOnError(err, "Failed to connect to stdin")

  if debug {
    fmt.Println("stdin buffer size = " + strconv.FormatInt(info.Size(), 10))
  }

  // removed the size test as this can return nothing but 0 on some OS's.
  // if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
  if info.Mode()&os.ModeCharDevice == os.ModeCharDevice {
    fmt.Println("The command is intended to work with pipes.")
    fmt.Println("Usage: echo mytext | stdin2rabbitmq")
    os.Exit(1)
  }

    // https://medium.com/golangspec/in-depth-introduction-to-bufio-scanner-in-golang-55483bb689b4
    scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
    //for scanner.Scan() {
    //    fmt.Println(scanner.Text())
    //}

    return scanner
}

func post_to_rabbitmq(debug bool, host string, payload string, port string, queue string, rabbituser string, rabbitpass string) {
  // http://www.agiratech.com/message-queues-golang-rabbitmq/

  if debug {
    fmt.Println("amqp://" + rabbituser + ":" + rabbitpass + "@" + host + ":" + port + "/")
  }
  conn, err := amqp.Dial("amqp://" + rabbituser + ":" + rabbitpass + "@" + host + ":" + port + "/")
  failOnError(err, "Failed to connect to RabbitMQ")
  defer conn.Close()

  ch, err := conn.Channel()
  failOnError(err, "Failed to open a channel")
  defer ch.Close()

  q, err := ch.QueueDeclare(
    queue, // name
    false,   // durable
    false,   // delete when unused
    false,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )
  failOnError(err, "Failed to declare a queue")

  body := payload
  err = ch.Publish(
    "",     // exchange
    q.Name, // routing key
    false,  // mandatory
    false,  // immediate
    amqp.Publishing{
      ContentType: "text/plain",
      Body:        []byte(body),
    })
  if debug {
    log.Printf(" [x] Sent %s", body)
  }
  failOnError(err, "Failed to publish a message")

}

func main() {

  debug      := flag.Bool("debug", false, "turn on debug messages, default false")
  host       := flag.String("host", "localhost", "the rabbitmq hostname, default: localhost")
  port       := flag.String("port", "5672", "the rabbitmq amqp port, default: 5672")
  queue      := flag.String("queue", "logs", "the name of the rabbitmq queue, default: logs")
  rabbituser := flag.String("user", "guest", "the rabbitmq username, default: guest")
  rabbitpass := flag.String("pass", "guest", "the rabbitmq password, default: guest")
  flag.Parse()

  log_line := read_in_stdin(*debug)
  for log_line.Scan() {
    post_to_rabbitmq(*debug, *host, log_line.Text(), *port, *queue, *rabbituser, *rabbitpass)
  }

}

