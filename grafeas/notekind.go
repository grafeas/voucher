package grafeas

import (
	"github.com/Shopify/voucher"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/common"
)

// getNoteKind translates the Voucher MetadataType into a Google Container Analysis NoteKind.
func getNoteKind(metadataType voucher.MetadataType) common.NoteKind {
	switch metadataType {
	case voucher.VulnerabilityType:
		return common.NoteKind_VULNERABILITY
	case voucher.BuildDetailsType:
		return common.NoteKind_BUILD
	case voucher.AttestationType:
		return common.NoteKind_ATTESTATION
	case DiscoveryType:
		return common.NoteKind_DISCOVERY
	case PackageType:
		return common.NoteKind_PACKAGE
	case ImageType:
		return common.NoteKind_IMAGE
	case DeploymentType:
		return common.NoteKind_DEPLOYMENT
	}
	return common.NoteKind_NOTE_KIND_UNSPECIFIED
}
