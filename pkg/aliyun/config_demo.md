```yaml
aliyun_oss:
  access_key_id: xxxxxxxx
  access_key_secret: yyyyyyyy
  bucket: BUCKET_ID1
  region_id: cn-shanghai
  site: "https://example.com"

aliyun_sts:
  access_key_id: xxxxxxxx
  access_key_secret: yyyyyyyy
  bucket: BUCKET_ID1
  region_id: cn-shanghai
  role_arn: acs:ram::12345678:role/ROLE_NAME
  expired_seconds: 3600
```

```toml
[aliyun_oss]
access_key_id = "xxxxxxxx"
access_key_secret = "yyyyyyyy"
bucket = "BUCKET_ID1"
region_id = "cn-shanghai"
site = ""

[aliyun_sts]
access_key_id = "xxxxxxxx"
access_key_secret = "yyyyyyyy"
bucket = "BUCKET_ID1"
region_id = "cn-shanghai"
role_arn = "acs:ram::12345678:role/ROLE_NAME"
expired_seconds = 3600
```
