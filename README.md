## sesgo

`sesgo` is a single, self-contained Go function that sends email via the AWS Simple Email Service (SES).

`sesgo` is useful when you want to send an email via SES, but you don't want to pull in a large dependency like `aws-sdk-go`.

## Usage

Try `sesgo` by running test/email.go:

    cd sesgo/test
    export AWS_ACCESS_KEY_ID="..." AWS_SECRET_ACCESS_KEY="..."
    go run email.go to@example.com 'John Doe <from@heytoaster.com>' 'Subject line' 'Body text'
