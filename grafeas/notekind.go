package grafeas

import (
	"github.com/Shopify/voucher"
	containeranalysispb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1alpha1"
)

// getNoteKind translates the Voucher MetadataType into a Google Container Analysis NoteKind.
func getNoteKind(metadataType voucher.MetadataType) containeranalysispb.Note_Kind {
	switch metadataType {
	case voucher.VulnerabilityType:
		return containeranalysispb.Note_PACKAGE_VULNERABILITY
	case DiscoveryType:
		return containeranalysispb.Note_DISCOVERY
	case voucher.BuildDetailsType:
		return containeranalysispb.Note_BUILD_DETAILS
	}
	return containeranalysispb.Note_KIND_UNSPECIFIED
}
