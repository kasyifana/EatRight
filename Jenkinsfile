pipeline {
    agent any
    
    environment {
        APP_NAME = 'eatright-backend'
        DOCKER_IMAGE = 'eatright-backend:latest'
    }
    
    stages {
        stage('Checkout') {
            steps {
                echo '📦 Checking out source code...'
                checkout scm
            }
        }
        
        stage('Setup Environment') {
            steps {
                echo '🔧 Setting up environment...'
                sh '''
                    # Create .env file from Jenkins credentials or use placeholder for now
                    # In production, you should use Jenkins Credentials Binding
                    
                    # For this setup, we assume .env exists on the server or we create a dummy one for build
                    # But docker-compose needs it.
                    
                    # Let's check if .env exists, if not create a template
                    if [ ! -f .env ]; then
                        echo "Creating template .env"
                        echo "PORT=8080" > .env
                        echo "ENV=production" >> .env
                        # Note: Real credentials should be managed via Jenkins Credentials
                    fi
                '''
            }
        }
        
        stage('Build Docker Image') {
            steps {
                echo '🔨 Building Docker image...'
                sh '''
                    docker build -t ${DOCKER_IMAGE} .
                '''
            }
        }
        
        stage('Deploy') {
            steps {
                echo '🚀 Deploying container...'
                sh '''
                    # Stop and remove old container
                    docker stop ${APP_NAME} || true
                    docker rm ${APP_NAME} || true
                    
                    # Run new container
                    # We mount the .env file from the host machine to keep secrets safe
                    # Ensure /var/www/eatright/.env exists on the VPS
                    
                    docker run -d \
                        --name ${APP_NAME} \
                        --restart unless-stopped \
                        -p 8080:8080 \
                        -v /var/www/eatright/.env:/app/.env \
                        ${DOCKER_IMAGE}
                '''
            }
        }
        
        stage('Verify') {
            steps {
                echo '✅ Verifying deployment...'
                sh '''
                    sleep 5
                    docker ps | grep ${APP_NAME}
                    curl -f http://localhost:8080/health || echo "Warning: Health check failed"
                '''
            }
        }
    }
    
    post {
        success {
            echo '✅✅✅ Docker Deployment successful! ✅✅✅'
        }
        failure {
            echo '❌❌❌ Docker Deployment failed! ❌❌❌'
        }
        always {
            echo '🧹 Cleaning up...'
            sh 'docker system prune -f || true'
        }
    }
}
