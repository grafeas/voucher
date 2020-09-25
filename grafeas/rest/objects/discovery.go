package objects

//DiscoveredAnalysisStatus https://github.com/grafeas/client-go/blob/master/0.1.0/model_discovered_analysis_status.go
type DiscoveredAnalysisStatus string

//DiscoveredContinuousAnalysis https://github.com/grafeas/client-go/blob/master/0.1.0/model_discovered_continuous_analysis.go
type DiscoveredContinuousAnalysis string

//consts
const (
	DiscoveredContinuousAnalysisUnspecified DiscoveredContinuousAnalysis = "CONTINUOUS_ANALYSIS_UNSPECIFIED"
	DiscoveredContinuousAnalysisActive      DiscoveredContinuousAnalysis = "ACTIVE"
	DiscoveredContinuousAnalysisInactive    DiscoveredContinuousAnalysis = "INACTIVE"

	DiscoveredAnalysisStatusUnspecified         DiscoveredAnalysisStatus = "ANALYSIS_STATUS_UNSPECIFIED"
	DiscoveredAnalysisStatusPending             DiscoveredAnalysisStatus = "PENDING"
	DiscoveredAnalysisStatusScanning            DiscoveredAnalysisStatus = "SCANNING"
	DiscoveredAnalysisStatusFinishedSuccess     DiscoveredAnalysisStatus = "FINISHED_SUCCESS"
	DiscoveredAnalysisStatusFinishedFailed      DiscoveredAnalysisStatus = "FINISHED_FAILED"
	DiscoveredAnalysisStatusFinishedUnsupported DiscoveredAnalysisStatus = "FINISHED_UNSUPPORTED"
)

//discovery for occurrence

//DiscoveryDetails https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1discovery_details.go
type DiscoveryDetails struct {
	Discovered *DiscoveryDiscovered `json:"discovered,omitempty"` //required
}

//DiscoveryDiscovered https://github.com/grafeas/client-go/blob/master/0.1.0/model_discovery_discovered.go
type DiscoveryDiscovered struct {
	ContinuousAnalysis  *DiscoveredContinuousAnalysis `json:"continuousAnalysis,omitempty"`
	AnalysisStatus      *DiscoveredAnalysisStatus     `json:"analysisStatus,omitempty"`
	AnalysisStatusError *RPCStatus                    `json:"analysisStatusError,omitempty"`
}

//discovery for note

//Discovery https://github.com/grafeas/client-go/blob/master/0.1.0/model_discovery_discovery.go
type Discovery struct {
	AnalysisKind *NoteKind `json:"analysisKind,omitempty"` //required
}
