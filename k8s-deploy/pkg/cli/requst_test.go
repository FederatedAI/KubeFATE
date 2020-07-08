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
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/config"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/logging"
	"github.com/dgrijalva/jwt-go"

	"github.com/spf13/viper"
)

func InitConfigForTest() {
	config.InitViper()
	viper.AddConfigPath("../../")
	viper.ReadInConfig()
}

func TestSend(t *testing.T) {
	InitConfigForTest()
	logging.InitLog()
	type args struct {
		r *Request
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				r: &Request{
					Type: "GET",
					Path: "/v1/job/",
					Body: nil,
				},
			},
			want:    200,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Send(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Code, tt.want) {
				//t.Errorf("Send() = %v, want %v", got, tt.want)
				t.Logf("Send() = %s", got.Body)
			}
		})
	}
}

func TestResponse_Unmarshal(t *testing.T) {
	type fields struct {
		Code int
		Body []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   *Result
	}{
		// TODO: Add test cases.
		{
			name: "",
			fields: fields{
				Code: 200,
				Body: []byte(`{"data":[{"uuid":"8f6b02fa-18ea-4428-b296-60b3fd56d5fe","start_time":"2020-02-14 11:23:02.7584522 +0800 CST m=+390.107448101","end_time":"","method":"ClusterInstall","Result":"","cluster_id":"","creator":"","sub-jobs":null,"status":"Running"},{"uuid":"55335543-3b56-4786-95f7-7cd09d1a9388","start_time":"2020-02-14 16:39:18.1314802 +0800 CST m=+148.350261201","end_time":"","method":"ClusterInstall","Result":"FateChart not exist","cluster_id":"a90b3906-8ed8-4a21-b14e-13d0a0f13f80","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"37c32f62-3209-4259-96b8-7ffd85d87fbb","start_time":"2020-02-14 16:56:07.3469193 +0800 CST m=+31.240052301","end_time":"","method":"ClusterInstall","Result":"create: failed to create: namespaces \"fate-cluster-10000\" not found","cluster_id":"87b20b8a-331c-406f-bfbe-a3822525c954","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"7b84be21-eed3-49bf-a574-74981751e241","start_time":"2020-02-14 17:00:00.0763482 +0800 CST m=+75.756018301","end_time":"","method":"ClusterInstall","Result":"FateChart not exist","cluster_id":"2faa3507-01ae-4da4-8c3a-ee55d3f46ba7","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"efa154a8-2d68-4789-83e2-0e2319050af5","start_time":"2020-02-15 21:50:09.2700917 +0800 CST m=+2.117446701","end_time":"","method":"ClusterInstall","Result":"","cluster_id":"","creator":"","sub-jobs":null,"status":"Running"},{"uuid":"a0f4536a-ac18-4293-8b15-1671b8e10e59","start_time":"2020-02-15 21:52:18.9889276 +0800 CST m=+2.106240301","end_time":"","method":"ClusterInstall","Result":"install success","cluster_id":"3b132495-3c5f-48d9-8c84-7f2475ace89b","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"fa46448c-05d2-422d-a629-923694480428","start_time":"2020-02-15 21:54:00.0743101 +0800 CST m=+2.109940501","end_time":"","method":"ClusterInstall","Result":"install success","cluster_id":"759610d5-3789-439a-ada6-cb47920cc502","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"8c4a7c8f-279e-4b9a-ac17-dbe2798438d2","start_time":"2020-02-15 21:56:31.6751711 +0800 CST m=+2.091454801","end_time":"","method":"ClusterInstall","Result":"install success","cluster_id":"d431d984-1b83-4902-85e8-5ae28368d91b","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"869d4c8e-ca53-41e6-8c8b-82026b1b4f51","start_time":"2020-02-15 22:00:12.4788371 +0800 CST m=+2.124143001","end_time":"","method":"ClusterInstall","Result":"install success","cluster_id":"52e0bd0b-44d8-4a49-b0df-f8fba73b1eb5","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"e947a6b0-6853-42b4-b2f5-a2f563351f72","start_time":"2020-02-15 22:03:49.5120709 +0800 CST m=+47.899916001","end_time":"","method":"ClusterInstall","Result":"create: failed to create: namespaces \"fate-cluster-10000\" not found","cluster_id":"2d437ff0-e6b8-428a-b243-eef0ebad1080","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"cbe5ab93-26cb-4b72-9851-4e689eb917d7","start_time":"2020-02-15 22:06:19.5070189 +0800 CST m=+197.894864001","end_time":"","method":"ClusterInstall","Result":"install success","cluster_id":"6b7fe334-f673-4e56-8ba2-9bc859f7f8a8","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"0ca2280a-be0a-45a8-9066-04553b2fbb7e","start_time":"2020-02-15 22:26:36.353319 +0800 CST m=+0.126001001","end_time":"","method":"ClusterInstall","Result":"install success","cluster_id":"bdf4ea5d-0287-4af2-b518-fa037dab9c68","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"13583e5e-6e2f-4239-b016-a12c83c7eee8","start_time":"2020-02-16 15:19:59.5916254 +0800 CST m=+1075.905328401","end_time":"2020-02-16 15:20:42.9963407 +0800 CST m=+1119.310043701","method":"ClusterInstall","Result":"install success","cluster_id":"2f41aabe-1610-4e4a-bc1c-9b24e9f8ec11","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"08ac7e52-5491-4b2f-8377-8d6bd415b70c","start_time":"2020-02-16 15:31:02.2466059 +0800 CST m=+121.843282801","end_time":"2020-02-16 15:31:04.6979669 +0800 CST m=+124.294643801","method":"ClusterInstall","Result":"cannot re-use a name that is still in use","cluster_id":"fba9f504-52f7-4625-b34c-f767ebcf2757","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"3ea0aeb0-6e18-4bc1-9b43-82deddb87cef","start_time":"2020-02-17 14:17:16.2831437 +0800 CST m=+62.687581901","end_time":"2020-02-17 14:17:18.7789293 +0800 CST m=+65.183367501","method":"ClusterInstall","Result":"cannot re-use a name that is still in use","cluster_id":"f6a94a0d-2b9a-4175-9e10-446e3efc60a3","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"d2bc9e20-d824-464d-aa47-c256021ebcbd","start_time":"2020-02-17 14:18:30.8854353 +0800 CST m=+137.289873501","end_time":"2020-02-17 14:18:46.2270401 +0800 CST m=+152.631478301","method":"ClusterInstall","Result":"install success","cluster_id":"8720fee1-3dc6-4df7-a03a-256cc493876c","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"ad4e7204-2a83-490e-86d7-a9d555fdeb3e","start_time":"2020-02-17 14:21:15.4652275 +0800 CST m=+301.869665701","end_time":"","method":"ClusterDelete","Result":"","cluster_id":"","creator":"","sub-jobs":null,"status":"Running"},{"uuid":"7d1be311-bd4c-4d0e-bab2-e8833b7b6e73","start_time":"2020-02-17 14:39:29.6391567 +0800 CST m=+42.051142701","end_time":"","method":"ClusterDelete","Result":"","cluster_id":"","creator":"","sub-jobs":null,"status":"Running"},{"uuid":"b7dd3d45-9bf7-42d9-a503-6f827d174c24","start_time":"2020-02-17 14:42:27.9785065 +0800 CST m=+66.362879801","end_time":"","method":"ClusterDelete","Result":"uninstall: Release not loaded: fate-10000: release: not found","cluster_id":"","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"1d940180-582d-4d61-9f3d-62c8811f066e","start_time":"2020-02-17 14:44:54.8065286 +0800 CST m=+213.190901901","end_time":"2020-02-17 14:45:08.3930903 +0800 CST m=+226.777463601","method":"ClusterInstall","Result":"install success","cluster_id":"2edd47aa-318b-4575-a4d8-b21befb699f0","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"22faa4f9-5c3f-4429-827b-928467f54304","start_time":"2020-02-17 15:25:37.9983955 +0800 CST m=+106.460258201","end_time":"2020-02-17 15:25:41.0169141 +0800 CST m=+109.478776801","method":"ClusterInstall","Result":"cannot re-use a name that is still in use","cluster_id":"aedff23b-2b3a-4367-ae26-d333ce40cc54","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"b232b44c-45b7-44e0-8b98-99fab9146fa6","start_time":"2020-02-17 15:27:23.4039901 +0800 CST m=+211.865852801","end_time":"2020-02-17 15:27:36.3299872 +0800 CST m=+224.791849901","method":"ClusterInstall","Result":"install success","cluster_id":"f41ffdc9-5f67-4dea-8257-be49be8be972","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"b730d00f-21a8-46e4-a227-d9a4b5c192be","start_time":"2020-02-18 11:30:37.6692103 +0800 CST m=+0.113000701","end_time":"2020-02-18 11:30:50.9305096 +0800 CST m=+13.374300001","method":"ClusterInstall","Result":"Service \"proxy\" is invalid: spec.ports[0].nodePort: Invalid value: 30010: provided port is already allocated","cluster_id":"a8838735-bc5a-427d-a376-37ade43d972c","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"79cc3bcf-fcfd-4fc1-b9bc-86c869c2e2ea","start_time":"2020-02-18 11:32:16.4349365 +0800 CST m=+0.115998301","end_time":"2020-02-18 11:32:35.3309205 +0800 CST m=+19.011982301","method":"ClusterInstall","Result":"install success","cluster_id":"f5248151-5d2d-409c-a88e-26a0703c05a2","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"c9722dc1-fd5a-4b7b-8415-f48d6268ebf5","start_time":"2020-02-18 11:46:57.5363403 +0800 CST m=+0.135000801","end_time":"","method":"ClusterInstall","Result":"","cluster_id":"","creator":"","sub-jobs":null,"status":"Running"},{"uuid":"dfdc8d3c-11e9-43e8-9c1c-8fe2b770d0a2","start_time":"2020-02-18 11:47:56.069797 +0800 CST m=+0.241000901","end_time":"","method":"ClusterDelete","Result":"","cluster_id":"","creator":"","sub-jobs":null,"status":"Running"},{"uuid":"78955586-ce3e-4463-a0e1-68b0884665b4","start_time":"2020-02-18 11:49:05.8571065 +0800 CST m=+0.090997601","end_time":"","method":"ClusterDelete","Result":"","cluster_id":"","creator":"","sub-jobs":null,"status":"Running"},{"uuid":"2d50163c-188c-49cf-ae06-9a7bb42bb603","start_time":"2020-02-18 11:57:45.7947982 +0800 CST m=+0.096005701","end_time":"2020-02-18 11:57:48.8689033 +0800 CST m=+3.170110801","method":"ClusterDelete","Result":"uninstall: Release not loaded: fate-8888: release: not found","cluster_id":"","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"057261fe-f06d-483d-b147-193e52850d30","start_time":"2020-02-18 12:09:41.5676 +0800 CST m=+0.112003301","end_time":"2020-02-18 12:09:56.8979373 +0800 CST m=+15.442340601","method":"ClusterInstall","Result":"install success","cluster_id":"993c55bc-c291-45ac-990a-e9ccc37373a0","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"d499b4dc-256e-4e10-a7fa-f19dffd27fbd","start_time":"2020-02-18 12:11:06.1289336 +0800 CST m=+0.101999301","end_time":"2020-02-18 12:11:13.3933879 +0800 CST m=+7.366453601","method":"ClusterDelete","Result":"uninstall success","cluster_id":"","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"d8ec587d-f386-4750-a98c-5b4702c26a3f","start_time":"2020-02-18 12:13:49.0199955 +0800 CST m=+0.095003401","end_time":"2020-02-18 12:13:52.1121072 +0800 CST m=+3.187115101","method":"ClusterDelete","Result":"uninstall: Release not loaded: fate-8888: release: not found","cluster_id":"","creator":"","sub-jobs":null,"status":"Retry"},{"uuid":"65f8bffe-eac4-4840-9588-944cf9314cee","start_time":"2020-02-18 12:14:46.9458639 +0800 CST m=+0.131999201","end_time":"2020-02-18 12:15:00.8571077 +0800 CST m=+14.043243001","method":"ClusterInstall","Result":"install success","cluster_id":"5f9a6c0d-22af-4186-89ed-8d96de919b09","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"7792e0e1-0f83-4df8-9c42-e894d3c8ffa2","start_time":"2020-02-18 12:17:17.3667675 +0800 CST m=+0.090001701","end_time":"2020-02-18 12:17:21.3797353 +0800 CST m=+4.102969501","method":"ClusterInstall","Result":"cannot re-use a name that is still in use","cluster_id":"1740d368-feef-4803-a0b1-27495e40d0a7","creator":"","sub-jobs":null,"status":"Retry"},{"uuid":"8e5357fd-68ea-4469-b903-b8cf2cfb2de9","start_time":"2020-02-18 12:18:24.9261303 +0800 CST m=+0.138004201","end_time":"2020-02-18 12:18:32.7373439 +0800 CST m=+7.949217801","method":"ClusterDelete","Result":"uninstall success","cluster_id":"","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"a69e26c4-56e2-41fe-9045-73ef2595f280","start_time":"2020-02-18 12:19:43.3418927 +0800 CST m=+0.106998401","end_time":"2020-02-18 12:19:46.4474331 +0800 CST m=+3.212538801","method":"ClusterDelete","Result":"uninstall: Release not loaded: fate-8888: release: not found","cluster_id":"","creator":"","sub-jobs":null,"status":"Retry"},{"uuid":"e9ccaebb-584c-4658-8f9b-b4cd8bfdbdf4","start_time":"2020-02-18 12:20:56.2128401 +0800 CST m=+0.099995701","end_time":"2020-02-18 12:21:10.4355801 +0800 CST m=+14.322735701","method":"ClusterInstall","Result":"install success","cluster_id":"7fa7cda4-3b09-4a69-ade1-80f67bd36076","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"1c042f5c-0f3e-4876-bc5c-05f423795059","start_time":"2020-02-18 12:47:58.8161725 +0800 CST m=+0.090996801","end_time":"2020-02-18 12:48:07.9424756 +0800 CST m=+9.217299901","method":"ClusterDelete","Result":"uninstall success","cluster_id":"","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"8ff44c19-d045-4b7d-baaa-29ef2c97b4f5","start_time":"2020-02-18 12:49:27.0520957 +0800 CST m=+0.091002101","end_time":"2020-02-18 12:49:30.617687 +0800 CST m=+3.656593401","method":"ClusterDelete","Result":"uninstall: Release not loaded: fate-8888: release: not found","cluster_id":"","creator":"","sub-jobs":null,"status":"Retry"},{"uuid":"d0a21b14-591f-46a2-a777-8edee24102b8","start_time":"2020-02-18 12:50:42.8185435 +0800 CST m=+0.085999901","end_time":"2020-02-18 12:50:57.2622939 +0800 CST m=+14.529750301","method":"ClusterInstall","Result":"install success","cluster_id":"01aa5b0c-c7ab-4a95-b92a-90b864130289","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"8228631f-a472-449d-bd10-d472d4df7699","start_time":"2020-02-18 12:51:35.6578933 +0800 CST m=+0.098997901","end_time":"2020-02-18 12:51:43.4941502 +0800 CST m=+7.935254801","method":"ClusterDelete","Result":"uninstall success","cluster_id":"","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"76887071-39ae-4685-93b8-eeb442e89f4e","start_time":"2020-02-18 12:53:32.8347262 +0800 CST m=+0.095999201","end_time":"2020-02-18 12:53:48.1450823 +0800 CST m=+15.406355301","method":"ClusterInstall","Result":"install success","cluster_id":"43d66a1d-21d2-4782-a305-53450acf4910","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"b2387723-7f08-47d2-a4c4-a3383dc7fda0","start_time":"2020-02-18 12:54:20.374074 +0800 CST m=+0.098005901","end_time":"2020-02-18 12:54:28.3816648 +0800 CST m=+8.105596701","method":"ClusterDelete","Result":"uninstall success","cluster_id":"","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"29998ed5-00fd-44a2-a47f-9cede9723ff5","start_time":"2020-02-18 12:56:20.6749909 +0800 CST m=+0.092001001","end_time":"","method":"ClusterDelete","Result":"","cluster_id":"","creator":"","sub-jobs":null,"status":"Success"},{"uuid":"fdaf44ba-3ac7-45f4-bb3b-820840bd9e6c","start_time":"2020-02-18 12:57:39.4084842 +0800 CST m=+0.088999401","end_time":"","method":"ClusterDelete","Result":"cluster no find","cluster_id":"","creator":"","sub-jobs":null,"status":"Retry"},{"uuid":"dfe6329c-35ee-4ab9-8438-fd96d6ed1a11","start_time":"2020-02-18 13:08:25.2802667 +0800 CST m=+0.097001201","end_time":"2020-02-18 13:08:27.3268279 +0800 CST m=+2.143562401","method":"ClusterDelete","Result":"cluster no find","cluster_id":"","creator":"","sub-jobs":null,"status":"Retry"},{"uuid":"d78502e5-b41c-4c37-a2df-aef38b27d30f","start_time":"2020-02-18 13:10:11.0094312 +0800 CST m=+0.163996901","end_time":"2020-02-18 13:10:25.4892006 +0800 CST m=+14.643766301","method":"ClusterInstall","Result":"install success","cluster_id":"5029628c-8886-4907-bced-6dbe3553c7ef","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"0f032ea5-e017-459b-925c-9457f10aa0b4","start_time":"2020-02-18 13:11:34.7557516 +0800 CST m=+0.092996901","end_time":"2020-02-18 13:11:42.4700184 +0800 CST m=+7.807263701","method":"ClusterDelete","Result":"uninstall success","cluster_id":"","creator":"","sub-jobs":null,"status":"Failed"},{"uuid":"f6a9ea74-fd99-44a7-85cf-7d4bc90c4531","start_time":"2020-02-18 13:12:37.0458631 +0800 CST m=+0.092997601","end_time":"2020-02-18 13:12:40.1603877 +0800 CST m=+3.207522201","method":"ClusterDelete","Result":"cluster no find","cluster_id":"","creator":"","sub-jobs":null,"status":"Retry"}],"msg":"getJobList success"}`),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				Code: tt.fields.Code,
				Body: tt.fields.Body,
			}
			if got := r.Unmarshal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response.Unmarshal() = %+v, want %v\n", got, tt.want)
				for _, v := range got.Data {
					t.Logf("%+v\n", v)
				}
			}
		})
	}
}

func TestGetToken(t *testing.T) {
	claims := &jwt.MapClaims{
		"id":       "admin",
		"exp":      time.Now().Add(300000 * time.Second).Unix(),
		"orig_iat": time.Now().Add(300000 * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //生成token
	accessToken, err := token.SignedString([]byte("secret key"))
	if err != nil {
		return
	}
	fmt.Println(accessToken)
}
