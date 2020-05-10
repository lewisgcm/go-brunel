{
    version: "v1",
    description: "Brunel CI/CD Build Pipeline",
    stages: [
        {
            name: "unit test",
            steps: [
                {
                    image: "golang:1.13.6-buster",
                    entryPoint: "sh",
                    workingDir: "/workspace/",
                    args: [
                        "-c",
                        "--",
                        "go test go-brunel/internal...",
                    ]
                }
            ]
        },
        {
            name: "data store integration test",
            services: [
                {
                    image: "mongo:4.2.1-bionic",
                    wait: {
                        output: "Listening on",
                        timeout: 30
                    },
                    environment: {
                        "MONGO_INITDB_DATABASE": "brunel",
                        "MONGO_INITDB_ROOT_USERNAME": "root",
                        "MONGO_INITDB_ROOT_PASSWORD": "example",
                    },
                    hostname: "mongo"
                },
            ],
            steps: [
                {
                    image: "golang:1.13.6-buster",
                    entryPoint: "sh",
                    workingDir: "/workspace/",
                    args: [
                        "-c",
                        "--",
                        'go test go-brunel/test/store... -mongo-db-uri="mongodb://root:example@mongo:27017"',
                    ]
                }
            ]
        },
        {
            name: "build",
            when: brunel.environment.variable("RELEASE") == "1",
            services: [
                {
                    image: "docker:dind",
                    privileged: true,
                    wait: {
                        output: "Daemon has completed initialization",
                        timeout: 30
                    },
                    environment: {
                        "DOCKER_TLS_CERTDIR": "",
                    },
                    hostname: "docker"
                },
            ],
            steps: [
                {
                    image: "docker:dind",
                    workingDir: "/workspace/",
                    entryPoint: "sh",
                    environment: {
                        "DOCKER_HOST": "tcp://docker:2375",
                    },
                    args: [
                        "-c",
                        "--",
                        "docker build -t lewisgcm/go-brunel:runner -f ./Dockerfile.runner .",
                    ]
                },
                {
                    image: "docker:dind",
                    workingDir: "/workspace/",
                    entryPoint: "sh",
                    environment: {
                        "DOCKER_HOST": "tcp://docker:2375",
                    },
                    args: [
                        "-c",
                        "--",
                        "docker build -t lewisgcm/go-brunel:latest -f ./Dockerfile .",
                    ]
                },
                {
                    image: "docker:dind",
                    workingDir: "/workspace/",
                    entryPoint: "sh",
                    environment: {
                        "DOCKER_HOST": "tcp://docker:2375",
                        "DOCKER_HUB_ACCESS_TOKEN": brunel.environment.variable('DOCKER_HUB_ACCESS_TOKEN')
                    },
                    args: [
                        "-c",
                        "--",
                        "docker login -u lewisgcm -p $DOCKER_HUB_ACCESS_TOKEN && docker push lewisgcm/go-brunel:runner && docker push lewisgcm/go-brunel:latest",
                    ]
                },
            ]
        },
    ]
}