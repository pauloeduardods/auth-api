package code

import (
	"auth-api/src/internal/shared/code/domain/code"
	"auth-api/src/pkg/logger"
	"context"
	"fmt"
	"strconv"
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

func NewCodeRepositoryDynamoDB(tableName string, dynamoClient *dynamodb.Client, logger logger.Logger) code.CodeRepository {
	return &CodeRepositoryDynamoDB{
		dynamoDBClient: dynamoClient,
		tableName:      tableName,
		logger:         logger,
	}
}

func (r *CodeRepositoryDynamoDB) Save(ctx context.Context, input *code.Code) error {
	_, err := r.dynamoDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item: map[string]types.AttributeValue{
			"identifier": &types.AttributeValueMemberS{Value: input.Identifier},
			"code":       &types.AttributeValueMemberS{Value: input.Value},
			"expires_at": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", input.ExpiresAt.Unix())},
		},
	})
	return err
}

func (r *CodeRepositoryDynamoDB) FindCode(ctx context.Context, identifier, codeValue string) (*code.Code, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(r.tableName),
		KeyConditions: map[string]types.Condition{
			"identifier": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: identifier},
				},
			},
			"code": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: codeValue},
				},
			},
		},
	}

	result, err := r.dynamoDBClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, code.ErrCodeNotFound
	}

	item := result.Items[0]
	expiresAtUnix, err := strconv.ParseInt(item["expires_at"].(*types.AttributeValueMemberN).Value, 10, 64)
	if err != nil {
		r.logger.Error("failed to parse time: %v", err)
		return nil, err
	}
	expiresAt := time.Unix(expiresAtUnix, 0)
	return &code.Code{
		Value:      item["code"].(*types.AttributeValueMemberS).Value,
		ExpiresAt:  expiresAt,
		Identifier: item["identifier"].(*types.AttributeValueMemberS).Value,
	}, nil
}

func (r *CodeRepositoryDynamoDB) Delete(ctx context.Context, code *code.Code) error {
	_, err := r.dynamoDBClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"identifier": &types.AttributeValueMemberS{Value: code.Identifier},
			"code":       &types.AttributeValueMemberS{Value: code.Value},
		},
	})
	return err
}
