pipeline {
    agent any
    stages {
        stage('prepare environment') {
            parallel {
                stage('shard-1') {
                    steps {
                        echo 'docker run --rm --name my_postgres_1 -e POSTGRESQL_USER=jenkins -e POSTGRESQL_PASSWORD=pass -e POSTGRESQL_DATABASE=db_test -p 5432:5432 centos/postgresql-96-centos7'
                    }
                }
                stage('shard-2') {
                    steps {
                        echo 'docker run --rm --name my_postgres_2 -e POSTGRESQL_USER=jenkins -e POSTGRESQL_PASSWORD=pass -e POSTGRESQL_DATABASE=db_test -p 5433:5432 centos/postgresql-96-centos7'
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

