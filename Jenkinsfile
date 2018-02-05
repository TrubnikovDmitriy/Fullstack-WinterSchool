pipeline {
    agent any
    stages {
        stage('begin') {
            parallel {
                stage('one-1') {
                    steps {
                        echo 'echo1'
                        echo 'echo2'
                        echo 'echo3'
                    }
                }    
                stage('one-2') {
                    steps {
                        sh 'go test ./tests/unit/match_test.go -v'
                        echo 'success'
                    }
                }
            }
        }
        
        stage('testing') {
            parallel {
                
                stage('match_test2') {
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
                
                
                stage('touney_test') {
                    steps {
                        sh 'go test ./tests/unit/tournament_test.go -v'
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
