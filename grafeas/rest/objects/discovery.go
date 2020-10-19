package objects

//DiscoveredAnalysisStatus based on
//https://github.com/grafeas/client-go/blob/master/0.1.0/model_discovered_analysis_status.go
type DiscoveredAnalysisStatus string

//consts
const (
	DiscoveredAnalysisStatusUnspecified         DiscoveredAnalysisStatus = "ANALYSIS_STATUS_UNSPECIFIED"
	DiscoveredAnalysisStatusPending             DiscoveredAnalysisStatus = "PENDING"
	DiscoveredAnalysisStatusScanning            DiscoveredAnalysisStatus = "SCANNING"
	DiscoveredAnalysisStatusFinishedSuccess     DiscoveredAnalysisStatus = "FINISHED_SUCCESS"
	DiscoveredAnalysisStatusFinishedFailed      DiscoveredAnalysisStatus = "FINISHED_FAILED"
	DiscoveredAnalysisStatusFinishedUnsupported DiscoveredAnalysisStatus = "FINISHED_UNSUPPORTED"
)

//discovery for occurrence

//DiscoveryDetails based on
//https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1discovery_details.go
type DiscoveryDetails struct {
	Discovered *DiscoveryDiscovered `json:"discovered,omitempty"` //required
}

//DiscoveryDiscovered based on
//https://github.com/grafeas/client-go/blob/master/0.1.0/model_discovery_discovered.go
type DiscoveryDiscovered struct {
	AnalysisStatus *DiscoveredAnalysisStatus `json:"analysisStatus,omitempty"`
}

//discovery for note

//Discovery based on
//https://github.com/grafeas/client-go/blob/master/0.1.0/model_discovery_discovery.go
type Discovery struct {
	AnalysisKind *NoteKind `json:"analysisKind,omitempty"` //required
}
