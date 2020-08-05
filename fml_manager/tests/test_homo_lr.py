
from fml_manager import *

# dsl
data_io = ComponentBuilder()\
    .with_name('dataio_0')\
    .with_module('DataIO')\
    .add_input_data('args.train_data')\
    .add_output_data('train')\
    .add_output_model('dataio').build()


homo_lr = ComponentBuilder()\
    .with_name('homo_lr_0')\
    .with_module('HomoLR')\
    .add_input_train_data('dataio_0.train')\
    .add_output_data('train')\
    .add_output_model('homolr').build()

evaluation = ComponentBuilder()\
    .with_name('evaluation_0')\
    .with_module('Evaluation')\
    .add_input_data('homo_lr_0.train')\
    .add_output_data('evaluate').build()

pipeline = Pipeline(
    data_io,
    homo_lr,
    evaluation
)

# Configuration
initiator = Initiator(role='guest', party_id=10000)

job_parameters = JobParametersBuilder()\
    .with_work_mode(1).build()

role = RoleBuilder()\
    .add_guest(party_id=10000)\
    .add_host(party_id=10000)\
    .add_host(party_id=9999)\
    .add_arbiter(party_id=10000).build()

eval_config = {
    'need_run': [False]
}

role_parameters = RoleParametersBuilder()\
    .add_guest_train_data(namespace='homo_breast_guest', name='homo_breast_guest')\
    .add_host_train_data(namespace='homo_breast_host', name='homo_breast_host_a')\
    .add_host_train_data(namespace='homo_breast_host', name='homo_breast_host_b')\
    .add_host_module_config(module='evaluation_0', config=eval_config).build()


homo_lr_params = {
    'penalty': 'L2',
    'optimizer': 'sgd',
    'eps': 1e-5,
    'alpha': 0.01,
    'max_iter': 10,
    'converge_func': 'diff',
    'batch_size': 500,
    'learning_rate': 0.15,
    'decay': 1,
    'decay_sqrt': True,
    'init_param': {
        'init_method': 'zeros'
    },
    'encrypt_param': {
        'method': 'Paillier'
    },
    'cv_param': {
        'n_splits': 4,
        'shuffle': True,
        'random_seed': 33,
        'need_cv': False
    }
}
dotaio_config = {
    'with_label': True,
    'label_name': 'y',
    'label_type': 'int',
    'output_format': 'dense'
}

algorithm_parameters = AlgorithmParametersBuilder()\
    .add_module_config(module='homo_lr_0', config=homo_lr_params)\
    .add_module_config(module='dataio_0', config=dotaio_config).build()

config = Config(
    initiator,
    job_parameters,
    role,
    role_parameters,
    algorithm_parameters
)

print(pipeline.to_dict())
print(config.to_dict())
