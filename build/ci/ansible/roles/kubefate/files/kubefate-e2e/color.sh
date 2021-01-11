#!/bin/bash
INFO="\033[36mInfo:\033[0m"
Warning="\033[33mWarning:\033[0m"
ERROR="\033[31mErrot:\033[0m"
SUCCESS="\033[32mSuccess:\033[0m"
DEBUG="\033[34mDebug:\033[0m"

log() {
    echo $@
}

loginfo() {
    echo -e $INFO $@
}

logwarning() {
    echo -e $Warning $@
}

logerror() {
    echo -e $ERROR $@
}

logsuccess() {
    echo -e $Success $@
}

logdebug() {
    echo -e $DEBUG $@
}
