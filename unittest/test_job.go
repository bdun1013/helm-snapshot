package unittest

import (
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/bpdunni/helm-unittest/unittest/common"
	"github.com/bpdunni/helm-unittest/unittest/snapshot"
	"github.com/bpdunni/helm-unittest/unittest/validators"
	yaml "gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/engine"
	"helm.sh/helm/v3/pkg/getter"
)

type orderedSnapshotComparer struct {
	cache   *snapshot.Cache
	test    string
	counter uint
}

func (s *orderedSnapshotComparer) CompareToSnapshot(content interface{}) *snapshot.CompareResult {
	s.counter++
	return s.cache.Compare(s.test, s.counter, content)
}

// TestJob definition of a test, including values and assertions
type TestJob struct {
	Name         string `yaml:"it"`
	Values       []string
	Set          map[string]interface{}
	Assertions   []*Assertion `yaml:"asserts"`
	Release      chartutil.ReleaseOptions
	Capabilities chartutil.Capabilities
	// route indicate which chart in the dependency hierarchy
	// like "parant-chart", "parent-charts/charts/child-chart"
	chartRoute string
	// where the test suite file located
	definitionFile string
	// template assertion should assert if not specified
	defaultTemplateToAssert string
}

// Run render the chart and validate it with assertions in TestJob
func (t *TestJob) Run(
	targetChart *chart.Chart,
	cache *snapshot.Cache,
	result *TestJobResult,
) *TestJobResult {
	t.polishAssertionsTemplate(targetChart)
	result.DisplayName = t.Name

	userValues, err := t.getUserValues()
	if err != nil {
		result.ExecError = err
		return result
	}

	outputOfFiles, err := t.renderChart(targetChart, userValues)
	if err != nil {
		result.ExecError = err
		return result
	}

	manifestsOfFiles, err := t.parseManifestsFromOutputOfFiles(outputOfFiles)
	if err != nil {
		result.ExecError = err
		return result
	}

	snapshotComparer := &orderedSnapshotComparer{cache: cache, test: t.Name}
	result.Passed, result.AssertsResult = t.runAssertions(
		manifestsOfFiles,
		snapshotComparer,
	)

	return result
}

// liberally borrows from helm-template
func (t *TestJob) getUserValues() (*values.Options, error) {

	valueOpts := &values.Options{
		ValueFiles: make([]string, 0),
	}

	for _, specifiedPath := range t.Values {
		var valueFilePath string
		if path.IsAbs(specifiedPath) {
			valueFilePath = specifiedPath
		} else {
			valueFilePath = filepath.Join(filepath.Dir(t.definitionFile), specifiedPath)
		}

		valueOpts.ValueFiles = append(valueOpts.ValueFiles, valueFilePath)
	}

	return valueOpts, nil
}

// render the chart and return result map
func (t *TestJob) renderChart(targetChart *chart.Chart, userValues *values.Options) (map[string]string, error) {
	options := *t.releaseOption()
	caps := *t.capabilityOption()

	settings := cli.New()

	p := getter.All(settings)
	vals, err := userValues.MergeValues(p)

	if err := chartutil.ProcessDependencies(targetChart, vals); err != nil {
		return nil, err
	}

	renderedVals, err := chartutil.ToRenderValues(targetChart, vals, options, &caps)
	if err != nil {
		return nil, err
	}

	outputOfFiles, err := engine.Render(targetChart, renderedVals)
	if err != nil {
		return nil, err
	}

	return outputOfFiles, nil
}

// get chartutil.ReleaseOptions ready for render
func (t *TestJob) releaseOption() *chartutil.ReleaseOptions {
	options := chartutil.ReleaseOptions{
		Name:      "RELEASE-NAME",
		Namespace: "NAMESPACE",
		Revision:  t.Release.Revision,
		IsInstall: !t.Release.IsUpgrade,
		IsUpgrade: t.Release.IsUpgrade,
	}
	if t.Release.Name != "" {
		options.Name = t.Release.Name
	}
	if t.Release.Namespace != "" {
		options.Namespace = t.Release.Namespace
	}
	return &options
}

// get chartutil.CapabilityOptions ready for render
// Only supports APIVersions for now
func (t *TestJob) capabilityOption() *chartutil.Capabilities {
	return chartutil.DefaultCapabilities
}

// parse rendered manifest if it's yaml
func (t *TestJob) parseManifestsFromOutputOfFiles(outputOfFiles map[string]string) (
	map[string][]common.K8sManifest,
	error,
) {
	manifestsOfFiles := make(map[string][]common.K8sManifest)

	for file, rendered := range outputOfFiles {
		decoder := yaml.NewDecoder(strings.NewReader(rendered))

		if filepath.Ext(file) == ".yaml" {
			manifests := make([]common.K8sManifest, 0)

			for {
				manifest := make(common.K8sManifest)
				if err := decoder.Decode(manifest); err != nil {
					if err == io.EOF {
						break
					} else {
						return nil, err
					}
				}

				if len(manifest) > 0 {
					manifests = append(manifests, manifest)
				}
			}

			manifestsOfFiles[file] = manifests
		}
	}

	return manifestsOfFiles, nil
}

// run Assert of all assertions of test
func (t *TestJob) runAssertions(
	manifestsOfFiles map[string][]common.K8sManifest,
	snapshotComparer validators.SnapshotComparer,
) (bool, []*AssertionResult) {
	testPass := true
	assertsResult := make([]*AssertionResult, len(t.Assertions))

	for idx, assertion := range t.Assertions {
		result := assertion.Assert(
			manifestsOfFiles,
			snapshotComparer,
			&AssertionResult{Index: idx},
		)

		assertsResult[idx] = result
		testPass = testPass && result.Passed
	}
	return testPass, assertsResult
}

// add prefix to Assertion.Template
func (t *TestJob) polishAssertionsTemplate(targetChart *chart.Chart) {
	if t.chartRoute == "" {
		t.chartRoute = targetChart.Metadata.Name
	}

	for _, assertion := range t.Assertions {
		var templateToAssert string

		if assertion.Template == "" {
			if t.defaultTemplateToAssert == "" {
				return
			}
			templateToAssert = t.defaultTemplateToAssert
		} else {
			templateToAssert = assertion.Template
		}

		// map the file name to the path of helm rendered result
		assertion.Template = filepath.ToSlash(
			filepath.Join(t.chartRoute, "templates", templateToAssert),
		)
	}
}
