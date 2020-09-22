/*
 * Copyright 2019-2020 VMware, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package cli

import (
	"encoding/json"
	"testing"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
)

func TestCluster_outPutInfo(t *testing.T) {
	type fields struct {
		all bool
	}
	type args struct {
		result interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:   "cluster info",
			fields: fields{},
			args: args{
				result: func() *ClusterResult {
					cluster := modules.Cluster{}
					_ = json.Unmarshal([]byte(result), &cluster)
					return &ClusterResult{Data: &cluster, Msg: ""}
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				all: tt.fields.all,
			}
			if err := c.outPutInfo(tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("Cluster.outPutInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var result = `{
        "uuid": "7bd6ddf2-9a65-40ba-8c37-004fc12d32ef",
        "name": "fate-9999",
        "namespaces": "fate-9999",
        "chart_name": "fate",
        "chart_version": "v1.4.4",
        "values": "{\"chartName\":\"fate\",\"chartVersion\":\"v1.4.4\",\"modules\":[\"rollsite\",\"clustermanager\",\"nodemanager\",\"mysql\",\"python\",\"client\"],\"mysql\":{\"accessMode\":\"ReadWriteOnce\",\"database\":\"eggroll_meta\",\"existingClaim\":\"\",\"ip\":\"mysql\",\"nodeSelector\":{},\"password\":\"fate_dev\",\"port\":3306,\"size\":\"1Gi\",\"storageClass\":\"mysql\",\"subPath\":\"\",\"user\":\"fate\"},\"name\":\"fate-9999\",\"namespace\":\"fate-9999\",\"nodemanager\":{\"count\":0,\"list\":[{\"accessMode\":\"ReadWriteOnce\",\"existingClaim\":\"\",\"name\":\"nodemanager\",\"nodeSelector\":{},\"sessionProcessorsPerNode\":4,\"size\":\"1Gi\",\"storageClass\":\"nodemanager\",\"subPath\":\"nodemanager\"}],\"sessionProcessorsPerNode\":4},\"partyId\":9999,\"persistence\":false,\"pullPolicy\":null,\"python\":{\"fateflowNodePort\":30109,\"fateflowType\":\"NodePort\",\"nodeSelector\":{}},\"registry\":\"\",\"rollsite\":{\"exchange\":{\"ip\":\"192.168.1.1\",\"port\":30000},\"nodePort\":30009,\"nodeSelector\":{},\"partyList\":[{\"partyId\":10000,\"partyIp\":\"10.186.212.217\",\"partyPort\":30010}],\"type\":\"NodePort\"},\"servingIp\":\"10.186.216.97\",\"servingPort\":30209}",
        "Spec": {
            "chartName": "fate",
            "chartVersion": "v1.4.4",
            "modules": [
                "rollsite",
                "clustermanager",
                "nodemanager",
                "mysql",
                "python",
                "client"
            ],
            "mysql": {
                "accessMode": "ReadWriteOnce",
                "database": "eggroll_meta",
                "existingClaim": "",
                "ip": "mysql",
                "nodeSelector": {},
                "password": "fate_dev",
                "port": 3306,
                "size": "1Gi",
                "storageClass": "mysql",
                "subPath": "",
                "user": "fate"
            },
            "name": "fate-9999",
            "namespace": "fate-9999",
            "nodemanager": {
                "count": 0,
                "list": [
                    {
                        "accessMode": "ReadWriteOnce",
                        "existingClaim": "",
                        "name": "nodemanager",
                        "nodeSelector": {},
                        "sessionProcessorsPerNode": 4,
                        "size": "1Gi",
                        "storageClass": "nodemanager",
                        "subPath": "nodemanager"
                    }
                ],
                "sessionProcessorsPerNode": 4
            },
            "partyId": 9999,
            "persistence": false,
            "pullPolicy": null,
            "python": {
                "fateflowNodePort": 30109,
                "fateflowType": "NodePort",
                "nodeSelector": {}
            },
            "registry": "",
            "rollsite": {
                "exchange": {
                    "ip": "192.168.1.1",
                    "port": 30000
                },
                "nodePort": 30009,
                "nodeSelector": {},
                "partyList": [
                    {
                        "partyId": 10000,
                        "partyIp": "10.186.212.217",
                        "partyPort": 30010
                    }
                ],
                "type": "NodePort"
            },
            "servingIp": "10.186.216.97",
            "servingPort": 30209
        },
        "revision": 1,
        "helm_revision": 0,
        "chart_values": null,
        "status": "Running",
        "Info": {
            "dashboard": [
                "9999.notebook.kubefate.net",
                "9999.fateboard.kubefate.net"
            ],
            "ip": "10.186.216.97",
            "modules": [
                "clustermanager-687b87d776-mlk28",
                "mysql-8568b959d-swhv6",
                "nodemanager-6d5d44d965-wfrt9",
                "python-7c595fffcd-ggd89",
                "rollsite-77f87d7f94-zxsws"
            ]
        },
        "ID": 24,
        "CreatedAt": "2020-09-14T02:23:51Z",
        "UpdatedAt": "2020-09-14T02:24:32Z",
        "DeletedAt": null
    }
`
