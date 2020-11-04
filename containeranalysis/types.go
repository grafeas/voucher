package containeranalysis

import "github.com/grafeas/voucher"

// DiscoveryType is a Grafeas specific type which refers to MetadataItems containing metadata discovery status.
const DiscoveryType voucher.MetadataType = "discovery"

// PackageType is a Grafeas specific type which refers to MetadataItems containing package information.
const PackageType voucher.MetadataType = "package"

// ImageType is a Grafeas specific type which refers to MetadataItems containing Image information.
const ImageType voucher.MetadataType = "image"

// DeploymentType is a Grafeas specific type which refers to MetadataItems containing deployment data.
const DeploymentType voucher.MetadataType = "deployment"
