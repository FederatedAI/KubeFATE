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

import json
import os
import tarfile
import requests
import base64
import random
import time
import subprocess
import tempfile

import pandas as pd
from contextlib import closing
from fml_manager.utils import file_utils
from fml_manager.utils.core import get_lan_ip

cFateFlowHostEnv = "FATE_FLOW_HOST"
cFateServingHostEnv = "FATE_SERVING_HOST"

cFateFlowServieName = "fateflow"
cFateServingServieName = "fateserving"
cFateFlowServiePort = 9380
cFateServingServiePort = 8059

cFateClusterCR = "fatecluster"


class FMLManager:
    """FMLManager is used to communicate with FATE cluster"""

    def __init__(self, server_conf=None, log_path="./"):
        """ Init the FMLManager with config and log path

        :param server_conf: Path to config file, default=None
        :type server_conf: string
        :param log_path: Path to log file, default=./
        :type log_path: string

        """

        #: Url of FATE Flow service
        self.server_url = None

        #: Url of FATE Serving service
        self.serving_url = None
        self.log_path = log_path

        if server_conf is not None:
            self._init_from_config(server_conf)
        elif os.getenv(cFateFlowHostEnv) is not None and os.getenv(cFateFlowHostEnv) != "":
            self._init_from_env()
        else:
            self._init_from_kube_api()

        # if the server url is still None, the initialization is failed
        if self.server_url is None:
            raise Exception(
                "Unable to find fate_flow url, failed to initialize the FML Manager")
        if self.serving_url is None:
            print("Unable to find fate_serving url but it's ok to continue")

    def _init_from_config(self, server_conf):
        self.server_conf = file_utils.load_json_conf(server_conf)
        self.ip = self.server_conf.get("servers").get("fateflow").get("host")
        self.serving_url = self.server_conf.get("servings")
        self.http_port = self.server_conf.get(
            "servers").get("fateflow").get("http.port")
        self.server_url = "http://{}:{}/{}".format(
            self.ip, self.http_port, "v1")

    def _init_from_env(self):
        server_host = os.getenv(cFateFlowHostEnv)
        serving_host = os.getenv(cFateServingHostEnv, "")

        self.server_url = "http://{}/{}".format(server_host, "v1")
        self.serving_url = "http://{}".format(serving_host)

    def _init_from_kube_api(self):
        args = "kubectl get {} -A -o json".format(cFateClusterCR).split(" ")
        try:
            data, err = subprocess.Popen(
                args, stdout=subprocess.PIPE).communicate()
            data_json = json.loads(data)
            if len(data_json["items"]) != 0:
                # fetch the first fatecluster by default
                fate_cluster_namespace = data_json["items"][0]["metadata"]["namespace"]
                self.server_url = "http://{}.{}:{}/{}".format(
                    cFateFlowServieName, fate_cluster_namespace, cFateFlowServiePort, "v1")
                self.serving_url = "http://{}.{}:{}/{}".format(
                    cFateServingServieName, fate_cluster_namespace, cFateServingServiePort, "v1")
        except Exception as e:
            print(e)

    # Job management

    def submit_job(self, dsl, config):
        """ Submit job to FATE cluster

        :param dsl: DSL definition
        :type dsl: Pipline
        :param config: Config definition
        :type config: Config

        :returns: response
        :rtype: dict

        """
        post_data = {'job_dsl': dsl,
                     'job_runtime_conf': config}
        response = requests.post(
            "/".join([self.server_url, "job", "submit"]), json=post_data)

        return self.prettify(response)

    def submit_job_by_files(self, dsl_path, config_path):
        """ Submit job with file to FATE cluster

        :param dsl_path: DSL definition path
        :type dsl_path: string
        :param config_path: Config definition path
        :type config_path: string

        :returns: response
        :rtype: dict

        """

        config_data = {}
        if config_path:
            config_path = os.path.abspath(config_path)
            with open(config_path, 'r') as f:
                config_data = json.load(f)
        else:
            raise Exception('Conf cannot be null.')
        dsl_data = {}
        if dsl_path:
            dsl_path = os.path.abspath(dsl_path)
            with open(dsl_path, 'r') as f:
                dsl_data = json.load(f)
        else:
            raise Exception('DSL_path cannot be null.')

        return self.submit_job(dsl_data, config_data)

    def query_job_status(self, query_conditions, max_tries=200):
        """ Fetch status of job

        :param query_conditions: Condition of the job
        :type query_conditions: dict

        :returns: response
        :rtype: dict

        """
        job_status = "failed"
        for i in range(max_tries):
            time.sleep(1)
            try:
                guest_status = self.query_job(query_conditions).json()[
                    "data"][0]["f_status"]
            except Exception as e:
                print("Failed to fetch status: ", e)

            print("Status: %s" % guest_status)
            if guest_status == "failed":
                job_status = "failed"
                raise Exception("Failed to upload data.")
            if guest_status == "success":
                job_status = "success"
                break
        return job_status

    def query_job(self, query_conditions):
        """ Fetch job

        :param query_conditions: Condition of the job
        :type query_conditions: QueryCondition

        :returns: response
        :rtype: dict

        """

        response = requests.post(
            "/".join([self.server_url, "job", "query"]), json=query_conditions.to_dict())
        return self.prettify(response)

    def query_job_conf(self, query_conditions):
        """ Fetch config of job

        :param query_conditions: Condition of the job
        :type query_conditions: dict

        :returns: response
        :rtype: dict

        """

        response = requests.post(
            "/".join([self.server_url, "job", "config"]), json=query_conditions)
        return self.prettify(response)

    def stop_job(self, job_id):
        """ Stop job

        :param job_id: job id of data
        :type job_id: string

        :returns: response
        :rtype: dict

        """
        post_data = {
            'job_id': job_id
        }
        response = requests.post(
            "/".join([self.server_url, "job", "stop"]), json=post_data)
        return self.prettify(response)

    def update_job(self, job_id, role, party_id, notes):
        post_data = {
            "job_id": job_id,
            "role": role,
            "party_id": party_id,
            "notes": notes
        }
        response = requests.post(
            "/".join([self.server_url, "job", "update"]), json=post_data)
        return self.prettify(response)

    def fetch_job_log(self, job_id):
        """ Fetch the log of job

        :param job_id: The UUID of job
        :type job_id: string

        :returns: response
        :rtype: dict

        """
        data = {
            "job_id": job_id
        }

        tar_file_name = 'job_{}_log.tar.gz'.format(job_id)
        extract_dir = os.path.join(self.log_path, 'job_{}_log'.format(job_id))
        with closing(requests.get("/".join([self.server_url, "job", "log"]), json=data,
                                  stream=True)) as response:
            if response.status_code == 200:
                self.__download_from_request(
                    http_response=response, tar_file_name=tar_file_name, extract_dir=extract_dir)
                response = {'retcode': 0,
                            'directory': extract_dir,
                            'retmsg': 'download successfully, please check {} directory, file name is {}'.format(extract_dir, tar_file_name)}

                return self.prettify(response, True)
            else:
                return self.prettify(response, True)

    # Data management
    def load_data(self, url, namespace, table_name, work_mode, head, partition, drop="1", api_version="1.4"):
        """ Upload data to FATE cluster

        :param url: URL of data to upload
        :type url: string

        :param namespace: Namespace of the data in FATE cluster
        :type namespace: string

        :param table_name: Table name of the data in FATE cluster
        :type table_name: string

        :param work_mode: The work mode of upload
        :type work_mode: int

        :param head: Head included flag. '1': with head, '0': without head.
        :type head: int

        :param partition: Partitions of the upload data
        :type partition: int

        :param drop: Flag to overwrite data with same identifier
        :type drop: string

        :returns: response
        :rtype: dict

        """
        if api_version == "1.4":
            temp_file = None
            if url.startswith("http://") or url.startswith("https://"):
                downloader = HttpDownloader(url)
                temp_file = downloader.download_to(
                    file_utils.get_project_base_directory())
                url = temp_file

            post_data = {
                "namespace": namespace,
                "table_name": table_name,
                "work_mode": work_mode,
                "head": head,
                "partition": partition,
                "drop": drop
            }
            data_files = {
                "file": open(url, "rb")
            }
            response = requests.post(
                "/".join([self.server_url, "data", "upload"]), params=post_data, files=data_files)

            if temp_file is not None and os.path.exists(temp_file):
                print("Delete temp file...")
                os.remove(temp_file)
        else:
            post_data = {
                "file": url,
                "namespace": namespace,
                "table_name": table_name,
                "work_mode": work_mode,
                "head": head,
                "partition": partition
            }
            response = requests.post(
                "/".join([self.server_url, "data", "upload"]), json=post_data)

        return self.prettify(response)

    def query_data(self, job_id, limit):
        """ Query data of job
        """
        post_data = {
            "job_id": job_id,
            "limit": limit
        }

        response = requests.post(
            "/".join([self.server_url, "data", "upload", "history"]), json=post_data)

        return self.prettify(response)

    # The data is download to fateflow. FATE not ready to download to local.
    def download_data(self, namespace, table_name, filename, work_mode, delimitor, output_folder="./"):
        """ Download data to local
        """

        DEFAULT_DATA_FOLDER = "/data/projects/fate/python/download_dir"
        output_path = "{}/{}".format(DEFAULT_DATA_FOLDER, filename)
        post_data = {
            "namespace": namespace,
            "table_name": table_name,
            "work_mode": work_mode,
            "delimitor": delimitor,
            "output_path": output_path
        }
        response = requests.post(
            "/".join([self.server_url, "data", "download"]), json=post_data)

        if response.status_code == 200:
            output = json.loads(response.content)
            job_id = output["jobId"]
            query_condition = {
                "job_id": job_id
            }
            for i in range(500):
                time.sleep(1)
                status = self.query_job(query_condition).json()[
                    "data"][0]["f_status"]
                if status == "failed":
                    print("Failed")
                    print(self.query_job(query_condition).json())
                    raise Exception("Failed to download data.")
                if status == "success":
                    return self.prettify(response, True)
        response = {
            'retcode': 1,
            'retmsg': 'Download failed'
        }

        return self.prettify(response, True)

    # Model management
    def load_model(self, initiator_party_id, federated_roles, work_mode, model_id, model_version):
        post_data = {
            "initiator": {
                "party_id": initiator_party_id,
                "role": "guest"
            },
            "role": federated_roles,
            "job_parameters": {
                "work_mode": work_mode,
                "model_id": model_id,
                "model_version": model_version
            }
        }
        response = requests.post(
            "/".join([self.server_url, "model", "load"]), json=post_data)

        return self.prettify(response)

    def bind_model(self, service_id, initiator_party_id, federated_roles, work_mode, model_id, model_version):
        if self.serving_url == nil:
            raise Exception(
                'Federated Serving is not deployed or not correctly configured yet. ')
        post_data = {
            "service_id": service_id,
            "initiator": {
                "party_id": initiator_party_id,
                "role": "guest"
            },
            "role": federated_roles,
            "job_parameters": {
                "work_mode": work_mode,
                "model_id": model_id,
                "model_version": model_version
            },
            "servings": self.serving_url
        }

        response = requests.post(
            "/".join([self.server_url, "model", "bind"]), json=post_data)
        return self.prettify(response)

    def print_model_version(self, role, party_id, model_id, api_version="1.4"):
        """ Print model version
        """

        action = "version"
        if api_version == "1.4":
            action = "version_history"

        namespace = "#".join([role, str(party_id), model_id])
        post_data = {
            "namespace": namespace
        }

        response = requests.post(
            "/".join([self.server_url, "model", action]), json=post_data)

        return self.prettify(response, True)

    def model_output(self, role, party_id, model_id, model_version, model_component):
        """ Output the model
        """

        namespace = "#".join([role, str(party_id), model_id])
        post_data = {
            "name": model_version,
            "namespace": namespace
        }
        response = requests.post(
            "/".join([self.server_url, "model", "transfer"]), json=post_data)
        model = json.loads(response.content)
        if model["data"] != "":
            en_model_metadata = model["data"]["%sMeta" % model_component]
            en_model_parameters = model["data"]["%sParam" % model_component]

        model = {
            "metadata": en_model_metadata,
            "parameters": en_model_parameters
        }

        return self.prettify(model, True)

    def offline_predict_on_dataset(self, is_vertical, initiator_party_role, initiator_party_id, work_mode, model_id, model_version, federated_roles, guest_data_name="", guest_data_namespace="", host_data_name="", host_data_namespace=""):
        if is_vertical:
            print("This API is not support vertical federated machine learning yet. ")
            return

        # For predict job, dsl is empty dict.
        dsl = {}

        config = {
            "initiator": {
                "role": initiator_party_role,
                "party_id": initiator_party_id
            },
            "job_parameters": {
                "work_mode": work_mode,
                "job_type": "predict",
                "model_id": model_id,
                "model_version": model_version
            },
            "role": federated_roles,
            "role_parameters": {}
        }

        if guest_data_name != "" or guest_data_namespace != "":
            if initiator_party_role != "guest":
                raise Exception("Initiator not has data sets.")

            guest_parameters = {
                "args": {
                    "data": {
                        "eval_data": [{"name": guest_data_name, "namespace": guest_data_namespace}]
                    }
                }
            }

            config["role_parameters"]["guest"] = guest_parameters

        if host_data_name != "" or host_data_namespace != "":
            host_parameters = {
                "args": {
                    "data": {
                        "eval_data": [{"name": host_data_name, "namespace": host_data_namespace}]
                    }
                }
            }
            config["role_parameters"]["host"] = guest_parameters

        return self.submit_job(dsl, config)

    # Task
    def query_task(self, query_conditions):
        """ Query task
        """

        response = requests.post(
            "/".join([self.server_url, "job", "task", "query"]), json=query_conditions)
        return self.prettify(response)

    # Tracking
    def track_job_data(self, job_id, role, party_id):
        """ Track job data
        """

        post_data = {
            "job_id": job_id,
            "role": role,
            "party_id": party_id
        }

        response = requests.post(
            "/".join([self.server_url, "tracking", "job", "data_view"]), json=post_data)
        return self.prettify(response, True)

    def track_component_all_metric(self, job_id, role, party_id, component_name):
        """ Track output all metric of component
        """
        post_data = {
            "job_id": job_id,
            "role": role,
            "party_id": party_id,
            "component_name": component_name
        }

        response = requests.post(
            "/".join([self.server_url, "tracking", "component", "metric", "all"]), json=post_data)
        return self.prettify(response, True)

    def track_component_metric_type(self, job_id, role, party_id, component_name):
        """ Track output metric type of component
        """
        post_data = {
            "job_id": job_id,
            "role": role,
            "party_id": party_id,
            "component_name": component_name
        }

        response = requests.post(
            "/".join([self.server_url, "tracking", "component", "metrics"]), json=post_data)
        return self.prettify(response, True)

    """
    metric_name and metric_namespace can be found in API track_component_metric_type
    e.g. response = manager.track_component_metric_type(jobId, "guest", "10000", "homo_lr_0")
        {
            "data": {
                "train": [
                    "loss"
                ]
            },
            "retcode": 0,
            "retmsg": "success"
        }

        The metric_name is "loss" and metric_namespace is "train"
    """

    def track_component_metric_data(self, job_id, role, party_id, component_name, metric_name, metric_namespace):
        """ Track output metric data of component
        """
        post_data = {
            "job_id": job_id,
            "role": role,
            "party_id": party_id,
            "component_name": component_name,
            "metric_name": metric_name,
            "metric_namespace": metric_namespace
        }

        response = requests.post(
            "/".join([self.server_url, "tracking", "component", "metric_data"]), json=post_data)
        return self.prettify(response, True)

    def track_component_parameters(self, job_id, role, party_id, component_name):
        """ Track output parameter of component
        """
        post_data = {
            "job_id": job_id,
            "role": role,
            "party_id": party_id,
            "component_name": component_name
        }

        response = requests.post(
            "/".join([self.server_url, "tracking", "component", "parameters"]), json=post_data)
        return self.prettify(response, True)

    def track_component_output_model(self, job_id, role, party_id, component_name):
        """ Track output model of component
        """
        post_data = {
            "job_id": job_id,
            "role": role,
            "party_id": party_id,
            "component_name": component_name
        }

        response = requests.post(
            "/".join([self.server_url, "tracking", "component", "output", "model"]), json=post_data)
        return self.prettify(response, True)

    def track_component_output_data(self, job_id, role, party_id, component_name):
        """ Track output data of component

        :rtype: pandas.DataFrame
        """
        post_data = {
            "job_id": job_id,
            "role": role,
            "party_id": party_id,
            "component_name": component_name
        }

        response = requests.post(
            "/".join([self.server_url, "tracking", "component", "output", "data"]), json=post_data)

        result = response.json()
        data = result['data']
        header = result['meta']['header']

        return pd.DataFrame(data, columns=header)

    # Utils
    def prettify(self, response, verbose=False):
        if verbose:
            if isinstance(response, requests.Response):
                if response.status_code == 200:
                    print("Success!")
                print(json.dumps(response.json(), indent=4, ensure_ascii=False))
            else:
                print(response)

        return response

    def __download_data_from_request(self, http_response, output):
        with open(output, 'wb') as fw:
            for chunk in http_response.iter_content(1024):
                if chunk:
                    fw.write(chunk)

    def __download_from_request(self, http_response, tar_file_name, extract_dir):
        with open(tar_file_name, 'wb') as fw:
            for chunk in http_response.iter_content(1024):
                if chunk:
                    fw.write(chunk)
        tar = tarfile.open(tar_file_name, "r:gz")
        file_names = tar.getnames()
        for file_name in file_names:
            tar.extract(file_name, extract_dir)
        tar.close()
        os.remove(tar_file_name)


class HttpDownloader:
    def __init__(self, url):
        self.url = url

    def download_to(self, path_to_save):
        r = requests.get(self.url, allow_redirects=True)
        filename = self.__get_filename_from_cd(
            r.headers.get('content-disposition'))
        temp_file_to_write = os.path.join(
            file_utils.get_project_base_directory(), filename)
        open(temp_file_to_write, 'wb').write(r.content)

        return temp_file_to_write

    def __get_filename_from_cd(self, cd):
        """
        Get filename from content-disposition
        """
        if not cd:
            # just return file name
            fname = self.url.split('/')[-1]
            if len(fname) == 0:
                return None
            return fname
        fname = re.findall('filename=(.+)', cd)
        if len(fname) == 0:
            return None
        return fname[0]
