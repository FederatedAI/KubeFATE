# Copyright 2019-2020 VMware, Inc.
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# you may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from fml_manager import ClusterManager
if __name__ == "__main__":
    # init the cluster manager
    cluster_manager = ClusterManager("fate-9999", "fatesample")

    # get route table
    route_table=cluster_manager.GetRouteTable()
    route_table["route_table"]["9999"] = {'default': [{'ip': '192.168.0.1', 'port': 9370}]}

    # update route table
    cluster_manager.SetRouteTable(route_table)

    # get entrypoint
    cluster_manager.GetEntrypoint()