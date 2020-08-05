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
from enum import Enum

"""Module to Build Complex Data Structural in FATE"""


class PartyType(Enum):

    """PartyType is used to identified the type of party"""

    #: For normal party
    NORMAL = 1

    #: For exchange
    EXCHANGE = 2


class Party():
    """Party is used to define route table record"""

    def __init__(self, p_id='', ip=None, port=None, p_type=PartyType.NORMAL) -> None:
        """ Return default party instance"""

        #: ID of the party
        if p_type is PartyType.EXCHANGE:
            self._id = 'default'
        else:
            self._id = p_id

        #: IP address of the party
        self._ip = ip

        #: Service port of the party
        self._port = port

        #: Type of the party
        self._type = p_type

    def set_id(self, party_id) -> None:
        """ Set party ID

        :param party_id: The ID of party
        :type party_id: string

        """

        self._id = party_id

    def set_ip(self, party_ip) -> None:
        """ Set IP address of party

        :param party_ip: The IP address of party
        :type party_ip: string

        """

        self._ip = party_ip

    def set_port(self, party_port) -> None:
        """ Set service port of the party

        :param party_port: The service port of party
        :type party_port: int

        """

        self._port = party_port

    def set_type(self, party_type) -> None:
        """ Set party type

        :param party_id: The type of party
        :type party_type: PartyType

        """

        self._type = party_type

    def get_id(self) -> str:
        """ Get party ID

        :rtype: string

        """

        return self._id

    def to_entry_point(self) -> dict:
        """ Get entrypoint of party

        :rtype: dict

        """
        return {'default': [{'ip': self._ip, 'port': self._port}]}


class PartyBuilder():
    """PartyBuilder is used to construct party instance"""

    def __init__(self) -> None:
        """ Get default party instance
        """

        self.reset()

    def reset(self) -> None:
        """ Reset the previous instance
        """

        self._party = Party()

    def with_id(self, party_id):
        """ Update ID

        :param party_id: The ID of party
        :type party_id: string

        """

        self._party.set_id(party_id)
        return self

    def with_ip(self, party_ip):
        """ Update IP

        :param party_ip: The IP of party
        :type party_ip: string

        """

        self._party.set_ip(party_ip)
        return self

    def with_port(self, party_port):
        """ Update port

        :param party_port: The port of party
        :type party_port: int

        """

        self._party.set_port(party_port)
        return self

    def with_type(self, party_type):
        """ Update type

        :param party_type: The type of party
        :type party_type: PartyType

        """

        self._party.set_type(party_type)
        return self

    def build(self):
        """ Return party instance with config

        :rtype: Party

        """
        party = self._party
        self.reset()

        # reset id to 'default' if party is exchange
        if party._type == PartyType.EXCHANGE:
            party._id = 'default'

        return party


class RouteTable():
    """RouteTable is used to communicate with other parties"""

    def __init__(self) -> None:
        """ Return instance with empty route table
        """

        #: The underlying route table config file
        self._route_table = None

    def add_parties(self, *parties) -> None:
        """ Append parties to route table

        :param parties: A list of party instance
        :type parties: list

        """

        for party in parties:
            self._route_table['route_table'][party.get_id()] = party.to_entry_point(
            )

    def update_parties(self, *parties) -> None:
        """ Update parties of route table

        :param parties: A list of party instance
        :type parties: list

        """

        for party in parties:
            self._route_table['route_table'][party.get_id()] = party.to_entry_point(
            )

    def remove_parties(self, *party_ids) -> None:
        """ Remove parties from route table

        :param parties: A list of party ID
        :type parties: list

        """

        for party_id in party_ids:
            if(self._route_table['route_table'].get(party_id) != None):
                self._route_table['route_table'].pop(party_id)

    def get_parties(self) -> dict:
        """ List all parties

        :rtype: dict

        """
        return self._route_table['route_table']

    def from_dict(self, route_table):
        """ Load route table config from dict

        :param route_table: The underlying route table
        :type route_table: dict
        :rtype: RouteTable

        """
        self._route_table = route_table
        return self

    def to_dict(self) -> dict:
        """ Return underlying route table

        :rtype: dict

        """
        return self._route_table


