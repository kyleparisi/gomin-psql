install:
	pip install -r ./schema/requirements.txt

run:
	bash run.sh

db:
	./schema/manage.py makemigrations
	./schema/manage.py migrate

kill:
	docker kill MyApp-mysql

test:
	bash test.sh
