package message

import (
    "log"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/sns"
    "github.com/aws/aws-sdk-go/aws/session"
    "encoding/json"
)

const (
    DO_NOTHING = false // メッセージ送信しない場合はtrueにしてください
    AWS_REGION = "ap-northeast-1"
    TOPIC_ARN  = "arn:aws:sns:ap-northeast-1:123456789012:SampleTopic"
)

type Message struct {
    OperationType string `json:"operationType"`
    ID string `json:"id"`
    ItemCode string `json:"itemCode"`
    Text   string   `json:"text"`
}

// SQSにメッセージを送信します
func SendMessage(operationType string, idStr string, itemCode string, text string, filterValue string){
    if DO_NOTHING{
        log.Print(filterValue)
        log.Print("DO_NOTHING = true のためメッセージは送信しません")
        msgBody := createMessage(operationType, idStr, itemCode, text)
        log.Print(msgBody)
        return
    }
    sess := session.Must(session.NewSession())
    svc := sns.New(sess, aws.NewConfig().WithRegion(AWS_REGION))
    // 送信内容を作成
    msgBody := createMessage(operationType, idStr, itemCode, text)
    
    params := &sns.PublishInput{
        Message:  aws.String(msgBody),
        MessageStructure: aws.String("json"),
        TopicArn: aws.String(TOPIC_ARN),
    }
    if filterValue != "" {
        params.SetMessageAttributes(map[string]*sns.MessageAttributeValue{
            "Key": {
                DataType:    aws.String("String"),
                StringValue: aws.String(filterValue),
            },
        })
    }

    sqsRes, err := svc.Publish(params)
    if err != nil {
        log.Fatal(err)
    }
    log.Print("SQSMessageID", *sqsRes.MessageId)
    log.Print("SQSMessageBody", msgBody)
}

// json形式でメッセージ本体を作成します
func createMessage(operationType string, idStr string, itemCode string, text string) string {
    sqsmsg := new(Message)
    sqsmsg.OperationType = operationType
    sqsmsg.ID = idStr
    sqsmsg.ItemCode = itemCode
    sqsmsg.Text = text
    sqsmsgJson, _ := json.Marshal(sqsmsg)

    message := map[string]string{
        "default": "This is default message.",
        "sqs": string(sqsmsgJson),
    }
    messageJson, _ := json.Marshal(message)
    return string(messageJson)
}
