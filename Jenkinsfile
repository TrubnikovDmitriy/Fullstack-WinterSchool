pipeline {
    agent any
    stages {
        stage('match_test') {
            steps {
                
            }
        }
        
        stage('testing') {
            parallel {
                
                stage('match_test') {
                    steps {
                        sh 'go test ./tests/unit/match_test.go -v'
                    }
                }
                
                stage('game_test') {
                    steps {
                        sh 'go test ./tests/unit/game_test.go -v'
                    }
                }

                stage('match_test') {
                    steps {
                        sh 'go test ./tests/unit/match_test.go -v'
                    }
                }
                
                stage('person_test') {
                    steps {
                        sh 'go test ./tests/unit/person_test.go -v'
                    }
                }
                
                stage('player_test') {
                    steps {
                        sh 'go test ./tests/unit/player_test.go -v'
                    }
                }
                
                stage('team_test') {
                    steps {
                        sh 'go test ./tests/unit/team_test.go -v'
                    }
                }
                
            }
        }
        
    }
    post { 
        always { 
            echo 'I will always say Hello again!'
        }
    }
}