class QueryCondition():
    """QueryCondition is context for job query"""

    def __init__(self, job_id):
        """ Init QueryCondition with job id

        :param job_id: The uuid of job
        :type job_id: string

        """
        self._job_id = job_id

    def get_job_id(self):
        """ Fetch the job id

        :rtype: dict

        """
        return {'job_id': self._job_id}

    def set_job_id(self, job_id):
        """ Set job id

        :param job_id: The uuid of job
        :type job_id: string

        """
        self._job_id = job_id

    def __str__(self):

        return self.get_job_id().__str__()

    def to_dict(self):
        return {'job_id': self._job_id}


class Component():
    # TODO: add setter/getter
    """Component is used to describe steps in a pipline"""

    def __init__(self, name='', module='', need_deploy=True):
        """ Init an empty component
        """
        self._name = name
        self._module = module

        #: Need deploy of data io
        self._need_deploy = need_deploy

        #: The input part contains two sub structures, for more details please refer to `DSL definition <https://github.com/FederatedAI/FATE/blob/master/doc/dsl_conf_setting_guide.rst>`_
        #: They should be list type
        self._input_data = []
        self._input_train_data = []
        self._input_eval_data = []
        self._input_model = []
        self._input_isometric_model = []

        #: The output part also contains two sub structures, for more details please refer to `DSL definition <https://github.com/FederatedAI/FATE/blob/master/doc/dsl_conf_setting_guide.rst>`_
        #: They should be list type
        self._output_data = []
        self._output_model = []

    def to_dict(self):
        """ Convert Component to dictionary

        :rtype: dict

        """
        name = self._name
        body = {}
        module = {'module': self._module}
        inputs = {'input': {}}
        outputs = {'output': {}}
        need_deploy = {'need_deploy': self._need_deploy}

        # check all input
        if len(self._input_data) != 0:
            inputs['input']['data'] = {
                'data': self._input_data
            }
        elif len(self._input_train_data) != 0:
            inputs['input']['data'] = {
                'train_data': self._input_train_data
            }
        elif len(self._input_eval_data) != 0:
            inputs['input']['data'] = {
                'eval_data': self._input_eval_data
            }

        if len(self._input_model) != 0:
            inputs['input']['model'] = self._input_model
        elif len(self._input_isometric_model) != 0:
            inputs['input']['isometric_model'] = self._input_isometric_model

        if len(self._output_data) != 0:
            outputs['output']['data'] = self._output_data
        if len(self._output_model) != 0:
            outputs['output']['model'] = self._output_model

        body.update(module)
        if inputs != {'input': {}}:
            body.update(inputs)
        if outputs != {'output': {}}:
            body.update(outputs)
        body.update(need_deploy)

        return {name: body}


class ComponentBuilder():
    """ComponentBuilder is used to build Component instance"""

    def __init__(self, name='', module='', need_deploy=True):
        """ Init ComponentBuilder instance
        """

        self.reset(name, module, need_deploy)

    def reset(self, name, module, need_deploy):
        self._component = Component(name, module, need_deploy)

    def with_need_deploy(self, need_deploy):
        """ Set 'need_deploy' for DataIO module

        :param need_deploy: Value of need_deploy
        :type need_deploy: bool

        """
        self._component._need_deploy = need_deploy
        return self

    def with_name(self, name):
        """ Set component name

        :param name: name of the component
        :type name: string

        """
        self._component._name = name
        return self

    def with_module(self, module):
        """ Set component module

        :param module: The available module in FATE, for more details please refer to `Modules <https://github.com/FederatedAI/FATE/tree/master/federatedml>`
        :type module: string
        """
        self._component._module = module
        return self

    def add_input_data(self, *data):
        """ Set input data
        """
        self._component._input_data.extend([d for d in data])
        return self

    def add_input_train_data(self, *train_data):
        """ Set input data for training
        """
        self._component._input_train_data.extend([d for d in train_data])
        return self

    def add_input_eval_data(self, *eval_data):
        """ Set input data for evaluation
        """
        self._component._input_eval_data.extend([d for d in eval_data])
        return self

    def add_input_model(self, *model):
        """ Set input model
        """
        self._component._input_model.extend([d for d in module])
        return self

    def add_input_isometric_model(self, *isometric_model):
        """ Set input isometric model
        """
        self._component._input_isometric_model.extend(
            [d for d in isometric_model])
        return self

    def add_output_data(self, *data):
        """ Set output data
        """
        self._component._output_data.extend([d for d in data])
        return self

    def add_output_model(self, *model):
        """ Set output model
        """

        self._component._output_model.extend([d for d in model])
        return self

    def build(self):
        component = self._component
        self.reset('', '', True)
        return component


