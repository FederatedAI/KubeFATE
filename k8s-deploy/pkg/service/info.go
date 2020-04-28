package service

func GetClusterInfo(name, namespace string) (map[string]interface{}, error) {
	ip, err := GetNodeIp()
	if err != nil {
		return nil, err
	}
	port, err := GetProxySvcNodePorts(name, namespace)
	if err != nil {
		return nil, err
	}
	podList, err := GetPodList(name, namespace)
	if err != nil {
		return nil, err
	}

	ingressUrlList, err := GetIngressUrl(name, namespace)
	if err != nil {
		return nil, err
	}

	info := make(map[string]interface{})

	if len(ip) > 0 {
		info["ip"] = ip[len(ip)-1]
	}
	if len(port) > 0 {
		info["port"] = port[0]
	}
	info["modules"] = podList

	info["dashboard"] = ingressUrlList

	return info, nil
}
