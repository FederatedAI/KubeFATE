import pprint
import json

from fml_manager import *

# hetero lr
pipeline_str = '''
{
    "components" : {
        "dataio_0": {
            "module": "DataIO",
            "input": {
                "data": {
                    "data": [
                        "args.train_data"
                    ]
                }
            },
            "output": {
                "data": ["train"],
                "model": ["dataio"]
            },
                        "need_deploy": true
         },
        "hetero_feature_binning_0": {
            "module": "HeteroFeatureBinning",
            "input": {
                "data": {
                    "data": [
                        "dataio_0.train"
                    ]
                }
            },
            "output": {
                "data": ["train"],
                "model": ["hetero_feature_binning"]
            }
        },
        "hetero_feature_selection_0": {
            "module": "HeteroFeatureSelection",
            "input": {
                "data": {"data": [
                        "hetero_feature_binning_0.train"
                    ]
                },
                "isometric_model": [
                    "hetero_feature_binning_0.hetero_feature_binning"
                ]
            },
            "output": {
                "data": ["train"],
                "model": ["selected"]
            }
        },
        "hetero_lr_0": {
            "module": "HeteroLR",
            "input": {
                "data": {
                    "train_data": ["hetero_feature_selection_0.train"]
                }
            },
            "output": {
                "data": ["train"],
                "model": ["hetero_lr"]
            }
        },
        "evaluation_0": {
            "module": "Evaluation",
            "input": {
                "data": {
                    "data": ["hetero_lr_0.train"]
                }
            },
            "output": {
                "data": ["evaluate"]
            }
        }
    }
}
'''

# Pipeline
data_io = ComponentBuilder(name='dataio_0',
                           module='DataIO')\
                           .add_input_data('args.train_data')\
                           .add_output_data('train')\
                           .add_output_model('dataio').build()
        
hetero_feature_binning = ComponentBuilder(name='hetero_feature_binning_0',
                                          module='HeteroFeatureBinning')\
                                          .add_input_data('dataio_0.train')\
                                          .add_output_data('train')\
                                          .add_output_model('hetero_feature_binning').build()

hetero_feature_selection = ComponentBuilder(name='hetero_feature_selection_0',
                                            module='HeteroFeatureSelection')\
                                            .add_input_data('hetero_feature_binning_0.train')\
                                            .add_input_isometric_model('hetero_feature_binning_0.hetero_feature_binning')\
                                            .add_output_data('train')\
                                            .add_output_model('selected').build()

hetero_lr = ComponentBuilder(name='hetero_lr_0',
                             module='HeteroLR')\
                             .add_input_train_data('hetero_feature_selection_0.train')\
                             .add_output_data('train')\
                             .add_output_model('hetero_lr').build()

evaluation = ComponentBuilder(name='evaluation_0',
                              module='Evaluation',
                              need_deploy=False)\
                              .add_input_data('hetero_lr_0.train')\
                              .add_output_data('evaluate').build()
pipeline = Pipeline(
    data_io, 
    hetero_feature_selection,  
    hetero_feature_binning, 
    hetero_lr, 
    evaluation)
    
lho = pipeline.to_dict()
rho = json.loads(pipeline_str)

pprint.pprint(lho)
print('------')
pprint.pprint(rho)

# config
config_str = '''
{
    "initiator": {
        "role": "guest",
        "party_id": 10000
    },
    "job_parameters": {
        "work_mode": 1
    },
    "role": {
        "guest": [10000],
        "host": [9999],
        "arbiter": [9999]
    },
    "role_parameters": {
        "guest": {
            "args": {
                "data": {
                    "train_data": [{"name": "breast_b", "namespace": "fate_flow_test_breast"}]
                }
            },
            "dataio_0":{
                "with_label": [true],
                "label_name": ["y"],
                "label_type": ["int"],
                "output_format": ["dense"]
            }
        },
        "host": {
            "args": {
                "data": {
                    "train_data": [{"name": "breast_a", "namespace": "fate_flow_test_breast"}]
                }
            },
             "dataio_0":{
                "with_label": [false],
                "output_format": ["dense"]
            }
        }
    },
    "algorithm_parameters": {
        "hetero_lr_0": {
            "penalty": "L2",
            "optimizer": "rmsprop",
            "eps": 1e-5,
            "alpha": 0.01,
            "max_iter": 3,
            "converge_func": "diff",
            "batch_size": 320,
            "learning_rate": 0.15,
            "init_param": {
                                "init_method": "random_uniform"
            }
        }
    }
}
'''

# Configuration
initiator = Initiator(role='guest', party_id=10000)

job_parameters = JobParameters(work_mode=1)

role = RoleBuilder()\
    .add_guest(10000)\
    .add_host(9999)\
    .add_arbiter(9999).build()

guest_data_io_config = {
    'with_label': [True],
    'label_name': ['y'],
    'label_type': ['int'],
    'output_format': ['dense']
}

host_data_io_config = {
    'with_label': [False],
    'output_format': ['dense']
}

role_parameters = RoleParametersBuilder()\
    .add_guest_train_data(namespace='fate_flow_test_breast', name='breast_b')\
    .add_guest_module_config(module='dataio_0', config=guest_data_io_config)\
    .add_host_train_data(namespace='fate_flow_test_breast', name='breast_a')\
    .add_host_module_config(module='dataio_0', config=host_data_io_config).build()

hetero_lr_params = {
    "penalty": "L2",
    "optimizer": "rmsprop",
    "eps": 1e-5,
    "alpha": 0.01,
    "max_iter": 3,
    "converge_func": "diff",
    "batch_size": 320,
    "learning_rate": 0.15,
    "init_param": {
        "init_method": "random_uniform"
    }
}

algorithm_parameters = AlgorithmParametersBuilder()\
    .add_module_config(module='hetero_lr_0', config=hetero_lr_params).build()

config = Config(initiator, job_parameters, role,
                role_parameters, algorithm_parameters)

lho = config.to_dict()
rho = json.loads(config_str)

pprint.pprint(lho)
print('------')
pprint.pprint(rho)