class Pipeline():
    """Pipline is used to described pipline in FATE"""

    def __init__(self, *components):
        self._name = 'components'
        self._components = {}

        for component in components:
            self._components.update(component.to_dict())

    def to_dict(self):
        return {self._name: self._components}


class PipelineBuilder():
    """PiplineBuilder is used to build pipline instance"""

    def __init__(self):
        self.reset()

    def reset(self):
        self._pipline = Pipeline()

    def with_components(self, *components):
        for component in components:
            self._pipline._components.update(component.to_dict())
        return self

    def build(self):
        return self._pipline


class Initiator():
    """Define initiator in configuration of job"""

    def __init__(self, role=None, party_id=None):
        self._role = role
        self._party_id = party_id

    def to_dict(self) -> dict:
        """ Return dictionary

        :rtype: dict

        """
        name = 'initiator'
        body = {
            'role': self._role,
            'party_id': self._party_id
        }

        return {name: body}


class InitiatorBuilder(object):
    """Build Initiator instance"""

    def __init__(self):
        self.reset()

    def reset(self):
        self._initiator = Initiator()

    def with_role(self, role: str) -> object:
        """

        :param role:
        :return:
        """
        self._initiator._role = role
        return self

    def with_party_id(self, part_id: int) -> object:
        """

        :param part_id:
        :return:
        """
        self._initiator._party_id = part_id
        return self

    def build(self) -> object:
        """

        :return:
        """
        initiator = self._initiator
        self.reset()
        return initiator


class JobParameters():
    """Define job parameters in configuration of job"""

    def __init__(self, work_mode=1, job_type=None, model_id=None, model_version=None):
        self._work_mode = 1
        self._job_type = job_type
        self._model_id = model_id
        self._model_version = model_version

    def to_dict(self):
        """ Return dictionary

        :rtype: dict

        """

        name = 'job_parameters'
        body = {
            'work_mode': self._work_mode
        }

        if self._job_type is not None:
            body['job_type'] = self._job_type

        if self._model_id is not None:
            body['model_id'] = self._model_id

        if self._model_version is not None:
            body['model_version'] = self._model_version

        return {name: body}


class JobParametersBuilder(object):
    """Build JobParameters instance"""

    def __init__(self):
        self.reset()

    def reset(self):
        self._job_parameters = JobParameters()

    def with_work_mode(self, work_mode: int) -> object:
        """

        :param work_mode:
        :type work_mode: int
        :return:

        """
        self._job_parameters._work_mode = work_mode
        return self

    def with_job_type(self, job_type: str) -> object:
        """

        :param job_type:
        :type job_type: string
        :return:

        """

        self._job_parameters._job_type = job_type
        return self

    def with_model_id(self, model_id: str) -> object:
        """

        :param model_id: id of model

        """
        self._job_parameters._model_id = model_id
        return self

    def with_model_version(self, model_version: str) -> object:
        """
        :param model_version: version of model

        """
        self._job_parameters._model_version = model_version
        return self

    def build(self) -> dict:
        """

        :return:
        """
        job_parameters = self._job_parameters
        self.reset()
        return job_parameters


class Role():
    """Define role in configuration of job"""

    def __init__(self):
        self._guest = []
        self._host = []
        self._arbiter = []

    def to_dict(self):
        name = 'role'
        body = {
            'guest': self._guest,
            'host': self._host,
            'arbiter': self._arbiter
        }

        return {name: body}


