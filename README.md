# AWS SQS Message Mover

CLI script to move messages from one AWS SQS queue into another AWS SQS Queue.

# Installation

    go install github.com/mattrx/sqs-message-mover

# Usage

    Usage:
      sqs-message-mover [flags]

    Flags:
      -h, --help                       help for sqs-message-mover
          --loop-count int32           number of loops for receive (default 1000)
      -p, --profile string             aws profile
      -s, --source-queue string        source queue url
      -t, --target-queue string        target queue url
          --visibility-timeout int32   visibility timeout in seconds (default 60)

You have to provide an AWS profile configured in `~/.aws/accounts` and the queue urls in the form of `https://sqs.eu-central-1.amazonaws.com/{accountID}/{queueName}`.
