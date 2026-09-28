pipeline {
  agent any
  options {
    timestamps()
    disableConcurrentBuilds()
  }
  parameters {
    booleanParam(
      name: 'DEPLOY_VPS',
      defaultValue: false,
      description: 'Deploy the validated commit to the VPS. Keep disabled for CI-only runs.'
    )
    string(
      name: 'VPS_HOST',
      defaultValue: '167.86.79.16',
      description: 'VPS address. SSH host keys must already be trusted by the Jenkins agent.'
    )
    string(
      name: 'VPS_USER',
      defaultValue: 'yacine',
      description: 'VPS SSH user. Prefer a restricted deployment user when available.'
    )
    string(
      name: 'VPS_APP_DIR',
      defaultValue: '/opt/aisha',
      description: 'Directory containing the deployed AISHA source and persistent Compose project.'
    )
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
    stage('Deploy VPS') {
      when {
        expression { params.DEPLOY_VPS }
      }
      steps {
        script {
          if (!(params.VPS_HOST ==~ /^[A-Za-z0-9.-]+$/)) {
            error('VPS_HOST contains unsupported characters')
          }
          if (!(params.VPS_USER ==~ /^[A-Za-z_][A-Za-z0-9_-]*$/)) {
            error('VPS_USER contains unsupported characters')
          }
          if (!(params.VPS_APP_DIR ==~ /^\/[A-Za-z0-9._\/-]+$/)) {
            error('VPS_APP_DIR must be an absolute safe path')
          }
        }
        withCredentials([
          file(credentialsId: 'aisha-vps-production-env', variable: 'PRODUCTION_ENV_FILE')
        ]) {
          sshagent(credentials: ['aisha-vps-ssh']) {
            sh '''#!/usr/bin/env bash
              set -Eeuo pipefail

              : "${GIT_COMMIT:?Jenkins must provide GIT_COMMIT}"
              remote="${VPS_USER}@${VPS_HOST}"
              remote_env="${VPS_APP_DIR}/.env.jenkins.${BUILD_NUMBER}"
              ssh_opts=(-o BatchMode=yes -o StrictHostKeyChecking=yes)

              ssh "${ssh_opts[@]}" "$remote" "mkdir -p '$VPS_APP_DIR'"
              rsync -az \
                --exclude='.git' \
                --exclude='.env' \
                --exclude='backups' \
                --exclude='apps/web/node_modules' \
                --exclude='apps/web/.next' \
                --exclude='apps/api/bin' \
                ./ "$remote:$VPS_APP_DIR/"

              scp "${ssh_opts[@]}" "$PRODUCTION_ENV_FILE" "$remote:$remote_env"
              ssh "${ssh_opts[@]}" "$remote" \
                "install -m 600 '$remote_env' '$VPS_APP_DIR/.env' && rm -f '$remote_env'"

              ssh "${ssh_opts[@]}" "$remote" \
                "cd '$VPS_APP_DIR' && chmod +x infrastructure/scripts/deploy-vps.sh && infrastructure/scripts/deploy-vps.sh '$VPS_APP_DIR' '$GIT_COMMIT'"
            '''
          }
        }
      }
    }
  }
}
