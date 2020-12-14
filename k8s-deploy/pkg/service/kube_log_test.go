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

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
)

func TestGetLogs(t *testing.T) {
	type args struct {
		args *LogChanArgs
	}
	tests := []struct {
		name    string
		args    args
		want    *bytes.Buffer
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				args: &LogChanArgs{
					Name:       "fate-9999",
					Namespace:  "fate-9999",
					Container:  "python",
					TailLines:  func() *int64 { a := int64(10); return &a }(),
					Timestamps: true,
					SinceTime:  time.Now().Add(-10*time.Hour),
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLogs(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLogs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readLogToString(t *testing.T) {
	logRead, err := getLogFollowOfModule("fate-9999", "fate-9999", "python")
	if err != nil {
		log.Debug().Err(err).Msg("GetLogFollowOfModule")
		return
	}
	log.Debug().Msg("start")
	type args struct {
		logRead io.ReadCloser
		prefix  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				logRead: logRead,
				prefix:  "python",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readLogToString(tt.args.logRead, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("readLogToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readLogToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
