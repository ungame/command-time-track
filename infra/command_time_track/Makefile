up:
	docker-compose up

down:
	- docker-compose down

clear: down
	- docker volume rm command_time_track_grafana_data
	- docker volume rm command_time_track_mysql_data
	- docker volume rm command_time_track_prometheus_data