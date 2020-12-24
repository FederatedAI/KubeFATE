# Overview

The jenkins CI workflow

## Prerequisites

A linux host with jenkins and ansible installed.

## Configuration

### Plugins

1. Install ansible plugin
2. Set global ansible tool path, for example `/usr/local/bin`

### Credential

1. Add SSH key as jenkins credential ssh-key with ID `ssh-ansible`
2. Add Inventory file as jenkins credential file with ID `ansible-inventory`