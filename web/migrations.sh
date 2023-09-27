#!/bin/sh

# set -e

cmd="$1"

migrate_tool="migrate"
migrate_source="file:///migrations"

m_mysql_address=${MYSQL_ADDRESS:-127.0.0.1}
m_mysql_username=${MYSQL_USERNAME:-root}
m_mysql_password=${MYSQL_PASSWORD:-root}
m_mysql_database=${MYSQL_DATABASE:-demo}

migrate_mysql()
{
    m_mysql_url="mysql://$m_mysql_username:$m_mysql_password@tcp($m_mysql_address)/$m_mysql_database"

    echo "migrating for mysql, address: $m_mysql_address, keyspace: $m_mysql_database"

    "$migrate_tool" -source=$migrate_source -database="$m_mysql_url" up || exit 1
}

case "$cmd" in
    migrate) echo "Running migrate..."
        migrate_mysql
        ;;
    *) echo "Unknown command: $cmd"
       echo "Usage: $0 migrate"
       exit 1
        ;;
esac
