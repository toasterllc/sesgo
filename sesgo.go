package sesgo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func SendEmail(awsAccessKey, awsSecretKey, to, from, subject, body string) error {
	const Region = "us-west-1"
	const Service = "ses"
	const Host = "email.us-west-1.amazonaws.com"
	const Method = http.MethodPost
	const URL = "https://" + Host + "/"
	const ContentType = "application/x-www-form-urlencoded"
	const CanonicalURI = "/"
	const CanonicalQueryString = ""
	const SignedHeaders = "content-type;host;x-amz-date"
	const Algorithm = "AWS4-HMAC-SHA256"

	sign := func(key []byte, msg string) []byte {
		mac := hmac.New(sha256.New, key)
		mac.Write([]byte(msg))
		return mac.Sum(nil)
	}

	getSigningKey := func(dateStamp string) []byte {
		return sign(sign(sign(sign([]byte("AWS4"+awsSecretKey), dateStamp), Region), Service), "aws4_request")
	}

	hashed := func(x string) string {
		b := sha256.Sum256([]byte(x))
		return hex.EncodeToString(b[:])
	}

	payload := url.Values{
		"Action":                    {"SendEmail"},
		"Destination.CcAddresses":   {},
		"Message.Body.Text.Charset": {"UTF-8"},
		"Message.Subject.Charset":   {"UTF-8"},
		"Version":                   {"2010-12-01"},

		"Destination.ToAddresses.member.1": {to},
		"Source":                           {from},
		"Message.Subject.Data":             {subject},
		"Message.Body.Text.Data":           {body},
	}

	payloadHash := hashed(payload.Encode())

	t := time.Now().UTC()
	amzDate := t.Format("20060102T150405Z")
	dateStamp := t.Format("20060102")

	canonicalHeaders := "content-type:" + ContentType + "\n" + "host:" + Host + "\n" + "x-amz-date:" + amzDate + "\n"
	canonicalRequest := Method + "\n" + CanonicalURI + "\n" + CanonicalQueryString + "\n" + canonicalHeaders + "\n" + SignedHeaders + "\n" + payloadHash
	credentialScope := dateStamp + "/" + Region + "/" + Service + "/" + "aws4_request"
	stringToSign := Algorithm + "\n" + amzDate + "\n" + credentialScope + "\n" + hashed(canonicalRequest)
	signature := hex.EncodeToString(sign(getSigningKey(dateStamp), stringToSign))
	authorizationHeader := Algorithm + " " + "Credential=" + awsAccessKey + "/" + credentialScope + ", " + "SignedHeaders=" + SignedHeaders + ", " + "Signature=" + signature

	client := &http.Client{}
	req, err := http.NewRequest(Method, URL, strings.NewReader(payload.Encode()))
	req.Header.Add("Content-Type", ContentType)
	req.Header.Add("Authorization", authorizationHeader)
	req.Header.Add("X-Amz-Date", amzDate)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error %v: %v", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
