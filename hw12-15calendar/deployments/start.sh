#!/bin/sh

set -e

if [ "$ENV" = 'DEV' ]; then
	# Есть возможность изменять код сервера без пересборки образа docker
	echo "Running Development Server"
	cd services/calendar/
	exec make run
elif [ "$ENV" = 'UNIT' ]; then
	echo "Running Unit Tests"
	make test
	cd services/calendar/
	exec make test
else
	# Запуск сервера собраного при сборке образа docker
	echo "Running Production Server"
	exec ./bin/calendar -config "configs/calendar.yml"
fi
