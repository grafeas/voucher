package grafeasos

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	grafeas "github.com/grafeas/client-go/0.1.0"
)

// PopulateGrafeasTestData populates grafeas with created notes and occurrences
func PopulateGrafeasTestData(basePath string) error {
	project := "projects/grafeasclienttest" // notes and occurrences are created in the same project
	config := grafeas.NewConfiguration()
	if basePath != "" {
		config.BasePath = basePath
	} else {
		config.BasePath = "http://localhost:8080"
	}
	gcli := grafeas.NewAPIClient(config)
	notes, occs := createNotesAndOccurrences(project)

	err := interpretResponse(
		gcli.GrafeasV1Beta1Api.BatchCreateNotes(
			context.Background(),
			project,
			grafeas.V1beta1BatchCreateNotesRequest{Parent: project, Notes: notes},
		),
	)
	if err != nil {
		return err
	}

	err = interpretResponse(
		gcli.GrafeasV1Beta1Api.BatchCreateOccurrences(
			context.Background(),
			project,
			grafeas.V1beta1BatchCreateOccurrencesRequest{Parent: project, Occurrences: occs},
		),
	)
	return err
}

func interpretResponse(result interface{}, response *http.Response, err error) error {
	if err != nil {
		fmt.Println("error: ", err)
		return err
	} else if response.StatusCode != 200 {
		fmt.Println(response)
		return errors.New("error: " + strconv.Itoa(response.StatusCode))
	} else {
		fmt.Printf("Successfully created: %T\n", result)
		return nil
	}
}

func createNotesAndOccurrences(parent string) (map[string]grafeas.V1beta1Note, []grafeas.V1beta1Occurrence) {
	notes := make(map[string]grafeas.V1beta1Note)
	var occurrences []grafeas.V1beta1Occurrence
	//vulnerability note
	note, occs := createVulnerabilityNote(parent)
	notes[note.Name] = note
	occurrences = append(occurrences, occs...)
	//build note
	note, occs = createBuildNote(parent)
	notes[note.Name] = note
	occurrences = append(occurrences, occs...)
	//attestation note
	note, occs = createAttestationNote(parent)
	notes[note.Name] = note
	occurrences = append(occurrences, occs...)
	//discovery note
	note, occs = createDiscoveryNote(parent)
	notes[note.Name] = note
	occurrences = append(occurrences, occs...)

	return notes, occurrences
}

func createGenericNote(noteKind *grafeas.V1beta1NoteKind, name string) grafeas.V1beta1Note {
	return grafeas.V1beta1Note{Name: name, ShortDescription: "short", LongDescription: "long", Kind: noteKind}
}

func createGenericOccurrence(noteKind *grafeas.V1beta1NoteKind, name, noteName, parent string) grafeas.V1beta1Occurrence {
	return grafeas.V1beta1Occurrence{
		Name:     getOccurrenceName(parent, "occurrences", name),
		Resource: &grafeas.V1beta1Resource{Uri: "https://gcr.io/project/image@sha256:foo"},
		NoteName: getOccurrenceName(parent, "notes", noteName), Kind: noteKind,
	}
}

func getOccurrenceName(parent, nameType, noteName string) string {
	return parent + "/" + nameType + "/" + noteName
}

func createDiscoveryNote(parent string) (grafeas.V1beta1Note, []grafeas.V1beta1Occurrence) {
	noteKindVuln := grafeas.VULNERABILITY_V1beta1NoteKind
	noteKind := grafeas.DISCOVERY_V1beta1NoteKind
	note := createGenericNote(&noteKind, "grafeasdiscovery")
	note.Discovery = &grafeas.DiscoveryDiscovery{AnalysisKind: &noteKindVuln}

	activeContAnalysis := grafeas.ACTIVE_DiscoveredContinuousAnalysis
	successStatus := grafeas.FINISHED_SUCCESS_DiscoveredAnalysisStatus
	finishedDiscoveryOcc := createGenericOccurrence(&noteKind, "occurdiscovery1", note.Name, parent)
	finishedDiscoveryOcc.Discovered = &grafeas.V1beta1discoveryDetails{
		Discovered: &grafeas.DiscoveryDiscovered{
			ContinuousAnalysis: &activeContAnalysis,
			AnalysisStatus:     &successStatus,
		},
	}

	packageKind := grafeas.NORMAL_VersionVersionKind
	vulnerabilityOcc := createGenericOccurrence(&noteKindVuln, "occurdiscovervuln", note.Name, parent)
	vulnerabilityOcc.Vulnerability = &grafeas.V1beta1vulnerabilityDetails{
		Type_: "vulnlow",
		PackageIssue: []grafeas.VulnerabilityPackageIssue{
			{
				AffectedLocation: &grafeas.VulnerabilityVulnerabilityLocation{
					CpeUri:   "uri",
					Package_: "package_test",
					Version: &grafeas.PackageVersion{
						Name:     "v0.1.0",
						Kind:     &packageKind,
						Revision: "re",
					},
				},
			},
		},
	}

	pendingStatus := grafeas.PENDING_DiscoveredAnalysisStatus
	pendingDiscoveryOcc := createGenericOccurrence(&noteKind, "occurdiscovery2", note.Name, parent)
	pendingDiscoveryOcc.Discovered = &grafeas.V1beta1discoveryDetails{
		Discovered: &grafeas.DiscoveryDiscovered{
			ContinuousAnalysis: &activeContAnalysis,
			AnalysisStatus:     &pendingStatus,
		},
	}

	return note, []grafeas.V1beta1Occurrence{finishedDiscoveryOcc, vulnerabilityOcc, pendingDiscoveryOcc}
}

