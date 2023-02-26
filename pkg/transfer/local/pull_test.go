/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package local

import (
	"testing"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/pkg/transfer"
	"github.com/containerd/containerd/pkg/unpack"
	"github.com/containerd/containerd/platforms"
)

func TestGetSupportedPlatform(t *testing.T) {
	supportedPlatforms := []unpack.Platform{
		{
			Platform:       platforms.OnlyStrict(platforms.MustParse("linux/amd64")),
			SnapshotterKey: "native",
		},
		{
			Platform:       platforms.OnlyStrict(platforms.MustParse("linux/amd64")),
			SnapshotterKey: containerd.DefaultSnapshotter,
		},
		{
			Platform:       platforms.OnlyStrict(platforms.MustParse("linux/amd64")),
			SnapshotterKey: "devmapper",
		},
		{
			Platform:       platforms.OnlyStrict(platforms.MustParse("linux/arm64")),
			SnapshotterKey: containerd.DefaultSnapshotter,
		},
		{
			Platform:       platforms.OnlyStrict(platforms.MustParse("linux/arm64")),
			SnapshotterKey: "native",
		},
		{
			Platform:       platforms.OnlyStrict(platforms.MustParse("linux/arm")),
			SnapshotterKey: "native",
		},
		{
			Platform:       platforms.OnlyStrict(platforms.MustParse("linux/arm")),
			SnapshotterKey: containerd.DefaultSnapshotter,
		},
	}

	for _, testCase := range []struct {
		// Name is the name of the test
		Name string

		//Input
		UnpackConfig       transfer.UnpackConfiguration
		SupportedPlatforms []unpack.Platform

		//Expected
		Match            bool
		ExpectedPlatform transfer.UnpackConfiguration
	}{
		{
			Name: "No match input linux/arm64 and devmapper",
			UnpackConfig: transfer.UnpackConfiguration{
				Platform:    platforms.MustParse("linux/arm64"),
				Snapshotter: "devmapper",
			},
			SupportedPlatforms: supportedPlatforms,
			Match:              false,
			ExpectedPlatform:   transfer.UnpackConfiguration{},
		},
		{
			Name: "No match input linux/386 and defaultSnapshotter",
			UnpackConfig: transfer.UnpackConfiguration{
				Platform:    platforms.MustParse("linux/386"),
				Snapshotter: containerd.DefaultSnapshotter,
			},
			SupportedPlatforms: supportedPlatforms,
			Match:              false,
			ExpectedPlatform:   transfer.UnpackConfiguration{},
		},
		{
			Name: "Match linux/amd64 and native snapshotter",
			UnpackConfig: transfer.UnpackConfiguration{
				Platform:    platforms.MustParse("linux/amd64"),
				Snapshotter: "native",
			},
			SupportedPlatforms: supportedPlatforms,
			Match:              true,
			ExpectedPlatform: transfer.UnpackConfiguration{
				Platform:    platforms.MustParse("linux/amd64"),
				Snapshotter: "native",
			},
		},
		{
			Name: "Match linux/arm64 and native snapshotter",
			UnpackConfig: transfer.UnpackConfiguration{
				Platform:    platforms.MustParse("linux/arm64"),
				Snapshotter: "native",
			},
			SupportedPlatforms: supportedPlatforms,
			Match:              true,
			ExpectedPlatform: transfer.UnpackConfiguration{
				Platform:    platforms.MustParse("linux/arm64"),
				Snapshotter: "native",
			},
		},
		{
			Name: "Match linux/arm and defaultSnapshotter",
			UnpackConfig: transfer.UnpackConfiguration{
				Platform:    platforms.MustParse("linux/arm"),
				Snapshotter: containerd.DefaultSnapshotter,
			},
			SupportedPlatforms: supportedPlatforms,
			Match:              true,
			ExpectedPlatform: transfer.UnpackConfiguration{
				Platform:    platforms.MustParse("linux/arm"),
				Snapshotter: containerd.DefaultSnapshotter,
			},
		},
		{
			Name: "Platform linux/amd64 input only match with defaultSnapshotter",
			UnpackConfig: transfer.UnpackConfiguration{
				Platform: platforms.MustParse("linux/amd64"),
			},
			SupportedPlatforms: supportedPlatforms,
			Match:              true,
			ExpectedPlatform: transfer.UnpackConfiguration{
				Platform:    platforms.MustParse("linux/amd64"),
				Snapshotter: containerd.DefaultSnapshotter,
			},
		},
		{
			Name: "Platform linux/arm64 input only match with defaultSnapshotter",
			UnpackConfig: transfer.UnpackConfiguration{
				Platform: platforms.MustParse("linux/arm64"),
			},
			SupportedPlatforms: supportedPlatforms,
			Match:              true,
			ExpectedPlatform: transfer.UnpackConfiguration{
				Platform:    platforms.MustParse("linux/arm64"),
				Snapshotter: containerd.DefaultSnapshotter,
			},
		},
		{
			Name: "Platform linux/arm input only match with defaultSnapshotter",
			UnpackConfig: transfer.UnpackConfiguration{
				Platform: platforms.MustParse("linux/arm"),
			},
			SupportedPlatforms: supportedPlatforms,
			Match:              true,
			ExpectedPlatform: transfer.UnpackConfiguration{
				Platform:    platforms.MustParse("linux/arm"),
				Snapshotter: containerd.DefaultSnapshotter,
			},
		},
	} {

		t.Run(testCase.Name, func(t *testing.T) {
			m, sp := getSupportedPlatform(testCase.UnpackConfig, testCase.SupportedPlatforms)
			if m == testCase.Match {
				if sp.SnapshotterKey == testCase.ExpectedPlatform.Snapshotter {
					if sp.Platform != nil && !sp.Platform.Match(testCase.ExpectedPlatform.Platform) {
						t.Fatalf("Expect Platform %v doesn't match", testCase.ExpectedPlatform.Platform)
					}
					if sp.Platform == nil && testCase.ExpectedPlatform.Platform.OS != "" {
						t.Fatalf("Expect Platform %v doesn't match", testCase.ExpectedPlatform.Platform)
					}
				} else {
					t.Fatalf("Expect SnapshotterKey %v, but got %v", testCase.ExpectedPlatform.Snapshotter, sp.SnapshotterKey)
				}
			} else {
				t.Fatalf("Expect match result %v, but got %v", testCase.Match, m)
			}
		})

	}

}