class RoleBuilder():
    """Build Role instance"""

    def __init__(self):
        self.reset()

    def reset(self):
        self._role = Role()

    def add_guest(self, party_id=''):
        self._role._guest.append(party_id)
        return self

    def add_host(self, party_id=''):
        self._role._host.append(party_id)
        return self

    def add_arbiter(self, party_id=''):
        self._role._arbiter.append(party_id)
        return self

    def with_guests(self, *guests) -> object:
        """

        :param guest: A list of guests
        :return:
        """

        self._role._guest.extend([d for d in guests])
        return self

    def with_hosts(self, *hosts) -> object:
        """

        :param host: A list of hosts
        :return:
        """
        self._role._host.extend([d for d in hosts])
        return self

    def with_arbiters(self, *arbiters) -> object:
        """

        :param arbiter: A list of arbiters
        :return:
        """
        self._role._arbiter.extend([d for d in arbiters])
        return self

    def build(self) -> object:
        """

        :return:
        """
        role = self._role
        self.reset()
        return role


class RoleParameters():
    """Define role_parameters in configuration of job"""

    def __init__(self):
        self._guest_data = []
        self._guest_module_config = []
        self._host_data = []
        self._host_module_config = []

    def to_dict(self):
        name = 'role_parameters'
        body = {
            'guest': {
                'args': {
                    'data': {}
                }
            },
            'host': {
                'args': {
                    'data': {}
                }
            }
        }

        for guest_data in self._guest_data:
            key = list(guest_data.keys())[0]
            value = guest_data[key]

            if body['guest']['args']['data'].get(key) != None:
                body['guest']['args']['data'][key].extend(value)
            else:
                body['guest']['args']['data'].update(guest_data)

        for host_data in self._host_data:
            key = list(host_data.keys())[0]
            value = host_data[key]

            if body['host']['args']['data'].get(key) != None:
                body['host']['args']['data'][key].extend(value)
            else:
                body['host']['args']['data'].update(host_data)

        for guest_module_config in self._guest_module_config:

            body['guest'].update(guest_module_config)

        for host_module_config in self._host_module_config:

            body['host'].update(host_module_config)

        if body['guest']['args']['data'] == {}:
            body['guest'].pop('args')

        if body['host']['args']['data'] == {}:
            body['host'].pop('args')

        return {name: body}


class RoleParametersBuilder(object):
    """Build RoleParameters instance"""

    def __init__(self):
        self.reset()

    def reset(self):
        self._role_parameters = RoleParameters()

    def _set_data(self, namespace='', name='', data_type='train_data', role='guest'):

        # simple check
        # if len(namespaces) != len(names):
        #    raise Exception("The number of namespaces and name is not matched")

        key = data_type
        body = [{'namespace': namespace,
                 'name': name}]

        # for ns, n in zip(namespaces, names):
        # body.append({'namespace': namespace, 'name': name})

        if role == 'guest':
            self._role_parameters._guest_data.append({key: body})

        elif role == 'host':
            self._role_parameters._host_data.append({key: body})

    def add_guest_train_data(self, namespace='', name='') -> object:
        """

        :param namespaces: The namespace of train data
        :type namespaces: string
        :param name: The name of train data
        :type name: string

        :return:
        """
        self._set_data(namespace, name, 'train_data', 'guest')

        return self

    def add_guest_eval_data(self, namespace='', name='') -> object:
        """

        :param namespaces: The namespace of train data
        :type namespaces: string
        :param name: The name of train data
        :type name: string

        :return:
        """
        self._set_data(namespace, name, 'eval_data', 'guest')

        return self

    def add_host_train_data(self, namespace='', name='') -> object:
        """

        :param namespaces: The namespace of train data
        :type namespaces: string
        :param name: The name of train data
        :type name: string

        :return:
        """
        self._set_data(namespace, name, 'train_data', 'host')

        return self

    def add_host_eval_data(self, namespace='', name='') -> object:
        """

        :param namespaces: The namespace of train data
        :type namespaces: string
        :param name: The name of train data
        :type name: string

        :return:
        """
        self._set_data(namespace, name, 'eval_data', 'host')

        return self

    def _set_config(self, module='', config={}, role='guest'):
        # simple check
        # if len(modules) != len(configs):
        #    raise Exception("The number of modules and configs in not matched")

        # for m, c in zip(modules, configs):
        key = module
        body = config

        if role == 'guest':
            self._role_parameters._guest_module_config.append({key: body})
        elif role == 'host':
            self._role_parameters._host_module_config.append({key: body})

    def add_guest_module_config(self, module='', config={}):
        """ Set guest module config

        :param modules: The modules
        :type modules: str
        :param configs: The configs
        :type configs: dict
        """

        self._set_config(module, config, 'guest')
        return self

    def add_host_module_config(self, module='', config={}):
        """ Set guest module config

        :param modules: The modules
        :type modules: str
        :param configs: The configs
        :type configs: dict
        """

        self._set_config(module, config, 'host')
        return self

    def build(self) -> object:
        """

        :return:
        """

        role_parameters = self._role_parameters
        self.reset()
        return role_parameters


