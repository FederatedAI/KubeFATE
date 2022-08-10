# FATE Algorithm and Computational Acceleration Selection

FATE作为联邦学习框架，支持很多算法，根据业务选择对应的算法和加速卡也是很多企业的需求，当前KubeFATE支持选择算法和加速卡的选择

不论在docker-compose还是k8s部署中，都可以对以下两个参数做选择：

- `algorithm` 算法选择
- `device` 计算设备选择

## 算法

当前算法的选择包含了两个选项：

- `Basic`
    Basic是默认选项，包含了除去nn（包括homo_nn和hetero_nn）算法相关的依赖组件。
- `NN`
    NN就包含了nn包括（homo_nn和hetero_nn）所需的所有依赖。***仅当computing是Eggroll的时候才可以使用NN***

## 计算设备选择

当前设备的选择包含了一个选项：

- `CPU`
    CPU是FATE计算主要的使用设备。
