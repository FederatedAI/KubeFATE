# Overview

This is a demo for jenkins CI workflow

## Usage

Assume you have built jenkins image called `jenkin:test`, use the following command to start jenkins service:

```bash
docker run -d -p 8888:8080 --name jenkins -v /home/luke/document/github/jenkins/data:/var/jenkins_home/  -u 1000 jenkins:test
```

## Configuration

### Plugins

1. Install ansible plugin
2. Set global ansible tool with path `/usr/local/bin`

### Credential

1. Add SSH key as jenkins credential ssh-key with ID `ssh-ansible`
2. Add Inventory file as jenkins credential file with ID `ansible-inventory`
