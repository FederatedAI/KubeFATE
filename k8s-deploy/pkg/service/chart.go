package service

import (
	"context"
	"helm.sh/helm/v3/pkg/chart/loader"
	"sync"

	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"

	"encoding/json"
	"encoding/xml"

	"fate-cloud-agent/pkg/db"
	"github.com/spf13/viper"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
	"sigs.k8s.io/yaml"
)

const latestChartVersion string = "v1.2.0"

type Chart interface {
	save(Chart) error
	read(version string) (Chart, error)
	load(version string) (Chart, error)
}

type FateChart struct {
	*db.HelmChart
}

func (fc *FateChart) save() error {
	helmUUID, err := db.Save(fc.HelmChart)
	if err != nil {
		return err
	}
	log.Debug().Str("helmUUID", helmUUID).Msg("helm chart save uuid")
	return nil
}

func (fc *FateChart) read(version string) (*FateChart, error) {

	helmChart, err := db.FindHelmByVersion(version)
	if err != nil {
		return nil, err
	}
	if helmChart == nil {
		return nil, errors.New("find chart error")
	}

	log.Debug().Interface("helmChart version", helmChart.Version).Msg("find chart from db success")

	return &FateChart{
		HelmChart: helmChart,
	}, nil
}
func (fc *FateChart) load(version string) (*FateChart, error) {
	chartPath := GetChartPath()
	settings := cli.New()

	cfg := new(action.Configuration)
	client := action.NewInstall(cfg)

	client.ChartPathOptions.Version = version

	cp, err := client.ChartPathOptions.LocateChart(chartPath, settings)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("FateChart chartPath:", cp).Msg("chartPath:")

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	helmChart, err := ChartRequestedTohelmChart(chartRequested)
	if err != nil {
		return nil, err
	}
	//helmChart, err := SaveChartFromPath(cp, "fate")
	//if err != nil {
	//	return nil, err
	//}

	return &FateChart{
		HelmChart: helmChart,
	}, nil
}

func ChartRequestedTohelmChart(chartRequested *chart.Chart) (*db.HelmChart, error) {
	if chartRequested == nil || chartRequested.Raw == nil {
		log.Error().Msg("chartRequested not exist")
		return nil, errors.New("chartRequested not exist")
	}

	var chartData string
	var valuesData string
	var ValuesTemplate string
	for _, v := range chartRequested.Raw {
		if v.Name == "Chart.yaml" {
			chartData = string(v.Data)
		}
		if v.Name == "values.yaml" {
			valuesData = string(v.Data)
		}
		if v.Name == "values-template.yaml" {
			ValuesTemplate = string(v.Data)
		}
	}

	helmChart := db.NewHelmChart(chartRequested.Name(),
		chartData, valuesData, chartRequested.Templates, chartRequested.Metadata.Version, chartRequested.AppVersion())

	helmChart.ValuesTemplate = ValuesTemplate
	return helmChart, nil
}

func (fc *FateChart) GetChartValuesTemplates() (string, error) {
	if fc.ValuesTemplate == "" {
		return "", errors.New("FateChart ValuesTemplate not exist")
	}
	return fc.ValuesTemplate, nil
}

func (fc *FateChart) GetChartValues(v map[string]interface{}) (map[string]interface{}, error) {
	// template to values
	template, err := fc.GetChartValuesTemplates()
	if err != nil {
		log.Err(err).Msg("GetChartValuesTemplates error")
		return nil, err
	}
	values, err := MapToConfig(v, template)
	if err != nil {
		log.Err(err).Interface("v", v).Interface("template", template).Msg("MapToConfig error")
		return nil, err
	}
	// values to map
	vals := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(values), &vals)
	if err != nil {
		log.Err(err).Msg("values yaml Unmarshal error")
		return nil, err
	}
	return vals, nil
}

// todo  get chart by version from repository
func GetFateChart(version string) (*FateChart, error) {
	chartVersion := version
	//if version == "" {
	//	chartVersion = latestChartVersion
	//}
	fc := new(FateChart)
	fc, err := fc.read(chartVersion)
	if err == nil {
		return fc, nil
	}
	log.Debug().Interface("read error", err).Msg("read version FateChart err")

	fc, err = fc.load(chartVersion)
	if err != nil {
		return nil, err
	}
	err = fc.save()
	if err != nil {
		return nil, err
	}
	return fc, nil
}

