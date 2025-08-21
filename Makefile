exec:
	go build -o go-dsp && ./go-dsp

dockerc:
	docker build -t go-dsp-image .
	docker stop go-dsp
	docker rm go-dsp
	docker container create -p 8080:8080 --name go-dsp go-dsp-image

dockerup:
	docker start go-dsp

dockerdown:
	docker stop go-dsp

gu:
	git add *
	git commit -m "autocommit"
	git push

gd:
	git pull
	make exec
