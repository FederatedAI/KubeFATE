import pprint
import json

from fml_manager import *

# dsl
secure_add_example = Component(name='secure_add_example_0',
                               module='SecureAddExample')

dsl = Pipeline(
    secure_add_example
)


# Configuration
initiator = Initiator(role='guest',
                      party_id=9999)


job_parameters = JobParameters(work_mode=1)

role = RoleBuilder()\
    .add_guest(9999)\
    .add_host(9999).build()

secure_add_example_guest_config = {
    "seed": [
        123
    ]
}


secure_add_example_host_config = {
    "seed": [
        321
    ]
}
role_parameters = RoleParametersBuilder()\
    .add_host_module_config(module='secure_add_example_0', config=secure_add_example_host_config)\
    .add_guest_module_config(module='secure_add_example_0', config=secure_add_example_guest_config).build()

secure_add_example = {
    "partition": 10,
    "data_num": 1000
}


algorithm_parameters = AlgorithmParametersBuilder()\
    .add_module_config(module='secure_add_example_0', config=secure_add_example).build()

config = Config(
    initiator,
    job_parameters,
    role,
    role_parameters,
    algorithm_parameters
)

pprint.pprint(config.to_dict())
pprint.pprint(dsl.to_dict())