func (fc *FateChart) ToHelmChart() (*chart.Chart, error) {
	if fc == nil || fc.HelmChart == nil {
		return nil, errors.New("FateChart not exist")
	}
	return ConvertToChart(fc.HelmChart)
}

func GetChartPath() string {
	ChartPath := viper.GetString("repo.name") + "/fate"
	log.Debug().Str("ChartPath", ChartPath).Msg("ChartPath")
	return ChartPath
}

type Value struct {
	Val []byte
	T   string // type json yaml yml
}

func (v *Value) Unmarshal() (map[string]interface{}, error) {
	si := make(map[string]interface{})
	switch v.T {
	case "yaml":
		err := yaml.Unmarshal(v.Val, &si)
		return si, err
	case "json":
		err := json.Unmarshal(v.Val, &si)
		return si, err
	case "xml":
		err := xml.Unmarshal(v.Val, &si)
		return si, err
	}
	return nil, errors.New("unrecognized type")
}

type repoAddOptions struct {
	name     string
	url      string
	username string
	password string
	noUpdate bool

	certFile string
	keyFile  string
	caFile   string

	repoFile  string
	repoCache string
}

func (o *repoAddOptions) run(settings *cli.EnvSettings) error {
	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(o.repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Error().Err(err).Msg("MkdirAll")
		return err
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(o.repoFile, filepath.Ext(o.repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		log.Error().Err(err).Msg("TryLockContext")
		return err
	}

	b, err := ioutil.ReadFile(o.repoFile)
	if err != nil && !os.IsNotExist(err) {
		log.Error().Err(err).Msg("ReadFile")
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		log.Error().Err(err).Msg("Unmarshal")
		return err
	}

	if o.noUpdate && f.Has(o.name) {
		return errors.Errorf("repository name (%s) already exists, please specify a different name", o.name)
	}

	c := repo.Entry{
		Name:     o.name,
		URL:      o.url,
		Username: o.username,
		Password: o.password,
		CertFile: o.certFile,
		KeyFile:  o.keyFile,
		CAFile:   o.caFile,
	}

	r, err := repo.NewChartRepository(&c, getter.All(settings))
	if err != nil {
		log.Error().Err(err).Msg("ReadFile")
		return err
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		return errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", o.url)
	}

	f.Update(&c)

	if err := f.WriteFile(o.repoFile, 0644); err != nil {
		log.Error().Err(err).Msg("WriteFile")
		return err
	}
	log.Debug().Msgf("%q has been added to your repositories\n", o.name)
	return nil
}

type repoUpdateOptions struct {
	update   func([]*repo.ChartRepository)
	repoFile string
}

func (o *repoUpdateOptions) run(settings *cli.EnvSettings) error {
	f, err := repo.LoadFile(o.repoFile)
	if isNotExist(err) || len(f.Repositories) == 0 {
		return errors.New("no repositories found. You must add one before updating")
	}
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(settings))
		if err != nil {
			return err
		}
		repos = append(repos, r)
	}

	o.update(repos)
	return nil
}
func isNotExist(err error) bool {
	return os.IsNotExist(errors.Cause(err))
}
func updateCharts(repos []*repo.ChartRepository) {
	log.Debug().Msg("Hang tight while we grab the latest from your chart repositories...")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				log.Debug().Msgf("...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
			} else {
				log.Debug().Msgf("...Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()
	log.Debug().Msg("Update Complete.")
}
func RepoAddAndUpdate() error {
	settings := cli.New()
	o := new(repoAddOptions)

	o.name = viper.GetString("repo.name")
	o.url = viper.GetString("repo.url")
	o.username = viper.GetString("repo.username")
	o.password = viper.GetString("repo.password")
	o.repoFile = settings.RepositoryConfig
	o.repoCache = settings.RepositoryCache
	err := o.run(settings)
	if err != nil {
		log.Error().Err(err).Msg("repoAdd")
		return err
	}
	log.Debug().Msg("repoAdd success")
	ou := &repoUpdateOptions{update: updateCharts}
	ou.repoFile = settings.RepositoryConfig
	err = ou.run(settings)
	return err
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}
