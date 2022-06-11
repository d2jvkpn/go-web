#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})


app_env=$1

# ~/.ansible.cfg
# [defaults]
# inventory = ~/.ansible/hosts.ini
# log_path = $PWD/ansible.log

# --become: root user
ansible-playbook -vv --inventory=hosts.ini \
  "${_path}/playbook.yaml" --extra-vars="app_env=$app_env"
