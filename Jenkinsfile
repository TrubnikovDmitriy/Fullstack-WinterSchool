pipeline {
  agent any
  stages {
    stage('Example') {
      parallel {
        stage('Example') {
          steps {
            echo 'Hello World'
          }
        }
        stage('') {
          steps {
            sh 'echo \'paral 2\''
          }
        }
      }
    }
    stage('paral 2-1') {
      parallel {
        stage('paral 2-1') {
          steps {
            sh 'echo \'paral 2-1\''
            sh 'echo \'paral 2-1 (step2)\''
          }
        }
        stage('paral2-2') {
          steps {
            sh 'echo \'paral 2-2\''
          }
        }
        stage('paral 2-3') {
          steps {
            sh 'echo \'paral 2-3\''
          }
        }
      }
    }
    stage('paral 3-1') {
      steps {
        sh 'echo \'paral 3-1\''
      }
    }
  }
  post {
    always {
      echo 'I will always say Hello again!'
      
    }
    
  }
}