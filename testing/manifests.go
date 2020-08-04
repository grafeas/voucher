package vtesting

import (
	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/libtrust"
)

// NewTestManifest creates a test schema2 manifest for our mock Docker API.
func NewTestManifest() *schema2.DeserializedManifest {
	manifest := schema2.Manifest{
		Config: distribution.Descriptor{
			MediaType: schema2.MediaTypeImageConfig,
			Size:      7023,
			Digest:    "sha256:b5b2b2c507a0944348e0303114d8d93aaaa081732b86451d9bce1f432a537bc7",
		},
		Layers: []distribution.Descriptor{
			{
				MediaType: schema2.MediaTypeLayer,
				Size:      32654,
				Digest:    "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f",
			},
			{
				MediaType: schema2.MediaTypeLayer,
				Size:      16724,
				Digest:    "sha256:3c3a4604a545cdc127456d94e421cd355bca5b528f4a9c1905b15da2eb4a4c6b",
			},
			{
				MediaType: schema2.MediaTypeLayer,
				Size:      73109,
				Digest:    "sha256:ec4b8955958665577945c89419d1af06b5f7636b4ac3da7f12184802ad867736",
			},
		},
	}

	manifest.SchemaVersion = 2
	manifest.MediaType = schema2.MediaTypeManifest

	newManifest, err := schema2.FromStruct(manifest)
	if nil != err {
		panic("failed to generate new schema2 manifest")
	}

	return newManifest
}

// NewTestSchema1Manifest creates a test schema1 manifest for our mock Docker API.
func NewTestSchema1Manifest() schema1.Manifest {
	manifest := schema1.Manifest{
		Versioned:    schema1.SchemaVersion,
		Architecture: "amd64",
		Name:         "schema1image",
		History: []schema1.History{
			{
				V1Compatibility: `{
	"architecture":"amd64",
	"author":"example@example.com",
	"config":{
		"Hostname":"test",
		"Domainname":"",
		"User":"nobody",
		"AttachStdin":false,
		"AttachStdout":false,
		"AttachStderr":false,
		"Tty":false,
		"OpenStdin":false,
		"StdinOnce":false,
		"Env":[
			"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
		],
		"Cmd":null,
		"Image":"sha256:03f65aeeb2e8e8db022b297cae4cdce9248633f551452e63ba520d1f9ef2eca0",
		"Volumes":null,
		"WorkingDir":"/",
		"Entrypoint":null,
		"OnBuild":null,
		"Labels":null
	},
	"container":"12394d276740df1e762dac127e132976faf1d213af81430cf3307553886e3425",
	"container_config":{
		"Hostname":"test",
		"Domainname":"",
		"User":"nobody",
		"AttachStdin":false,
		"AttachStdout":false,
		"AttachStderr":false,
		"Tty":false,
		"OpenStdin":false,
		"StdinOnce":false,
		"Env":[
			"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
		],
		"Cmd":null,
		"Image":"sha256:03f65aeeb2e8e8db022b297cae4cdce9248633f551452e63ba520d1f9ef2eca0",
		"Volumes":null,
		"WorkingDir":"/",
		"Entrypoint":null,
		"OnBuild":null,
		"Labels":null
	},
	"created":"2020-04-09T20:09:22.704847544Z",
	"docker_version":"19.03.5",
	"id":"2d115a58c6b903363bd41f10f1fcf5005a393118c4e4862ee5bcb5f30471a15f",
	"os":"linux",
	"parent":"",
	"throwaway":true
}`,
			},
		},
		FSLayers: []schema1.FSLayer{
			{BlobSum: "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f"},
			{BlobSum: "sha256:3c3a4604a545cdc127456d94e421cd355bca5b528f4a9c1905b15da2eb4a4c6b"},
			{BlobSum: "sha256:ec4b8955958665577945c89419d1af06b5f7636b4ac3da7f12184802ad867736"},
		},
	}
	return manifest
}

// NewPrivateKey creates a private key that can be used to sign Docker
// manifests.
func NewPrivateKey() libtrust.PrivateKey {
	pk, err := libtrust.GenerateRSA2048PrivateKey()
	if nil != err {
		panic("failed to generate private key for signing test manifest")
	}
	return pk
}

// NewTestSchema1SignedManifest creates a test schema1 manifest, and signs it
// with the passed private key, for use with our mock Docker API.
func NewTestSchema1SignedManifest(pk libtrust.PrivateKey) *schema1.SignedManifest {
	m := NewTestSchema1Manifest()
	signed, err := schema1.Sign(&m, pk)
	if nil != err {
		panic("failed to sign manifest")
	}

	return signed
}
