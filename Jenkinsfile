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
                        echo 'Развернуть постгрес в докере №1'
                    }
                }
                stage('go-get') {
                    steps {
                        sh 'go get github.com/valyala/fasthttp'
                        sh 'go get github.com/buaazp/fasthttprouter'
                        sh 'go get github.com/jackc/pgx'
                        sh 'go get github.com/satori/go.uuid'
                        sh 'go get github.com/liderman/text-generator'
                        sh 'go get github.com/dgrijalva/jwt-go'
                        sh 'go get github.com/garyburd/redigo/redis'
                    }
                }
                stage('prepare config') {
                    steps {
                        sh 'cp ./application.cfg ./tests/unit/'
                        echo 'ls -la'
                        echo 'ls -la ./tests'
                        echo 'ls -la ./tests/unit'
                    }
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
}

