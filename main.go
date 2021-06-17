package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/spf13/cobra"
)

var (
	profile           string
	sourceQueue       string
	targetQueue       string
	loopCount         int32
	visibilityTimeout int32
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "aws profile")
	rootCmd.PersistentFlags().StringVarP(&sourceQueue, "source-queue", "s", "", "source queue url")
	rootCmd.PersistentFlags().StringVarP(&targetQueue, "target-queue", "t", "", "target queue url")
	rootCmd.PersistentFlags().Int32VarP(&loopCount, "loop-count", "", 1000, "number of loops for receive")
	rootCmd.PersistentFlags().Int32VarP(&visibilityTimeout, "visibility-timeout", "", 60, "visibility timeout in seconds")
}

var rootCmd = &cobra.Command{
	Use: "sqs-message-mover",
	RunE: func(cmd *cobra.Command, args []string) error {
		sess, err := session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Profile:           profile,
		})

		if err != nil {
			return err
		}

		svc := sqs.New(sess)

		for i := 1; i <= int(loopCount); i++ {
			receiveResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(sourceQueue),
				MaxNumberOfMessages: aws.Int64(10),
				MessageAttributeNames: []*string{
					aws.String(sqs.QueueAttributeNameAll),
				},
				WaitTimeSeconds:   aws.Int64(10),
				VisibilityTimeout: aws.Int64(int64(visibilityTimeout)),
			})

			if err != nil {
				return err
			}

			for _, message := range receiveResult.Messages {
				if _, err := svc.SendMessage(&sqs.SendMessageInput{
					QueueUrl:          aws.String(targetQueue),
					MessageBody:       message.Body,
					MessageAttributes: message.MessageAttributes,
				}); err != nil {
					return err
				}

				if _, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      aws.String(sourceQueue),
					ReceiptHandle: message.ReceiptHandle,
				}); err != nil {
					return err
				}
			}
		}

		return nil
	},
}
