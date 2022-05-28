set -eu;

FDB_CLUSTER_FILE=${FDB_CLUSTER_FILE:-/etc/foundationdb/fdb.cluster}
mkdir -p "$(dirname $FDB_CLUSTER_FILE)"

if [[ -n $FDB_COORDINATOR ]]; then
    coordinator_ip=$(dig +short "$FDB_COORDINATOR")
    if [[ -z "$coordinator_ip" ]]; then
        echo "Failed to look up coordinator address for $FDB_COORDINATOR" 1>&2
        exit 1
    fi
    coordinator_port=${FDB_COORDINATOR_PORT:-4500}
    echo "docker:docker@$coordinator_ip:$coordinator_port" > "$FDB_CLUSTER_FILE"
else
    echo "FDB_COORDINATOR environment variable not defined" 1>&2
    exit 1
fi

if ! /usr/bin/fdbcli -C $FDB_CLUSTER_FILE --exec status --timeout 3 ; then
    echo "creating the database"
    if ! fdbcli -C $FDB_CLUSTER_FILE --exec "configure new single memory ; status" --timeout 10 ; then
        echo "Unable to configure new FDB cluster."
        exit 1
    fi
fi

cd gen && make test