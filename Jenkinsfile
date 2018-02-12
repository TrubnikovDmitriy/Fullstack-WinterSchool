pipeline {
    agent any
    stages {
        stage('prepare environment') {
            parallel {
                stage('clear image-1') {
                    steps {
                        sh 'docker rm -f my_postgres_1 || true'
                    }
                }
                stage('clear image-2') {
                    steps {
                        sh 'docker rm -f my_postgres_2 || true'
                    }
                }
                stage('whoami') {
                    steps {
                        sh 'whoami'
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


        stage('deploy postgres in docker') {
            parallel {
                stage('shard-1') {
                    steps {
                        sh 'docker run -d --name my_postgres_1 -e POSTGRESQL_USER=jenkins -e POSTGRESQL_PASSWORD=pass -e POSTGRESQL_DATABASE=db_test -p 5432:5432 centos/postgresql-96-centos7'
                    }
                }
                stage('shard-2') {
                    steps {
                        sh 'docker run -d --name my_postgres_2 -e POSTGRESQL_USER=jenkins -e POSTGRESQL_PASSWORD=pass -e POSTGRESQL_DATABASE=db_test -p 5433:5432 centos/postgresql-96-centos7'
                    }
                }
                stage('application config') {
                    steps {
                        sh 'cp ./application.cfg tests/unit/application.cfg'
                    }
                }
            }
        }

        stage('create schema') {
            parallel {
                stage('schema-1') {
                    steps {
                        sh 'psql -h localhost -p 5433 -d db_test < ./migrations/V1__init.sql'
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

                stage('auth_test') {
                    steps {
                        sh 'go test ./tests/unit/match_test.go -v'
                    }
                }

                stage('game_test') {
                    steps {
                        sh 'go test ./tests/unit/game_test.go -v'
                    }
                }

                stage('MatchTreeForm_test') {
                    steps {
                        sh 'go test ./tests/unit/MatchTreeForm_test.go -v'
                    }
                }

                stage('person_test') {
                    steps {
                        sh 'go test ./tests/unit/person_test.go -v'
                    }
                }

                stage('tourney_test') {
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

