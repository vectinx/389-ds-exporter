#!/bin/bash

LDAP_URI="ldap://localhost:3389"
BIND_DN="cn=Directory Manager"
BIND_PW="$DS_DM_PASSWORD"
# SUFFIX_NAME="$DS_SUFFIX_NAME"
MAX_ATTEMPTS=24
SLEEP_INTERVAL=5


function log() {
    echo "INITIALIZATION_SCRIPT - $@"
}

function terminateDs() {
    log "Kill 389-ds server ..."
    kill -2 $(pgrep -f dscontainer)
    exit 1
}

function waitForStart()
{
    log "Waiting for LDAP server start ..."

    for ((i=1; i<=3; i++)); do
        for ((j=1; j<=MAX_ATTEMPTS; j++)); do
            /usr/sbin/dsconf localhost backend suffix list
            if [ $? -eq 0 ]; then
                log "Successfuly connected to LDAP on $j attempt. Increment success count"
                SUCCESS_CONNECT=1
                break
            else
                log "Attempt $j/$MAX_ATTEMPTS: LDAP server is not available yet. Repeat in $SLEEP_INTERVAL seconds"
                sleep $SLEEP_INTERVAL
            fi
        done
    done
    if [ $SUCCESS_CONNECT -ne 1 ]; then
        log "Failed to connect to LDAP server after $MAX_ATTEMPTS attempts."
        exit 1
    else
        return 0
    fi
}

log "Start 389-ds instance"
/usr/lib/dirsrv/dscontainer -r &

zypper --non-interactive install openldap2-client

log "Wait for instance start"

waitForStart

log "Check if suffix already exists"

dsconf localhost backend suffix list | grep -w "$SUFFIX_NAME"

if [ $? -eq 0 ]; then
    log "Suffix $SUFFIX_NAME already exists. Skipping creation ..."
else
    sleep 5
    /usr/sbin/dsconf localhost backend create --be-name userRoot --suffix "$SUFFIX_NAME"
    if [ $? -eq 0 ]; then
        log "Successfuly created suffix"
    else
        log "Error creating suffix"
        terminateDs
    fi
fi

waitForStart

ldapadd -x -H "$LDAP_URI" -D "$BIND_DN" -w "$BIND_PW" <<EOF
dn: dc=example,dc=com
objectClass: top
objectClass: dcObject
objectClass: organization
o: Example Organization
dc: example

dn: dc=test,dc=example,dc=com
objectClass: top
objectClass: dcObject
objectClass: organization
o: Test Organization
dc: test
EOF

tail -f /dev/null