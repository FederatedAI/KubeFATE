package modules

import (
	"os"
	"testing"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/orm"
	"github.com/jinzhu/gorm"
	"sigs.k8s.io/yaml"
)

func TestCluster_Install(t *testing.T) {
	InitConfigForTest()
	mysql := new(orm.Mysql)
	mysql.Setup()
	DB = orm.DBCLIENT
	//DB.LogMode(true)

	// Create Table
	e := &Cluster{}
	e.InitTable()

	// Drop Table
	defer e.DropTable()
	hc := &HelmChart{}
	hc.InitTable()
	//defer hc.DropTable()

	var spec MapStringInterface
	err := yaml.Unmarshal([]byte(cluster), &spec)
	if err != nil {
		t.Errorf("yaml.Unmarshal() error = %v", err)
	}
	type fields struct {
		Uuid         string
		Name         string
		NameSpace    string
		ChartName    string
		ChartVersion string
		Values       string
		Spec         MapStringInterface
		Revision     int8
		HelmRevision int8
		ChartValues  MapStringInterface
		Status       ClusterStatus
		Info         MapStringInterface
		Model        gorm.Model
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "fate-9999",
			fields: fields{
				Uuid:         "",
				Name:         "fate-9999",
				NameSpace:    "fate-9999",
				ChartName:    "fate",
				ChartVersion: "v1.4.0",
				Values:       cluster,
				Spec:         spec,
				Revision:     0,
				HelmRevision: 0,
				ChartValues:  nil,
				Status:       0,
				Info:         nil,
				Model:        gorm.Model{},
			},
			wantErr: false,
		},
		{
			name: "namespace not found",
			fields: fields{
				Uuid:         "",
				Name:         "fate-not",
				NameSpace:    "fate-not",
				ChartName:    "fate",
				ChartVersion: "v1.4.0",
				Values:       cluster,
				Spec:         spec,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Cluster{
				Uuid:         tt.fields.Uuid,
				Name:         tt.fields.Name,
				NameSpace:    tt.fields.NameSpace,
				ChartName:    tt.fields.ChartName,
				ChartVersion: tt.fields.ChartVersion,
				Values:       tt.fields.Values,
				Spec:         tt.fields.Spec,
				Revision:     tt.fields.Revision,
				HelmRevision: tt.fields.HelmRevision,
				ChartValues:  tt.fields.ChartValues,
				Status:       tt.fields.Status,
				Info:         tt.fields.Info,
				Model:        tt.fields.Model,
			}
			if err := e.HelmInstall(); (err != nil) != tt.wantErr {
				t.Errorf("Cluster.Install() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var cluster = `name: fate-9999
namespace: fate-9999
chartName: fate
chartVersion: v1.4.0
partyId: 9999
registry: ""
pullPolicy: 
persistence: false
modules:
  - rollsite
  - clustermanager
  - nodemanager
  - mysql
  - python
# - client

rollsite: 
  type: NodePort
  nodePort: 30009
  exchange:
    ip: 192.168.1.1
    port: 30000
  partyList:
  - partyId: 10000
    partyIp: 192.168.10.1
    partyPort: 30010
  nodeSelector: {}

nodemanager:
  count: 3
  sessionProcessorsPerNode: 4
  list:
  - name: nodemanager
    nodeSelector: {}
    sessionProcessorsPerNode: 2
    subPath: "nodemanager"
    existingClaim: ""
    storageClass: "nodemanager"
    accessMode: ReadWriteOnce
    size: 1Gi

python:
  fateflowType: NodePort
  fateflowNodePort: 30109
  nodeSelector: {}

mysql: 
  nodeSelector: {}
  ip: mysql
  port: 3306
  database: eggroll_meta
  user: fate
  password: fate_dev
  subPath: ""
  existingClaim: ""
  storageClass: "mysql"
  accessMode: ReadWriteOnce
  size: 1Gi

# If use external MySQL, uncomment and change this section
# externalMysqlIp: mysql
# externalMysqlPort: 3306
# externalMysqlDatabase: eggroll_meta
# externalMysqlUser: fate
# externalMysqlPassword: fate_dev

# If FATE-Serving deployed, uncomment and change
# servingIp: 192.168.9.1
# servingPort: 30209`

var updateCluster = `name: fate-9999
namespace: fate-9999
chartName: fate
chartVersion: v1.4.0
partyId: 9999
registry: ""
pullPolicy: 
persistence: false
istio:
  enabled: false
modules:
  - rollsite
  - clustermanager
  - nodemanager
  - mysql
  - python
  - client

rollsite: 
  type: NodePort
  nodePort: 30009
  exchange:
    ip: 192.168.1.1
    port: 30000
  partyList:
  - partyId: 10000
    partyIp: 192.168.10.1
    partyPort: 30010
  nodeSelector: {}

nodemanager:
  count: 1
  sessionProcessorsPerNode: 2
  list:
  - name: nodemanager
    nodeSelector: {}
    sessionProcessorsPerNode: 4
    subPath: "nodemanager"
    existingClaim: ""
    storageClass: "nodemanager"
    accessMode: ReadWriteOnce
    size: 1Gi

python:
  fateflowType: NodePort
  fateflowNodePort: 30109
  nodeSelector: {}

mysql: 
  nodeSelector: {}
  ip: mysql
  port: 3306
  database: eggroll_meta
  user: fate
  password: fate_dev
  subPath: ""
  existingClaim: ""
  storageClass: "mysql"
  accessMode: ReadWriteOnce
  size: 1Gi

# If use external MySQL, uncomment and change this section
# externalMysqlIp: mysql
# externalMysqlPort: 3306
# externalMysqlDatabase: eggroll_meta
# externalMysqlUser: fate
# externalMysqlPassword: fate_dev

# If FATE-Serving deployed, uncomment and change
# servingIp: 192.168.9.1
# servingPort: 30209`

var upgradeCluster = `name: fate-9999
namespace: fate-9999
chartName: fate
chartVersion: v1.4.0-a
partyId: 9999
registry: ""
pullPolicy: 
persistence: false
istio:
  enabled: false
modules:
  - rollsite
  - clustermanager
  - nodemanager
  - mysql
  - python
  - client

rollsite: 
  type: NodePort
  nodePort: 30009
  exchange:
    ip: 192.168.1.1
    port: 30000
  partyList:
  - partyId: 10000
    partyIp: 192.168.10.1
    partyPort: 30010
  nodeSelector: {}

nodemanager:
  count: 1
  sessionProcessorsPerNode: 2
  list:
  - name: nodemanager
    nodeSelector: {}
    sessionProcessorsPerNode: 4
    subPath: "nodemanager"
    existingClaim: ""
    storageClass: "nodemanager"
    accessMode: ReadWriteOnce
    size: 1Gi

python:
  fateflowType: NodePort
  fateflowNodePort: 30109
  nodeSelector: {}

mysql: 
  nodeSelector: {}
  ip: mysql
  port: 3306
  database: eggroll_meta
  user: fate
  password: fate_dev
  subPath: ""
  existingClaim: ""
  storageClass: "mysql"
  accessMode: ReadWriteOnce
  size: 1Gi

# If use external MySQL, uncomment and change this section
# externalMysqlIp: mysql
# externalMysqlPort: 3306
# externalMysqlDatabase: eggroll_meta
# externalMysqlUser: fate
# externalMysqlPassword: fate_dev

# If FATE-Serving deployed, uncomment and change
# servingIp: 192.168.9.1
# servingPort: 30209`

var upgradeCluster_1 = `name: fate-9999
namespace: fate-9999
chartName: fate
chartVersion: v1.4.3
partyId: 9999
registry: ""
pullPolicy: 
persistence: false
modules:
  - rollsite
  - clustermanager
  - nodemanager
  - mysql
  - python
  - client

rollsite: 
  type: NodePort
  nodePort: 30009
  exchange:
    ip: 192.168.1.1
    port: 30000
  partyList:
  - partyId: 10000
    partyIp: 192.168.10.1
    partyPort: 30010
  nodeSelector: {}

nodemanager:
  count: 1
  sessionProcessorsPerNode: 2
  list:
  - name: nodemanager
    nodeSelector: {}
    sessionProcessorsPerNode: 4
    subPath: "nodemanager"
    existingClaim: ""
    storageClass: "nodemanager"
    accessMode: ReadWriteOnce
    size: 1Gi

python:
  fateflowType: NodePort
  fateflowNodePort: 30109
  nodeSelector: {}

mysql: 
  nodeSelector: {}
  ip: mysql
  port: 3306
  database: eggroll_meta
  user: fate
  password: fate_dev
  subPath: ""
  existingClaim: ""
  storageClass: "mysql"
  accessMode: ReadWriteOnce
  size: 1Gi
`

func TestCluster_HelmDelete(t *testing.T) {
	type fields struct {
		Uuid         string
		Name         string
		NameSpace    string
		ChartName    string
		ChartVersion string
		Values       string
		Spec         MapStringInterface
		Revision     int8
		HelmRevision int8
		ChartValues  MapStringInterface
		Status       ClusterStatus
		Info         MapStringInterface
		Model        gorm.Model
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "fate-9999",
			fields: fields{
				Name:      "fate-9999",
				NameSpace: "fate-9999",
			},
			wantErr: false,
		},
		{
			name: "release: not found",
			fields: fields{
				Name:      "fate-not",
				NameSpace: "fate-not",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Cluster{
				Uuid:         tt.fields.Uuid,
				Name:         tt.fields.Name,
				NameSpace:    tt.fields.NameSpace,
				ChartName:    tt.fields.ChartName,
				ChartVersion: tt.fields.ChartVersion,
				Values:       tt.fields.Values,
				Spec:         tt.fields.Spec,
				Revision:     tt.fields.Revision,
				HelmRevision: tt.fields.HelmRevision,
				ChartValues:  tt.fields.ChartValues,
				Status:       tt.fields.Status,
				Info:         tt.fields.Info,
				Model:        tt.fields.Model,
			}
			if err := e.HelmDelete(); (err != nil) != tt.wantErr {
				t.Errorf("Cluster.HelmDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCluster_Helm(t *testing.T) {
	InitConfigForTest()
	mysql := new(orm.Mysql)
	mysql.Setup()
	DB = orm.DBCLIENT
	//DB.LogMode(true)

	// Create Table
	e := &Cluster{}
	e.InitTable()

	// Drop Table
	defer e.DropTable()
	hc := &HelmChart{}
	hc.InitTable()
	//defer hc.DropTable()

	var spec MapStringInterface
	err := yaml.Unmarshal([]byte(cluster), &spec)
	if err != nil {
		t.Errorf("yaml.Unmarshal() error = %v", err)
	}

	e = &Cluster{
		Name:         "fate-9999",
		NameSpace:    "fate-9999",
		ChartName:    "fate",
		ChartVersion: "v1.4.0",
		Values:       cluster,
		Spec:         spec,
	}
	t.Run("HelmInstall", func(t *testing.T) {
		t.Log("Normal install")
		if err := e.HelmInstall(); (err != nil) != false {
			t.Errorf("Cluster.Install() error = %v, wantErr %v", err, false)
		}
		t.Log("Repeat install")
		if err := e.HelmInstall(); (err != nil) != true {
			t.Errorf("Cluster.Install() error = %v, wantErr %v", err, true)
		}
	})

	var upgradeSpec MapStringInterface
	err = yaml.Unmarshal([]byte(updateCluster), &upgradeSpec)
	if err != nil {
		t.Errorf("yaml.Unmarshal() error = %v", err)
	}

	upgrade := &Cluster{
		Name:         "fate-9999",
		NameSpace:    "fate-9999",
		ChartName:    "fate",
		ChartVersion: "v1.4.0-a",
		Values:       upgradeCluster,
		Spec:         upgradeSpec,
	}
	t.Run("HelmUpgrade", func(t *testing.T) {
		t.Log("Repeat install")
		if err := upgrade.HelmUpgrade(); (err != nil) != false {
			t.Errorf("Cluster.HelmUpgrade() error = %v, wantErr %v", err, false)
		}
	})

	var upgradeSpec_1 MapStringInterface
	err = yaml.Unmarshal([]byte(upgradeCluster_1), &upgradeSpec_1)
	if err != nil {
		t.Errorf("yaml.Unmarshal() error = %v", err)
	}

	upgrade_1 := &Cluster{
		Name:         "fate-9999",
		NameSpace:    "fate-9999",
		ChartName:    "fate",
		ChartVersion: "v1.4.3",
		Values:       upgradeCluster_1,
		Spec:         upgradeSpec_1,
	}
	t.Run("HelmUpgrade", func(t *testing.T) {
		t.Log("Repeat install")
		if err := upgrade_1.HelmUpgrade(); (err != nil) != false {
			t.Errorf("Cluster.HelmUpgrade() error = %v, wantErr %v", err, false)
		}
	})
	t.Run("HelmDelete", func(t *testing.T) {
		if err := e.HelmDelete(); (err != nil) != false {
			t.Errorf("Cluster.HelmDelete() error = %v, wantErr %v", err, false)
		}
	})

}

func TestCluster_HelmUpgrade(t *testing.T) {
	InitConfigForTest()
	mysql := new(orm.Mysql)
	mysql.Setup()
	DB = orm.DBCLIENT
	//DB.LogMode(true)

	// Create Table
	e := &Cluster{}
	e.InitTable()

	// Drop Table
	defer e.DropTable()
	hc := &HelmChart{}
	hc.InitTable()
	//defer hc.DropTable()
	var upgradeSpec_1 MapStringInterface
	err := yaml.Unmarshal([]byte(upgradeCluster_1), &upgradeSpec_1)
	if err != nil {
		t.Errorf("yaml.Unmarshal() error = %v", err)
	}
	type fields struct {
		Uuid         string
		Name         string
		NameSpace    string
		ChartName    string
		ChartVersion string
		Values       string
		Spec         MapStringInterface
		Revision     int8
		HelmRevision int8
		ChartValues  MapStringInterface
		Status       ClusterStatus
		Info         MapStringInterface
		Model        gorm.Model
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Different versions of chart upgrade fate",
			fields: fields{
				Name:         "fate-9999",
				NameSpace:    "fate-9999",
				ChartName:    "fate",
				ChartVersion: "v1.4.0-a",
				Values:       upgradeCluster_1,
				Spec:         upgradeSpec_1,
			},
			wantErr: true,
		},
		{
			name: "Different versions of chart upgrade fate",
			fields: fields{
				Name:         "fate-9999",
				NameSpace:    "fate-9999",
				ChartName:    "fate",
				ChartVersion: "v1.4.3",
				Values:       upgradeCluster_1,
				Spec:         upgradeSpec_1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Cluster{
				Uuid:         tt.fields.Uuid,
				Name:         tt.fields.Name,
				NameSpace:    tt.fields.NameSpace,
				ChartName:    tt.fields.ChartName,
				ChartVersion: tt.fields.ChartVersion,
				Values:       tt.fields.Values,
				Spec:         tt.fields.Spec,
				Revision:     tt.fields.Revision,
				HelmRevision: tt.fields.HelmRevision,
				ChartValues:  tt.fields.ChartValues,
				Status:       tt.fields.Status,
				Info:         tt.fields.Info,
				Model:        tt.fields.Model,
			}
			if err := e.HelmUpgrade(); (err != nil) != tt.wantErr {
				t.Errorf("Cluster.HelmUpgrade() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCluster_HelmInstall(t *testing.T) {
	InitConfigForTest()
	mysql := new(orm.Mysql)
	mysql.Setup()
	DB = orm.DBCLIENT
	//DB.LogMode(true)

	// Create Table
	e := &Cluster{}
	e.InitTable()

	// Drop Table
	defer e.DropTable()
	hc := &HelmChart{}
	hc.InitTable()
	//defer hc.DropTable()
	var spec MapStringInterface
	err := yaml.Unmarshal([]byte(cluster), &spec)
	if err != nil {
		t.Errorf("yaml.Unmarshal() error = %v", err)
	}
	type fields struct {
		Uuid         string
		Name         string
		NameSpace    string
		ChartName    string
		ChartVersion string
		Values       string
		Spec         MapStringInterface
		Revision     int8
		HelmRevision int8
		ChartValues  MapStringInterface
		Status       ClusterStatus
		Info         MapStringInterface
		Model        gorm.Model
	}
	type args struct {
		dryRun bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			fields: fields{
				Name:         "fate-9999",
				NameSpace:    "fate-9999",
				ChartName:    "fate",
				ChartVersion: "v1.4.0",
				Values:       cluster,
				Spec:         spec,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Cluster{
				Uuid:         tt.fields.Uuid,
				Name:         tt.fields.Name,
				NameSpace:    tt.fields.NameSpace,
				ChartName:    tt.fields.ChartName,
				ChartVersion: tt.fields.ChartVersion,
				Values:       tt.fields.Values,
				Spec:         tt.fields.Spec,
				Revision:     tt.fields.Revision,
				HelmRevision: tt.fields.HelmRevision,
				ChartValues:  tt.fields.ChartValues,
				Status:       tt.fields.Status,
				Info:         tt.fields.Info,
				Model:        tt.fields.Model,
			}
			if err := e.HelmInstall(); (err != nil) != tt.wantErr {
				t.Errorf("Cluster.HelmInstall() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCluster_HelmInstall_Chart(t *testing.T) {
	os.Setenv("FATECLOUD_REPO_URL", "")
	InitConfigForTest()
	mysql := new(orm.Mysql)
	mysql.Setup()
	DB = orm.DBCLIENT
	//DB.LogMode(true)

	// Create Table
	e := &Cluster{}
	e.InitTable()

	// Drop Table
	defer e.DropTable()
	hc := &HelmChart{}
	hc.InitTable()
	//defer hc.DropTable()
	var spec MapStringInterface
	err := yaml.Unmarshal([]byte(cluster), &spec)
	if err != nil {
		t.Errorf("yaml.Unmarshal() error = %v", err)
	}
	type fields struct {
		Uuid         string
		Name         string
		NameSpace    string
		ChartName    string
		ChartVersion string
		Values       string
		Spec         MapStringInterface
		Revision     int8
		HelmRevision int8
		ChartValues  MapStringInterface
		Status       ClusterStatus
		Info         MapStringInterface
		Model        gorm.Model
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "chart does not exist",
			fields: fields{
				Name:         "fate-9999",
				NameSpace:    "fate-9999",
				ChartName:    "fate",
				ChartVersion: "v1.4.3",
				Values:       cluster,
				Spec:         spec,
			},
			wantErr: true,
		},
		{
			name: "chart exist",
			fields: fields{
				Name:         "fate-9999",
				NameSpace:    "fate-9999",
				ChartName:    "fate",
				ChartVersion: "v1.4.0",
				Values:       cluster,
				Spec:         spec,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Cluster{
				Uuid:         tt.fields.Uuid,
				Name:         tt.fields.Name,
				NameSpace:    tt.fields.NameSpace,
				ChartName:    tt.fields.ChartName,
				ChartVersion: tt.fields.ChartVersion,
				Values:       tt.fields.Values,
				Spec:         tt.fields.Spec,
				Revision:     tt.fields.Revision,
				HelmRevision: tt.fields.HelmRevision,
				ChartValues:  tt.fields.ChartValues,
				Status:       tt.fields.Status,
				Info:         tt.fields.Info,
				Model:        tt.fields.Model,
			}
			if err := e.HelmInstall(); (err != nil) != tt.wantErr {
				t.Errorf("Cluster.HelmInstall() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
