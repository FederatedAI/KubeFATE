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
package service

//func TestUpgrades(t *testing.T) {
//	_ = config.InitViper()
//	viper.AddConfigPath("../../")
//	_ = viper.ReadInConfig()
//	logging.InitLog()
//	_ = os.Setenv("FATECLOUD_CHART_PATH", "../../")
//	type args struct {
//		namespace string
//		name      string
//		version   string
//		value     *Value
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//		{
//			name:    "",
//			args:    args{
//				namespace: "fate-10000",
//				name:      "fate-10000",
//				version:   "v1.2.0",
//				value:     &Value{
//					Val: []byte(`{ "partyId":10000,"endpoint": { "ip":"192.168.100.123","port":30030}}`),
//					T:   "json",
//				},
//			},
//			wantErr: false,
//		},
//	}
//	//for _, tt := range tests {
//	//	t.Run(tt.name, func(t *testing.T) {
//	//		if err := Upgrade(tt.args.namespace, tt.args.name, tt.args.version, tt.args.value); (err != nil) != tt.wantErr {
//	//			t.Errorf("Upgrade() error = %v, wantErr %v", err, tt.wantErr)
//	//		}
//	//	})
//	//}
//}
