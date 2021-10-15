pkgDir=${PWD}/pkg/...
internalDir=${PWD}/internal/...
outputDir=${PWD}/_dist
testOutputDir=${outputDir}/tests
envsubstImage=bhgedigital/envsubst # or cmd.cat/envsubst
grafanaInfluxDBSource=./docker/grafana/provisioning/datasources/datasource.yaml

# Development
up: create-dev-env create-dev-config
	@docker compose up --build -d

down:
	@docker compose down -v

create-dev-env:
	@test -e .env || cp .env.example .env

create-dev-config:
	@docker run --rm \
		--env-file .env \
		-v `pwd`:`pwd` -w `pwd` \
		-i $(envsubstImage) \
		envsubst '$$DOCKER_INFLUXDB_INIT_ORG,$$DOCKER_INFLUXDB_INIT_BUCKET,$$DOCKER_INFLUXDB_INIT_ADMIN_TOKEN' < $(grafanaInfluxDBSource).tmpl > $(grafanaInfluxDBSource)

# CI
ci-test:
	@docker compose exec golang sh -c 'make cli-test'

# DEMO
monitor:
	@docker compose exec golang go run cmd/endpoint-monitor.go status -f $(config)

# CLI
cli-test:
	@make cli-lint
	@make cli-unit

cli-lint:
	@golangci-lint run ${pkgDir} ${internalDir} -v

cli-unit:
	@mkdir -p ${testOutputDir}
	@go clean -testcache
	@go test \
        -cover \
        -coverprofile=cp.out \
        -outputdir=${testOutputDir} \
        -race \
        -v \
        -failfast \
        ${pkgDir} \
        ${internalDir}
	@go tool cover -html=${testOutputDir}/cp.out -o ${testOutputDir}/cp.html