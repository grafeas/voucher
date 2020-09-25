package objects

import "time"

//BuildSignatureKeyType https://github.com/grafeas/client-go/blob/master/0.1.0/model_build_signature_key_type.go
type BuildSignatureKeyType string

//AliasContextKind https://github.com/grafeas/client-go/blob/master/0.1.0/model_alias_context_kind.go
type AliasContextKind string

//HashType https://github.com/grafeas/client-go/blob/master/0.1.0/model_hash_hash_type.go
type HashType string

//consts
const (
	BuildSignatureUnspecified BuildSignatureKeyType = "KEY_TYPE_UNSPECIFIED"
	BuildSignaturePgpASCII    BuildSignatureKeyType = "PGP_ASCII_ARMORED"
	BuildSignaturePkixPem     BuildSignatureKeyType = "PKIX_PEM"

	AliasContextKindUnspecified AliasContextKind = "KIND_UNSPECIFIED"
	AliasContextKindFixed       AliasContextKind = "FIXED"
	AliasContextKindMovable     AliasContextKind = "MOVABLE"
	AliasContextKindOther       AliasContextKind = "OTHER"

	HashTypeUnspecified HashType = "HASH_TYPE_UNSPECIFIED"
	HashTypeSHA256      HashType = "SHA256"
)

//note objects

//Build https://github.com/grafeas/client-go/blob/master/0.1.0/model_build_build.go
type Build struct {
	BuilderVersion string          `json:"builderVersion,omitempty"` //required
	Signature      *BuildSignature `json:"signature,omitempty"`
}

//BuildSignature https://github.com/grafeas/client-go/blob/master/0.1.0/model_build_build_signature.go
type BuildSignature struct {
	PublicKey string                 `json:"publicKey,omitempty"`
	Signature string                 `json:"signature,omitempty"` //required
	KeyID     string                 `json:"keyId,omitempty"`
	KeyType   *BuildSignatureKeyType `json:"keyType,omitempty"`
}

//occurrence objects

//BuildDetails https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1build_details.go
type BuildDetails struct {
	Provenance      *ProvenanceBuild `json:"provenance,omitempty"` //required
	ProvenanceBytes string           `json:"provenanceBytes,omitempty"`
}

//ProvenanceBuild https://github.com/grafeas/client-go/blob/master/0.1.0/model_provenance_build_provenance.go
type ProvenanceBuild struct {
	ID               string               `json:"id,omitempty"` //required
	ProjectID        string               `json:"projectId,omitempty"`
	Commands         []ProvenanceCommand  `json:"commands,omitempty"`
	BuiltArtifacts   []ProvenanceArtifact `json:"builtArtifacts,omitempty"`
	CreateTime       time.Time            `json:"createTime,omitempty"`
	StartTime        time.Time            `json:"startTime,omitempty"`
	EndTime          time.Time            `json:"endTime,omitempty"`
	Creator          string               `json:"creator,omitempty"` //email address
	LogsURI          string               `json:"logsUri,omitempty"`
	SourceProvenance *ProvenanceSource    `json:"sourceProvenance,omitempty"`
	TriggerID        string               `json:"triggerId,omitempty"`
	BuildOptions     map[string]string    `json:"buildOptions,omitempty"`
	BuilderVersion   string               `json:"builderVersion,omitempty"`
}

//ProvenanceArtifact https://github.com/grafeas/client-go/blob/master/0.1.0/model_provenance_artifact.go
type ProvenanceArtifact struct {
	Checksum string   `json:"checksum,omitempty"`
	ID       string   `json:"id,omitempty"`
	Names    []string `json:"names,omitempty"`
}

//ProvenanceSource https://github.com/grafeas/client-go/blob/master/0.1.0/model_provenance_source.go
type ProvenanceSource struct {
	ArtifactStorageSourceURI string                          `json:"artifactStorageSourceUri,omitempty"`
	FileHashes               map[string]ProvenanceFileHashes `json:"fileHashes,omitempty"`
	Context                  *SourceContext                  `json:"context,omitempty"`
	AdditionalContexts       []SourceContext                 `json:"additionalContexts,omitempty"`
}

//SourceContext https://github.com/grafeas/client-go/blob/master/0.1.0/model_source_source_context.go
type SourceContext struct {
	CloudRepo *CloudRepoSourceContext `json:"cloudRepo,omitempty"`
	Gerrit    *GerritSourceContext    `json:"gerrit,omitempty"`
	Git       *GitSourceContext       `json:"git,omitempty"`
	Labels    map[string]string       `json:"labels,omitempty"`
}

//CloudRepoSourceContext https://github.com/grafeas/client-go/blob/master/0.1.0/model_source_cloud_repo_source_context.go
type CloudRepoSourceContext struct { //SourceCloudRepoSourceContext
	RepoID       *SourceRepoID       `json:"repoId,omitempty"`
	RevisionID   string              `json:"revisionId,omitempty"`
	AliasContext *SourceAliasContext `json:"aliasContext,omitempty"`
}

//GitSourceContext https://github.com/grafeas/client-go/blob/master/0.1.0/model_source_git_source_context.go
type GitSourceContext struct { //SourceGitSourceContext
	URL        string `json:"url,omitempty"`
	RevisionID string `json:"revisionId,omitempty"`
}

//SourceRepoID https://github.com/grafeas/client-go/blob/master/0.1.0/model_source_repo_id.go
type SourceRepoID struct {
	ProjectRepoID *SourceProjectRepoID `json:"projectRepoId,omitempty"`
	UID           string               `json:"uid,omitempty"`
}

//SourceAliasContext https://github.com/grafeas/client-go/blob/master/0.1.0/model_source_alias_context.go
type SourceAliasContext struct {
	Kind *AliasContextKind `json:"kind,omitempty"`
	Name string            `json:"name,omitempty"`
}

//SourceProjectRepoID https://github.com/grafeas/client-go/blob/master/0.1.0/model_source_project_repo_id.go
type SourceProjectRepoID struct {
	ProjectID string `json:"projectId,omitempty"`
	RepoName  string `json:"repoNme,omitempty"`
}

//ProvenanceCommand https://github.com/grafeas/client-go/blob/master/0.1.0/model_provenance_command.go
type ProvenanceCommand struct {
	Name    string   `json:"name,omitempty"` //required
	Env     []string `json:"env,omitempty"`
	Args    []string `json:"args,omitempty"`
	Dir     string   `json:"dir,omitempty"`
	ID      string   `json:"id,omitempty"`
	WaitFor []string `json:"waitFor,omitempty"`
}

//ProvenanceFileHashes https://github.com/grafeas/client-go/blob/master/0.1.0/model_provenance_file_hashes.go
type ProvenanceFileHashes struct {
	FileHash []ProvenanceHash `json:"fileHash,omitempty"`
}

//ProvenanceHash https://github.com/grafeas/client-go/blob/master/0.1.0/model_provenance_hash.go
type ProvenanceHash struct {
	Type  *HashType `json:"type,omitempty"`  //required
	Value string    `json:"value,omitempty"` //required
}

//GerritSourceContext https://github.com/grafeas/client-go/blob/master/0.1.0/model_source_gerrit_source_context.go
type GerritSourceContext struct {
	HostURI       string              `json:"hostUri,omitempty"`
	GerritProject string              `json:"gerritProject,omitempty"`
	RevisionID    string              `json:"revisionId,omitempty"`
	AliasContext  *SourceAliasContext `json:"aliasContext,omitempty"`
}
