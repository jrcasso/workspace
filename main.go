package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jrcasso/conduit/conduit"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() { cancel() }()
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			S3ForcePathStyle: aws.Bool(true),
			Region:           aws.String("us-east-1"),
			Endpoint:         aws.String("http://localstack:4566"),
		},
	}))
	t := conduit.NewConduit(*sess, myTransform, conduit.Config{
		S3Egress: os.Getenv("CONDUIT_S3_EGRESS_BUCKET"),
		QueueUrl: os.Getenv("CONDUIT_QUEUE_URL"),
	})

	if err := t.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func myTransform(t conduit.Transformable, uploadQueue chan<- conduit.Upload) {
	fmt.Println("Transforming record...")
	// Do some transformation on t.Data
	newData := fmt.Sprintf("new data %v", t.Data)

	uploadQueue <- conduit.Upload{
		Transformable: conduit.Transformable{
			Data:   newData,
			Record: t.Record,
		},
		Key: fmt.Sprintf("transformed-%v", t.Record.S3.Object.Key),
	}
	fmt.Println("Transformed record!")
}
