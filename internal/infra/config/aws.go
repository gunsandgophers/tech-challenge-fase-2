package config

var (
	AWS_REGION       = GetEnv("AWS_REGION", "us-east-1")
	AWS_USER_POOL_ID = GetEnv("AWS_USER_POOL_ID", "us-east-1_3ofqHwfxr")
)
