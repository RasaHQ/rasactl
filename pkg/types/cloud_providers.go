package types

type CloudProvider string

const (
	CloudProviderUnknown      CloudProvider = "Unknown"
	CloudProviderGoogle       CloudProvider = "GCP"
	CloudProviderAmazon       CloudProvider = "AWS"
	CloudProviderAlibaba      CloudProvider = "Alibaba Cloud"
	CloudProviderAzure        CloudProvider = "Azure"
	CloudProviderDigitalOcean CloudProvider = "Digital Ocean"
	CloudProviderOracle       CloudProvider = "Oracle Cloud"
)
