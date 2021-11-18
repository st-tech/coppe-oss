package coppe

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adhocore/gronx"
	"github.com/st-tech/coppe-oss/src/app/services"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func NotifyPubsub(ctx context.Context, m PubSubMessage) error {
	fileNameToDataMap, err := services.ReadFiles("./yaml")
	if err != nil {
		return fmt.Errorf("file retrieval failed. %v", err)
	}

	var validations []services.Validation

	for fileName, data := range fileNameToDataMap {
		validationsInFile, err := services.UnmarshalToValidations(data)
		if err != nil {
			return fmt.Errorf("yaml unmarshal failed. %v", err)
		}
		for i, v := range validationsInFile {
			v.Identifier.FoundAt.FileName = fileName
			v.Identifier.FoundAt.Index = i
			validations = append(validations, v)
		}
	}

	publishData, err := generatePublishDataFromValidations(validations)
	if err != nil {
		return err
	}

	projectId := os.Getenv("GCP_PROJECT_ID")
	topicId := os.Getenv("TOPIC_NAME_VALIDATOR")

	err = services.PublishTopicMessages(ctx, projectId, topicId, publishData); if err != nil {
		msg := err.Error()
		services.SendErrorLog(msg)
		services.PostSlackMessage(msg, "")
	}

	return err
}

func CheckRule(ctx context.Context, m PubSubMessage) error {
	v, err := services.UnmarshalToValidation(m.Data)
	if err != nil {
		msg := fmt.Sprintf("unmarshal pubsub message failed: %v", err)
		services.SendErrorLog(msg)
		services.PostSlackMessage(msg, "")
		return err
	}

	err = services.Validate(ctx, v); if err != nil {
		msg := err.Error()
		services.SendErrorLog(msg)
		services.PostSlackMessage(msg, "")
	}
	return err
}


func generatePublishDataFromValidations(validations []services.Validation) ([][]byte, error) {
	var publishData [][]byte

	gron := gronx.New()
	timezone := os.Getenv("TIMEZONE")
	local, _ := time.LoadLocation(timezone)
	localTimeNow := time.Now().In(local)

	for _, v := range validations {
		scheduled, err := gron.IsDue(v.Schedule, localTimeNow)
		if err != nil {
			log.Printf("wrong schedule syntax: %s", v.Schedule)
			continue
		}

		if scheduled {
			v.Sql, err = services.GetSql(v.Sql, v.SqlFile)
			if err != nil {
				return nil, err
			}
			matrixCombinations, err := services.GenerateMatrixCombinations(v.Matrix)
			if err != nil {
				return nil, err
			}
			for _, matrixCombination := range matrixCombinations {
				v.Identifier.MatrixValue = services.ParamMap{}
				for key, value := range matrixCombination {
					v.Identifier.MatrixValue[key] = value
				}
				data, err := services.MarshalValidation(v)
				if err != nil {
					return nil, err
				}
				publishData = append(publishData, data)
			}
		}
	}

	return publishData, nil
}
