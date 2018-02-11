pipeline {
    agent any
    stages {
        stage('prepare environment') {
            parallel {
                stage('shard-1') {
                    steps {
                        echo 'Развернуть постгрес в докере №1'
                    }
                }
                stage('shard-2') {
                    steps {
                        echo 'Развернуть постгрес в докере №2'
                    }
                }
            }
        }

        stage('prepare config') {
            stage('copy') {
                steps {
                    sh 'cp ./application.cfg' './tests/unit/'
                }
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
            echo 'CI completed!'
        }
    }
