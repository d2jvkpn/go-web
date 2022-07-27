package aliyun

import (
	"fmt"
	// "time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type StsClient struct {
	*sts.Client
	config Config
}

type StsResult struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SecurityToken   string `json:"securityToken"`
	Expiration      string `json:"expiration"`
	RegionId        string `json:"regionId"`
	Bucket          string `json:"bucket"`
}

func (client *StsClient) AssumeRole(userId string) (response *sts.AssumeRoleResponse, err error) {
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"

	request.RoleArn = client.config.RoleArn
	request.RoleSessionName = userId // 用于统计
	if client.config.ExpiredSeconds > 0 {
		request.DurationSeconds = requests.NewInteger(client.config.ExpiredSeconds)
	}

	if response, err = client.Client.AssumeRole(request); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *StsClient) GetSTS(userId, key string) (result *StsResult, err error) {
	var response *sts.AssumeRoleResponse

	if response, err = client.AssumeRole(userId); err != nil {
		return nil, err
	}

	result = &StsResult{
		AccessKeyId:     response.Credentials.AccessKeyId,
		AccessKeySecret: response.Credentials.AccessKeySecret,
		SecurityToken:   response.Credentials.SecurityToken,
		Expiration:      response.Credentials.Expiration,
		RegionId:        client.config.RegionId,
		Bucket:          client.config.Bucket,
	}

	return result, nil
}

func (result *StsResult) Upload(fp, subpath string, options ...oss.Option) (
	link string, err error) {
	var (
		urlpath string
		bucket  *oss.Bucket
		client  *oss.Client
	)

	if client, err = oss.New(
		fmt.Sprintf("https://oss-%s.aliyuncs.com", result.RegionId),
		result.AccessKeyId, result.AccessKeySecret,
		oss.SecurityToken(result.SecurityToken)); err != nil {
		return "", err
	}

	if bucket, err = client.Bucket(result.Bucket); err != nil {
		return "", err
	}

	// urlpath = strings.Trim(fmt.Sprintf("%s/%s", strings.Trim(result.Path, "/"), subpath), "/")
	if subpath, err = ValidSubpath(subpath); err != nil {
		return "", err
	}
	if err = bucket.PutObjectFromFile(urlpath, fp, options...); err != nil {
		return "", err
	}

	// https://fileserver-cim.oss-cn-hangzhou.aliyuncs.com/meshes/PrivateModels/hello.txt
	link = fmt.Sprintf("https://%s.oss-%s.aliyuncs.com/%s", result.Bucket, result.RegionId, urlpath)
	return link, nil
}

func (client *StsClient) Upload(userId string, fp, subpath string, options ...oss.Option) (
	link string, err error) {
	var result *StsResult

	if result, err = client.GetSTS(userId, ""); err != nil {
		return "", err
	}
	return result.Upload(fp, subpath, options...)
}
