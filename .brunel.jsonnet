/*local shared = brunel.shared({
    repository: "/Users/lewis/Documents/Projects/test_repo",
    branch: "master",
    file: "file.jsonnet"
});*/

{
    version: "v1",
    description: "This is my supper cool DSL",
    stages: {
        test: {
            services: [
                {
                    image: "golang:1.12.0",
                    wait: {
                        output: "3",
                        timeout: 10
                    },
                    entrypoint: "sh",
                    args: [ "-c", "--", "for i in $(seq 1 30); do echo $i; sleep 1; done" ]
                },
                {
                    image: "nginx:latest",
                    hostname: "nginx"
                }
            ],
            steps: [
                {
                    image: "golang:1.12.0",
                    resources: {
                        limits: {
                            CPU: 0.5,
                            memory: "200m"
                        }
                    },
                    working_dir: "/workspace/src/go-brunel",
                    entrypoint: "sh",
                    args: [ "-c", "--", "for i in $(seq 1 10); do echo $i; sleep 1; done; echo 'sdasd' > /dev/stderr; cat asdasd; echo '\u001b[0m\u001b[4m\u001b[42m\u001b[31mfoo\u001b[39m\u001b[49m\u001b[24mfoo\u001b[0m';" ]
                },
                {
                    image: "byrnedo/alpine-curl",
                    entrypoint: "sh",
                    args: [ "-c", "--", "curl http://nginx" ]
                },
                {
                    image: "byrnedo/alpine-curl",
//                    environment: {
//                        GOPATH : brunel.environment('MY_CONFIG'),
//                        USERNAME : brunel.secret('MY_SECRET'),
//                    },
                    entrypoint: "sh",
                    args: [ "-c", "--",  "echo $GOPATH; echo $USERNAME; echo \"\u001b[30;1m A \u001b[31;1m B \u001b[32;1m C \u001b[33;1m D \u001b[0m\"" ]
                }
            ],
        },
        build: {
            services: [
                {
                    image: "nginx:latest",
                    hostname: "nginx"
                }
            ],
            steps: [
                {
                    image: "byrnedo/alpine-curl",
                    entrypoint: "sh",
                    args: [ "-c", "--", "curl http://nginx" ]
                },
            ],
        }
    }
}