pipeline {
  agent any
  options {
    timestamps()
    disableConcurrentBuilds()
  }
  stages {
    stage('Repository validation') {
      steps {
        sh 'test -s .env.example'
        sh 'test ! -f .env'
        sh 'docker compose -f infrastructure/compose/docker-compose.yml config --quiet'
        sh 'docker compose -f infrastructure/compose/docker-compose.yml -f infrastructure/compose/docker-compose.prod.yml config --quiet'
      }
    }
    stage('Backend') {
      steps {
        dir('apps/api') {
          sh 'go mod download'
          sh 'test -z "$(gofmt -l .)"'
          sh 'go vet ./...'
          sh 'go test -race ./...'
          sh 'go build ./...'
        }
      }
    }
    stage('Database migrations') {
      steps {
        sh 'docker compose -f infrastructure/compose/docker-compose.yml up -d postgres'
        sh "docker compose -f infrastructure/compose/docker-compose.yml exec -T postgres sh -ec 'dropdb --if-exists -U aisha aisha_migration_test && createdb -U aisha aisha_migration_test'"
        sh 'make test-migrations'
        sh 'make migrate-up'
        sh 'make sqlboiler'
        sh 'git diff --exit-code -- apps/api/db/sqlboiler/models'
      }
    }
    stage('Frontend') {
      steps {
        dir('apps/web') {
          sh 'npm ci'
          sh 'npm run lint'
          sh 'npm run type-check'
          sh 'npm test'
          sh 'npm run build'
        }
      }
    }
    stage('Contracts and images') {
      steps {
        sh 'npx --yes @redocly/cli lint apps/api/openapi/openapi.yaml'
        sh 'docker compose -f infrastructure/compose/docker-compose.yml build'
      }
    }
  }
}
