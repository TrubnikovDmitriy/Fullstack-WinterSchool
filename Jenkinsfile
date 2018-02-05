pipeline {
    agent any
    stages {
        stage('Example') {
            steps {
                sh 'go run main.go'
            }
        }
    }
    post { 
        always { 
            echo 'I will always say Hello again!'
        }
    }
}
