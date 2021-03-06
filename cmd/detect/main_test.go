package main

import (
	"fmt"
	"path/filepath"
	"ruby-cnb/ruby"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/helper"

	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/google/go-cmp/cmp"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitDetect(t *testing.T) {
	spec.Run(t, "Detect", testDetect, spec.Report(report.Terminal{}))
}

func testDetect(t *testing.T, when spec.G, it spec.S) {
	var factory *test.DetectFactory

	it.Before(func() {
		RegisterTestingT(t)
		factory = test.NewDetectFactory(t)
	})

	it("always passes", func() {
		code, err := runDetect(factory.Detect)
		if err != nil {
			t.Error("Err in detect : ", err)
		}

		if diff := cmp.Diff(code, detect.PassStatusCode); diff != "" {
			t.Error("Problem : ", diff)
		}
	})

	when("testing versions", func() {
		when("there is no buildpack.yml", func() {
			it("shouldn't set the version in the buildplan", func() {
				code, err := runDetect(factory.Detect)
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(detect.PassStatusCode))

				Expect(factory.Plans.Plan).To(Equal(buildplan.Plan{
					Requires: []buildplan.Required{
						{
							Name:     ruby.Dependency,
							Version:  "",
							Metadata: buildplan.Metadata{"build": true, "launch": true},
						},
					},
					Provides: []buildplan.Provided{
						{ruby.Dependency},
					},
				}))
			})
		})

		when("there is a buildpack.yml", func() {
			const version string = "1.2.3"

			it.Before(func() {
				buildpackYAMLString := fmt.Sprintf("ruby:\n  version: %s", version)
				Expect(helper.WriteFile(filepath.Join(factory.Detect.Application.Root, "buildpack.yml"), 0666, buildpackYAMLString)).To(Succeed())
			})

			it("should pass with the requested version of ruby", func() {
				code, err := runDetect(factory.Detect)
				Expect(err).NotTo(HaveOccurred())
				Expect(code).To(Equal(detect.PassStatusCode))

				Expect(factory.Plans.Plan).To(Equal(buildplan.Plan{
					Requires: []buildplan.Required{
						{
							Name:     ruby.Dependency,
							Version:  version,
							Metadata: buildplan.Metadata{"build": true, "launch": true},
						},
					},
					Provides: []buildplan.Provided{
						{ruby.Dependency},
					},
				}))
			})
		})
	})
}
