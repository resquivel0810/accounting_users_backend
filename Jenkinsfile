pipeline {
    agent any

    environment {
        GITHUB_TOKEN = credentials('github-token-id')  // GitHub token for API access
        GITHUB_REPO = 'accounting_users_be'           // GitHub repository name
        GITHUB_USER = 'github-username'               // GitHub username
        DOCKERFILE_NAME = 'Dockerfile'                // Default Dockerfile name
    }

    stages {
        stage('Set Environment') {
            steps {
                script {
                    // Determine the Dockerfile based on the Git branch
                    if (env.GIT_BRANCH == 'origin/develop') {
                        env.DOCKERFILE_NAME = 'Dockerfile'
                        env.agentLabel = 'dev'
                        sh "cp './.env.dev' './.env'"  // Copy development environment variables
                    } else if (env.GIT_BRANCH == 'origin/main') {
                        env.DOCKERFILE_NAME = 'Dockerfile.prod'
                        env.agentLabel = 'prod'
                        sh "cp './.env.prod' './.env'"  // Copy production environment variables
                    }
                }
                dir("${workspace}") {
                    stash name: 'modifiedenv', includes: '.env'  // Stash the .env file for later use
                }
            }
        }

        stage('Build') {
            steps {
                script {
                    // Notify GitHub of a pending build status
                    sh "curl -H 'Authorization: token ${GITHUB_TOKEN}' " +
                       "-d '{\"state\": \"pending\", \"description\": \"Build in progress\", " +
                       "\"context\": \"continuous-integration/jenkins\"}' " +
                       "https://api.github.com/repos/${GITHUB_USER}/${GITHUB_REPO}/statuses/${env.GIT_COMMIT}"
                }
                dir("${workspace}") {
                    unstash 'modifiedenv'  // Retrieve the stashed .env file
                }
                script {
                    try {
                        // Execute build script with appropriate Dockerfile
                        sh "./build.sh ${agentLabel} ${DOCKERFILE_NAME}"
                        // Notify GitHub of a successful build status
                        sh "curl -H 'Authorization: token ${GITHUB_TOKEN}' " +
                           "-d '{\"state\": \"success\", \"description\": \"Build successful\", " +
                           "\"context\": \"continuous-integration/jenkins\"}' " +
                           "https://api.github.com/repos/${GITHUB_USER}/${GITHUB_REPO}/statuses/${env.GIT_COMMIT}"
                    } catch (Exception e) {
                        // Notify GitHub of a failed build status
                        sh "curl -H 'Authorization: token ${GITHUB_TOKEN}' " +
                           "-d '{\"state\": \"failure\", \"description\": \"Build failed\", " +
                           "\"context\": \"continuous-integration/jenkins\"}' " +
                           "https://api.github.com/repos/${GITHUB_USER}/${GITHUB_REPO}/statuses/${env.GIT_COMMIT}"
                        currentBuild.result = 'FAILURE'  // Mark build as a failure
                    }
                }
            }
        }
    }
}
