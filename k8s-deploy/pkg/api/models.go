/*
 * Copyright 2019-2021 VMware, Inc.
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

package api

// JSONResult Success Result
type JSONResult struct {
	Message string `json:"msg" example:"Success"`
}

// JSONERRORResult error Result
type JSONERRORResult struct {
	Code  int    `json:"code" `
	Error string `json:"error"`
}

// JSONEMSGResult 401 Result
type JSONEMSGResult struct {
	Code    int    `json:"code" example:"401"`
	Message string `json:"msg" example:"cookie token is empty"`
}

// VersionResult GET Version Result
type VersionResult struct {
	Version string `example:"v1.3.0"`
	Message string `json:"msg" example:"getVersion Success"`
}

//TokenResult LOgin Result
type TokenResult struct {
	Code   int    `example:"200"`
	Expire string `example:"2021-01-28T14:48:53+08:00"`
	Token  string `example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Login body
type Login struct {
	Username string `example:"admin"`
	Password string `example:"admin"`
}
