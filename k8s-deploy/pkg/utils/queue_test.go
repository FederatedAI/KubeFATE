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
package utils

import (
	"testing"
)

func TestQueue(t *testing.T) {
	q := NewQueue(8)
	a := "1"
	ok, quantity := q.Put(&a)
	if !ok {
		t.Error("TestStack Get.Fail")
		return
	} else {
		t.Logf("TestStack Put value:%d[%v], quantity:%v\n", &a, a, quantity)
	}
	b := "2"
	ok, quantity = q.Put(&b)
	if !ok {
		t.Error("TestStack Get.Fail")
		return
	} else {
		t.Logf("TestStack Put value:%d[%v], quantity:%v\n", &b, b, quantity)
	}
	c := "3"
	ok, quantity = q.Put(&c)
	if !ok {
		t.Error("TestStack Get.Fail")
		return
	} else {
		t.Logf("TestStack Put value:%d[%v], quantity:%v\n", &c, c, quantity)
	}

	val, ok, quantity := q.Get()
	if !ok {
		t.Error("TestStack Get.Fail")
		return
	} else {
		t.Logf("TestStack Get value:%d[%v], quantity:%v\n", val, *(val.(*string)), quantity)
	}
	val, ok, quantity = q.Get()
	if !ok {
		t.Error("TestStack Get.Fail")
		return
	} else {
		t.Logf("TestStack Get value:%d[%v], quantity:%v\n", val, *(val.(*string)), quantity)
	}
	val, ok, quantity = q.Get()
	if !ok {
		t.Error("TestStack Get.Fail")
		return
	} else {
		t.Logf("TestStack Get value:%d[%v], quantity:%v\n", val, *(val.(*string)), quantity)
	}
	if q := q.Quantity(); q != 0 {
		t.Errorf("Quantity Error: [%v] <>[%v]", q, 0)
	}
}
