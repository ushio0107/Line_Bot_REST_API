MONGO_CONTAINER = mongo4.4
MONGO_INITDB_ROOT_USERNAME = root
MONGO_INITDB_ROOT_PASSWORD = 123456
DB_PATH = ${HOME}

docker_mongo_image:
	docker pull mongo:4.4

docker_install:
	docker run --name ${MONGO_CONTAINER} \
	-e MONGO_INITDB_ROOT_USERNAME=${MONGO_INITDB_ROOT_USERNAME} \
	-e MONGO_INITDB_ROOT_PASSWORD=${MONGO_INITDB_ROOT_PASSWORD} \
	-v ${DB_PATH}/data:/data/db \
	-p 27017:27017 \
	-d mongo:4.4 \

docker_start:
	docker start ${MONGO_CONTAINER}

docker_rm:
	docker stop ${MONGO_CONTAINER} && docker rm ${MONGO_CONTAINER}

echo:
	echo ${DB_PATH}/data
