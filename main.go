//

package main

import (
  "github.com/streadway/amqp"
  "bufio"
  "flag"
  "fmt"
  "io"
  "log"
  "os"
  "strings"
)

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}

func read_in_stdin(debug bool) string {
  // massive copy and paste
  // https://flaviocopes.com/go-shell-pipes/

  info, err := os.Stdin.Stat()
  failOnError(err, "Failed to connect to RabbitMQ")

  if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
    fmt.Println("The command is intended to work with pipes.")
    fmt.Println("Usage: fortune | gocowsay")
    os.Exit(1)
  }

  reader := bufio.NewReader(os.Stdin)
  var output []rune

  for {
    input, _, err := reader.ReadRune()
    if err != nil && err == io.EOF {
      break
    }
    output = append(output, input)
  }

  if debug {
    fmt.Print("The message is: ")
    for j := 0; j < len(output); j++ {
      fmt.Printf("%c", output[j])
    }
  }

  stringout := strings.TrimSuffix(string(output), "\n")
  return stringout

}

func post_to_rabbitmq(debug bool, host string, payload string, port string, queue string, rabbituser string, rabbitpass string) {

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
  post_to_rabbitmq(*debug, *host, log_line, *port, *queue, *rabbituser, *rabbitpass)

}

