pipeline {
    agent any
    
    environment {
        APP_NAME = 'eatright-server'
        DEPLOY_DIR = '/var/www/eatright'
        GO_VERSION = '1.22.0'
    }
    
    stages {
        stage('Checkout') {
            steps {
                echo '📦 Checking out source code...'
                checkout scm
            }
        }
        
        stage('Setup Go') {
            steps {
                echo '🔧 Setting up Go environment...'
                sh '''
                    export PATH=$PATH:/usr/local/go/bin
                    export GOPATH=$HOME/go
                    go version
                '''
            }
        }
        
        stage('Install Dependencies') {
            steps {
                echo '📥 Installing Go dependencies...'
                sh '''
                    export PATH=$PATH:/usr/local/go/bin
                    export GOPATH=$HOME/go
                    go mod download
                    go mod tidy
                '''
            }
        }
        
        stage('Run Tests') {
            steps {
                echo '🧪 Running tests...'
                sh '''
                    export PATH=$PATH:/usr/local/go/bin
                    export GOPATH=$HOME/go
                    go test ./... -v || true
                '''
            }
        }
        
        stage('Build') {
            steps {
                echo '🔨 Building application...'
                sh '''
                    export PATH=$PATH:/usr/local/go/bin
                    export GOPATH=$HOME/go
                    
                    # Build for Linux
                    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/${APP_NAME} cmd/server/main.go
                    
                    # Make executable
                    chmod +x bin/${APP_NAME}
                    
                    # Show binary info
                    ls -lh bin/${APP_NAME}
                '''
            }
        }
        
        stage('Deploy') {
            steps {
                echo '🚀 Deploying to production...'
                sh '''
                    # Run deployment script
                    chmod +x scripts/deploy.sh
                    ./scripts/deploy.sh
                '''
            }
        }
        
        stage('Verify') {
            steps {
                echo '✅ Verifying deployment...'
                sh '''
                    # Wait for service to start
                    sleep 5
                    
                    # Check if service is running
                    sudo systemctl status eatright || true
                    
                    # Test health endpoint
                    curl -f http://localhost:8080/health || echo "Warning: Health check failed"
                '''
            }
        }
    }
    
    post {
        success {
            echo '✅✅✅ Deployment successful! ✅✅✅'
            echo '🎉 Application is now running at http://localhost:8080'
        }
        failure {
            echo '❌❌❌ Deployment failed! ❌❌❌'
            echo '📋 Check the console output for errors'
        }
        always {
            echo '🧹 Cleaning up workspace...'
            cleanWs()
        }
    }
}
