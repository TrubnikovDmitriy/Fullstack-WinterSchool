pipeline {
    agent any
    stages {
        stage('Example') {
            steps {
                sh 'go test ./tests/unit/ -v'
            }
        }
    }
    post { 
        always { 
            echo 'I will always say Hello again!'
        }
    }
}
