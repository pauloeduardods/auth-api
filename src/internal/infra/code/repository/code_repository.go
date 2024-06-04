package code_repo

import (
	"auth-api/src/internal/domain/code"
	"auth-api/src/pkg/logger"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type CodeRepositoryDynamoDB struct {
	dynamoDBClient *dynamodb.Client
	tableName      string
	logger         logger.Logger
}

func NewCodeRepositoryDynamoDB(tableName string, dynamoClient *dynamodb.Client, logger logger.Logger) *CodeRepositoryDynamoDB {
	return &CodeRepositoryDynamoDB{
		dynamoDBClient: dynamoClient,
		tableName:      tableName,
		logger:         logger,
	}
}

func (r *CodeRepositoryDynamoDB) Save(input *code.Code) error {
	_, err := r.dynamoDBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item: map[string]types.AttributeValue{
			"identifier": &types.AttributeValueMemberS{Value: input.Identifier},
			"value":      &types.AttributeValueMemberS{Value: input.Value},
			"expires_at": &types.AttributeValueMemberS{Value: input.ExpiresAt.Format(time.RFC3339)},
		},
	})
	return err
}

func (r *CodeRepositoryDynamoDB) FindByIdentifier(identifier string) (*[]code.Code, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(r.tableName),
		KeyConditions: map[string]types.Condition{
			"identifier": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: identifier},
				},
			},
		},
	}

	result, err := r.dynamoDBClient.Query(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, code.ErrCodeNotFound
	}

	codes := make([]code.Code, 0)
	for _, item := range result.Items {
		expiresAt, err := time.Parse(time.RFC3339, item["expires_at"].(*types.AttributeValueMemberS).Value)
		if err != nil {
			r.logger.Error("failed to parse time: %v", err)
			continue
		}
		codes = append(codes, code.Code{
			Value:      item["value"].(*types.AttributeValueMemberS).Value,
			ExpiresAt:  expiresAt,
			Identifier: item["identifier"].(*types.AttributeValueMemberS).Value,
		})
	}

	return &codes, nil
}

func (r *CodeRepositoryDynamoDB) Delete(code *code.Code) error {
	_, err := r.dynamoDBClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"identifier": &types.AttributeValueMemberS{Value: code.Identifier},
		},
	})
	return err
}
