package service

import (
	"fate-cloud-agent/pkg/db"
	"fate-cloud-agent/pkg/utils/logging"
	"os"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"helm.sh/helm/v3/pkg/chart"
)

func TestFateChart_save(t *testing.T) {
	type fields struct {
		version   string
		Chart     *chart.Chart
		HelmChart *db.HelmChart
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "test",
			fields:  fields{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FateChart{
				HelmChart: tt.fields.HelmChart,
			}
			if err := fc.save(); (err != nil) != tt.wantErr {
				t.Errorf("FateChart.save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetChartPath(t *testing.T) {
	InitConfigForTest()
	logging.InitLog()
	_ = os.Setenv("FATECLOUD_CHART_PATH", "./")
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				version: "v1.2.0",
			},
			want: viper.GetString("chart.path") + "fate/v1.2.0/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetChartPath(); got != tt.want {
				t.Errorf("GetChartPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFateChart(t *testing.T) {
	InitConfigForTest()
	logging.InitLog()
	_ = os.Setenv("FATECLOUD_CHART_PATH", "../../")
	type args struct {
		version string
	}
	var tests = []struct {
		name    string
		args    args
		want    *FateChart
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				version: "v1.2.0",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFateChart(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFateChart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFateChart() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFateChart_read(t *testing.T) {
	InitConfigForTest()
	type fields struct {
		HelmChart *db.HelmChart
	}
	type args struct {
		version string
	}
	var tests = []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			fields: fields{
				HelmChart: new(db.HelmChart),
			},
			args: args{
				version: "v1.2.0",
			},
			want:    "v1.2.0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FateChart{
				HelmChart: tt.fields.HelmChart,
			}
			got, err := fc.read(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("FateChart.read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Version != tt.want {
				t.Errorf("FateChart.read() Version = %v, want %v", got.Version, tt.want)
			}
		})
	}
}