class AlgorithmParameters():
    """Class to define algorithm_parameters in configuration of job"""

    def __init__(self):
        self._parameters_list = []

    def to_dict(self):
        name = 'algorithm_parameters'
        body = {}

        for parameters in self._parameters_list:
            body.update(parameters)

        return {name: body}


class AlgorithmParametersBuilder():
    """Class to build AlgorithmParameters instance"""

    def __init__(self):
        self.reset()

    def reset(self):
        self._algorithm_parameters = AlgorithmParameters()

    def add_module_config(self, module='', config={}):
        """

        :param modules: The available module in FATE, for more details please refer to `Modules <https://github.com/FederatedAI/FATE/tree/master/federatedml>`
        :type modules: list

        :param configs: A list of configs, the available module config in FATE, for more details please refer to `parameters <https://github.com/FederatedAI/FATE/blob/master/federatedml/conf/default_runtime_conf>`

        :type confgis: list
        """
        # if len(modules) != len(configs):
        #    raise Exception("The number of modules and configs is not matched")

        # for m, c in zip(modules, configs):
        #    self._algorithm_parameters._parameters_list.append({m: c})
        self._algorithm_parameters._parameters_list.append({module: config})

        return self

    def build(self):
        algorithm_parameters = self._algorithm_parameters
        self.reset()
        return algorithm_parameters


class Config(object):
    """Define configuration of job"""

    def __init__(self, initiator=None, job_parameters=None,
                 role=None, role_parameters=None, algorithm_parameters=None):
        self._initiator = initiator
        self._job_parameters = job_parameters
        self._role = role
        self._role_parameters = role_parameters
        self._algorithm_parameters = algorithm_parameters

    def to_dict(self) -> dict:

        result = {}

        if self._initiator is not None:
            result.update(self._initiator.to_dict())

        if self._job_parameters is not None:
            result.update(self._job_parameters.to_dict())

        if self._role is not None:
            result.update(self._role.to_dict())

        if self._role_parameters is not None:
            result.update(self._role_parameters.to_dict())

        if self._algorithm_parameters is not None:
            result.update(self._algorithm_parameters.to_dict())

        return result


class ConfigBuilder(object):
    """Class to build Config instance"""

    def __init__(self):
        self.reset()

    def reset(self):
        self._config = Config()

    def with_initiator(self, initiator: Initiator) -> object:
        self._config._initiator = initiator
        return self

    def with_job_parameters(self, job_parameters: JobParameters) -> object:
        self._config._job_parameters = job_parameters
        return self

    def with_role(self, role: Role) -> object:
        self._config._role = role
        return self

    def with_role_parameters(self, role_parameters: RoleParameters) -> object:
        self._config._role_parameters = role_parameters
        return self

    def with_algorithm_parameters(self, algorithm_parameters: AlgorithmParameters) -> object:
        self._config._algorithm_parameters = algorithm_parameters
        return self

    def build(self) -> object:
        config = self._config
        self.reset()
        return config
