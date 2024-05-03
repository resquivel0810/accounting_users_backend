pipeline {
    agent any

    environment {
        GITHUB_TOKEN = credentials('github-token-id')
        GITHUB_REPO = 'accounting_users_be'
        GITHUB_USER = 'github-username'
    }

    stages {
        stage('Develop') {
            when {
                expression {env.GIT_BRANCH == 'origin/develop'}
            }
            steps {
                sh "cp './.env.dev' './.env'"
                dir("${workspace}") {
                    stash name: 'modifiedenv', includes: '.env'
                }
                script {
                    env.agentLabel='dev'
                }
            }
        }

        stage('Production') {
            when {
                expression {env.GIT_BRANCH == 'origin/main'}
            }
            steps {
                sh "cp './.env.prod' './.env'"
                dir("${workspace}") {
                    stash name: 'modifiedenv', includes: '.env'
                }
                script {
                    env.agentLabel='prod'
                }
            }
        }

        stage('Build') {
            steps {
                script {
                    sh "curl -H 'Authorization: token ${GITHUB_TOKEN}' " +
                       "-d '{\"state\": \"pending\", \"description\": \"Build in progress\", " +
                       "\"context\": \"continuous-integration/jenkins\"}' " +
                       "https://api.github.com/repos/${GITHUB_USER}/${GITHUB_REPO}/statuses/${env.GIT_COMMIT}"
                }
                dir("${workspace}") {
                    unstash 'modifiedenv'
                }
                script {
                    try {
                        sh './build.sh ${agentLabel}'
                        sh "curl -H 'Authorization: token ${GITHUB_TOKEN}' " +
                           "-d '{\"state\": \"success\", \"description\": \"Build successful\", " +
                           "\"context\": \"continuous-integration/jenkins\"}' " +
                           "https://api.github.com/repos/${GITHUB_USER}/${GITHUB_REPO}/statuses/${env.GIT_COMMIT}"
                    } catch (Exception e) {
                        sh "curl -H 'Authorization: token ${GITHUB_TOKEN}' " +
                           "-d '{\"state\": \"failure\", \"description\": \"Build failed\", " +
                           "\"context\": \"continuous-integration/jenkins\"}' " +
                           "https://api.github.com/repos/${GITHUB_USER}/${GITHUB_REPO}/statuses/${env.GIT_COMMIT}"
                        currentBuild.result = 'FAILURE'
                    }
                }
            }
        }
    }
}
