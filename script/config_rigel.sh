#!/bin/sh
rigelctl -e localhost:2379 -a CVL-KRA -m KRA -v 1 schema add ../config/config_rigel.json
rigelctl -a CVL-KRA -m KRA -v 1 --config dev config set app_server_port "8080"
rigelctl -a CVL-KRA -m KRA -v 1 --config dev config set db_host "localhost"
rigelctl -a CVL-KRA -m KRA -v 1 --config dev config set db_port "5432"
rigelctl -a CVL-KRA -m KRA -v 1 --config dev config set db_user "postgres"
rigelctl -a CVL-KRA -m KRA -v 1 --config dev config set db_password "postgres"
rigelctl -a CVL-KRA -m KRA -v 1 --config dev config set db_name "cvl_kra_db"
rigelctl -a CVL-KRA -m KRA -v 1 --config dev config set "logger_priority" "Debug2"
rigelctl -a CVL-KRA -m KRA -v 1 --config dev config set "CALC.OUT.PRECISION" "3"
rigelctl -a CVL-KRA -m KRA -v 1 --config dev config set "mockoon_panvalid" "http://localhost:3000/panvalidation"
rigelctl -a CVL-KRA -m KRA -v 1 --config dev config set "mockoon_getpassword" "http://localhost:3000/getpassword"
rigelctl -a CVL-KRA -m KRA -v 1 --config dev config set "mockoon_pandownload_I" "http://localhost:3000/individualpandownload"
rigelctl -a CVL-KRA -m KRA -v 1 --config dev config set "mockoon_pandownload_NI" "http://localhost:3000/nonindividualpandownload"
rigelctl -a CVL-KRA -m KRA -v 1 --config dev config set "mockoon_pandataupload" "http://localhost:3000/pandataupload"

