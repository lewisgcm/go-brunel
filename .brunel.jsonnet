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
//            services: [
//                {
//                    image: "golang:1.12.0",
//                    wait: {
//                        output: "3",
//                        timeout: 10
//                    },
//                    entrypoint: "sh",
//                    args: [ "-c", "--", "for i in $(seq 1 30); do echo $i; sleep 1; done" ]
//                },
////                {
////                    image: "nginx:latest",
////                    hostname: "nginx"
////                }
//            ],
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
                    args: [ "-c", "--", "for i in $(seq 1 1); do echo $i; sleep 1; done" ]
                },
//                {
//                    image: "byrnedo/alpine-curl",
//                    entrypoint: "sh",
//                    args: [ "-c", "--", "curl http://nginx" ]
//                },
//                {
//                    image: "byrnedo/alpine-curl",
//                    environment: {
//                        GOPATH : brunel.environment('MY_CONFIG'),
//                        USERNAME : brunel.secret('MY_SECRET'),
//                        //HELLO: shared.hello,
//                    },
//                    entrypoint: "sh",
//                    args: [ "-c", "--",  "echo $GOPATH; echo $USERNAME; echo $HELLO" ]
//                }
            ],
        }
//        build: {
//            services: [
//                {
//                    image: "nginx:latest",
//                    hostname: "nginx"
//                }
//            ],
//            steps: [
//                {
//                    image: "byrnedo/alpine-curl",
//                    entrypoint: "sh",
//                    args: [ "-c", "--", "curl http://nginx" ]
//                },
//            ],
//        }
    }
}