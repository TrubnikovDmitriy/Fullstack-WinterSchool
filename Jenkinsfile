pipeline {
    agent any
    stages {
        parallel stage('game_test') {
            steps {
                sh 'go test ./tests/unit/game_test.go -v'
            }
        }
        
        parallel stage('match_test') {
            steps {
                sh 'go test ./tests/unit/match_test.go -v'
            }
        }
        parallel stage('person_test') {
            steps {
                sh 'go test ./tests/unit/person_test.go -v'
            }
        }
        parallel stage('player_test') {
            steps {
                sh 'go test ./tests/unit/player_test.go -v'
            }
        }
        parallel stage('team_test') {
            steps {
                sh 'go test ./tests/unit/team_test.go -v'
            }
        }
    }
    post { 
        always { 
            echo 'I will always say Hello again!'
        }
    }
}
