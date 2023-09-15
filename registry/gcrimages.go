package registry

import (
	"github.com/pkg/errors"
	"fmt"
	"net/http"
	// "net/http"
	// "strconv"
	"crypto/tls"
	"runtime"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	// v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)


// https://github.com/tap8stry/orion/blob/911ba084852eda91b31b59bc718bcda6d6bb99b3/pkg/imagefs/imagefs.go#L52
// getImageReferences gets a reference string and returns all image
func GetImageReferences(imageName string) ([]struct {
	Digest string
	Arch   string
	OS     string
}, error) {
	//auth remote.Option
	auth := remote.WithAuth(authn.FromConfig(authn.AuthConfig{
		Username: user,
		Password: pass,
	}))
	// https, skip_key_validate
	if nil==transport { //
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	trans := remote.WithTransport(transport)

	ref, err := name.ParseReference(imageName)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing image reference %s", imageName)
	}
	//descr, err := remote.Get(ref)
	descr, err := remote.Get(ref, auth, trans) //, remote.WithAuthFromKeychain(authn.DefaultKeychain)
	if err != nil {
		return nil, errors.Wrap(err, "fetching remote descriptor")
	}

	images := []struct {
		Digest string
		Arch   string
		OS     string
	}{}

	// If we got a digest, we reuse it as is
	if _, ok := ref.(name.Digest); ok {
		images = append(images, struct {
			Digest string
			Arch   string
			OS     string
		}{Digest: ref.(name.Digest).String(), OS: runtime.GOOS, Arch:runtime.GOARCH}) //TODO 采用registry的arch?
		return images, nil
	}

	// If the reference is not an image, it has to work as a tag
	tag, ok := ref.(name.Tag)
	if !ok {
		return nil, errors.Errorf("could not cast tag from reference %s", imageName)
	}
	// If the reference points to an image, return it
	if descr.MediaType.IsImage() {
		fmt.Printf("Reference %s points to a single image", imageName)
		// Check if we can get an image
		im, err := descr.Image()
		if err != nil {
			return nil, errors.Wrap(err, "getting image from descriptor")
		}

		imageDigest, err := im.Digest()
		if err != nil {
			return nil, errors.Wrap(err, "while calculating image digest")
		}

		dig, err := name.NewDigest(
			fmt.Sprintf(
				"%s/%s@%s:%s",
				tag.RegistryStr(), tag.RepositoryStr(),
				imageDigest.Algorithm, imageDigest.Hex,
			),
		)
		if err != nil {
			return nil, errors.Wrap(err, "building single image digest")
		}

		fmt.Printf("Adding image digest %s from reference", dig.String())
		return append(images, struct {
			Digest string
			Arch   string
			OS     string
		}{Digest: dig.String(), OS: runtime.GOOS, Arch:runtime.GOARCH}), nil //registry's arch?
	}

	// Get the image index
	index, err := descr.ImageIndex()
	if err != nil {
		return nil, errors.Wrapf(err, "getting image index for %s", imageName)
	}
	indexManifest, err := index.IndexManifest()
	if err != nil {
		return nil, errors.Wrapf(err, "getting index manifest from %s", imageName)
	}
	msg:= fmt.Sprintf("Reference image index points to %d manifests", len(indexManifest.Manifests))
	fmt.Println(msg)

	for _, manifest := range indexManifest.Manifests {
		dig, err := name.NewDigest(
			fmt.Sprintf(
				"%s/%s@%s:%s",
				tag.RegistryStr(), tag.RepositoryStr(),
				manifest.Digest.Algorithm, manifest.Digest.Hex,
			))
		if err != nil {
			return nil, errors.Wrap(err, "generating digest for image")
		}

		arch, osid := "", ""
		if manifest.Platform != nil {
			arch = manifest.Platform.Architecture
			osid = manifest.Platform.OS
		}
		msg:= fmt.Sprintf(
			"Adding image %s/%s@%s:%s (%s/%s)",
			tag.RegistryStr(), tag.RepositoryStr(), manifest.Digest.Algorithm, manifest.Digest.Hex,
			arch, osid, //avoid: manifest.Platform == nil
		)
		fmt.Println(msg)
		if "unknown"!=arch { //attestation@buildx v1.10+ >> https://docs.docker.com/build/attestations/
			images = append(images,
				struct {
					Digest string
					Arch   string
					OS     string
				}{
					Digest: dig.String(),
					Arch:   arch,
					OS:     osid,
				})
		}
	}
	return images, nil
}
