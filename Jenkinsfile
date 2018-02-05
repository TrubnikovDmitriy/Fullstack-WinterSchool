pipeline {
    agent any
    stages {
        stage('Example') {
            steps {
                bash 'go run main'
            }
        }
    }
    post { 
        always { 
            echo 'I will always say Hello again!'
        }
    }
}