func createAttestationNote(parent string) (grafeas.V1beta1Note, []grafeas.V1beta1Occurrence) {
	noteKind := grafeas.ATTESTATION_V1beta1NoteKind
	note := createGenericNote(&noteKind, "grafeasattestation")
	note.AttestationAuthority = &grafeas.AttestationAuthority{Hint: &grafeas.AuthorityHint{HumanReadableName: note.Name}}

	contentType := grafeas.CONTENT_TYPE_UNSPECIFIED_AttestationPgpSignedAttestationContentType
	occ := createGenericOccurrence(&noteKind, "occurdiscovery", note.Name, parent)
	occ.Attestation = &grafeas.V1beta1attestationDetails{
		Attestation: &grafeas.AttestationAttestation{
			PgpSignedAttestation: &grafeas.AttestationPgpSignedAttestation{
				Signature:   "signature",
				PgpKeyId:    "1234",
				ContentType: &contentType,
			},
		},
	}

	return note, []grafeas.V1beta1Occurrence{occ}
}

func createBuildNote(parent string) (grafeas.V1beta1Note, []grafeas.V1beta1Occurrence) {
	noteKind := grafeas.BUILD_V1beta1NoteKind
	note := createGenericNote(&noteKind, "grafeasbuild")
	note.Build = &grafeas.BuildBuild{BuilderVersion: "v0.0.0"}

	occ1 := createGenericOccurrence(&noteKind, "occurbuild1", note.Name, parent)
	occ1.Build = &grafeas.V1beta1buildDetails{
		Provenance: &grafeas.ProvenanceBuildProvenance{
			Id:        "id",
			Creator:   "shopify",
			ProjectId: "id1",
			LogsUri:   "some_url",
			SourceProvenance: &grafeas.ProvenanceSource{
				Context: &grafeas.SourceSourceContext{
					Git: &grafeas.SourceGitSourceContext{
						Url:        "github.com/Shopify/voucher",
						RevisionId: "q1q",
					},
				},
			},
		},
	}
	occ2 := createGenericOccurrence(&noteKind, "occurbuild2", note.Name, parent)
	occ2.Build = &grafeas.V1beta1buildDetails{
		Provenance: &grafeas.ProvenanceBuildProvenance{
			Id:        "provenceid",
			Creator:   "shopify",
			ProjectId: "id2",
			LogsUri:   "some_url2",
			BuiltArtifacts: []grafeas.ProvenanceArtifact{
				{Checksum: "fe43gf42f2", Id: "some_id", Names: []string{"name1", "name2"}},
			},
			SourceProvenance: &grafeas.ProvenanceSource{
				Context: &grafeas.SourceSourceContext{
					Git: &grafeas.SourceGitSourceContext{
						Url:        "github.com/Shopify/voucher",
						RevisionId: "2",
					},
				},
			},
		},
	}

	return note, []grafeas.V1beta1Occurrence{occ1, occ2}
}

func createVulnerabilityNote(parent string) (grafeas.V1beta1Note, []grafeas.V1beta1Occurrence) {
	vulnSeverity := grafeas.MEDIUM_VulnerabilitySeverity
	noteKind := grafeas.VULNERABILITY_V1beta1NoteKind
	note := createGenericNote(&noteKind, "grafeasvulnerability")
	note.Vulnerability = &grafeas.VulnerabilityVulnerability{
		CvssScore: 4.3,
		Severity:  &vulnSeverity,
		Details: []grafeas.VulnerabilityDetail{
			{CpeUri: "test_url", Package_: "package", SeverityName: "medium"},
		},
		CvssV3: nil, WindowsDetails: []grafeas.VulnerabilityWindowsDetail{},
	}

	packageKind := grafeas.NORMAL_VersionVersionKind
	occ := createGenericOccurrence(&noteKind, "occurvulnerability", note.Name, parent)
	occ.Vulnerability = &grafeas.V1beta1vulnerabilityDetails{
		Type_:    "vulnmedium",
		Severity: &vulnSeverity,
		PackageIssue: []grafeas.VulnerabilityPackageIssue{
			{
				AffectedLocation: &grafeas.VulnerabilityVulnerabilityLocation{
					CpeUri:   "uri",
					Package_: "package_test",
					Version:  &grafeas.PackageVersion{Name: "v0.0.0", Kind: &packageKind, Revision: "r"},
				},
			},
		},
	}

	return note, []grafeas.V1beta1Occurrence{occ}
}
