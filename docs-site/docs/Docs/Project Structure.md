**The ArmurAI backend follows a modular design, with code organized into distinct directories and packages based on their functionalities. This approach promotes code reusability, maintainability, and separation of concerns.**

## Directory Structure

```bash
Project Directory Structure:
├── Dockerfile
├── LICENSE
├── Makefile
├── cmd
│   ├── server
│   │   └── main.go
├── docker-compose.yml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal
│   ├── api
│   │   ├── handlers.go
│   │   └── routes.go
│   ├── redis
│   │   └── redis.go
│   ├── tasks
│   │   ├── queue_task.go
│   │   ├── result_store.go
│   │   └── tasks.go
│   ├── tools
│   │   ├── bandit.go
│   │   ├── checkov.go
│   │   ├── eslint.go
│   │   ├── jscpd.go
│   │   ├── osvscanner.go
│   │   ├── pydocstyle.go
│   │   ├── radon.go
│   │   ├── runpylint.go
│   │   ├── semgrep.go
│   │   ├── trivy.go
│   │   ├── trufflehog.go
│   │   └── vulture.go
│   ├── worker
│   │   └── worker.go
├── pkg
│   ├── common
│   │   ├── constants.go
│   │   └── cwd.json
│   └── utils.go
├── postman
│   └── armur-codescan.postman_collection.json
├── rule_config
│   ├── eslint
│   │   ├── eslint.config.js
│   │   ├── eslint_deadcode.config.js
│   │   ├── eslint_jsdoc.config.js
│   │   ├── eslint_security.config.js
│   │   └── mappings.json
└── shared_tmp

```

**Here's a breakdown of the main directories and their purposes:**

1. api: This directory includes routes and handler functions of the backend.
2. redis: contains redis logic for setting and retrieving data from the redis store
3. tasks: responsible for managing and implementing tasks, ensuring scan tasks are processed
4. tools: tools directory contains all implementation for scanning tools for each programming language
5. worker: using async, this worker is reponsible for making sure tasks are processed
6. cmd: main app server start logic
7. docs & postman: contains documentation and detailed usage of routes
8. shared_tmp: this folder will be mounted in the container and you can place local repositories for scanning 