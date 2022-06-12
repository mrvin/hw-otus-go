#!/bin/bash

set -e

if [ "$ENV" = 'DEV' ]; then
	# Есть возможность изменять код сервера без пересборки образа docker
	echo "Running Development Server"
	exec make run
elif [ "$ENV" = 'UNIT' ]; then
	echo "Running Unit Tests"
	exec make test
else
	# Запуск сервера собраного при сборке образа docker
	echo "Running Production Server"
	exec ./bin/calendar -config "configs/config.yml"
fi
