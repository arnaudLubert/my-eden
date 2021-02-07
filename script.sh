#!/bin/bash
# Editor: arnaud.lubert@epitech.eu

USAGE="\r\n📋 Usage:\t./script.sh <command> [arguments]\r\n\r\nThe commands are:\r\n\tinstall\t\tinstall project dependecies\r\n\tbuild\t\tbuild services\r\n\tstart\t\tlaunch services\r\n\tstop\t\tstop services\r\n\tlog\t\tdisplay container logs\r\n\trestart\t\trestart golang container (when mysql has too much leaks)\r\n\treboot\t\treboot the host at 4:00 am\r\n\r\nThe arguments are:\r\n\t-d\t\tdebug mode (to buid|start golang outside docker container)\r\n\t-go\t\tOnly affect golang service\r\n\t-sql\t\tOnly affect MySQL service\r\n\t-ftp\t\tOnly affect pure-ftpd service\r\n\t-ngx\t\tOnly affect NginxMonitor service\r\n\t-cad\t\tOnly affect Cadvisor service\r\n\t-exp\t\tOnly affect sysexporter service\r\n\t-lo\t\tlocal execution (http)\r\n"

if [ $# -eq 0 ]; then
    echo -e $USAGE
elif [ -r ./env.config ]; then
    if [ $1 == "install" -o $1 == "build" -o $1 == "start" -o $1 == "restart" -o $1 == "stop" -o $1 == "log" -o $1 == "reboot" ]; then
        source ./env.config

        export RELATIVE_PATH=$RELATIVE_PATH
        export URL_SCHEME=$URL_SCHEME
        export MYSQL_PORT=$MYSQL_PORT
        export LOCAL_IP=$LOCAL_IP

        EXEC_NAME=server
        PROJECT=./containers/golang

        for ac in $@
        do
            if [ $ac == "-d" ]; then
                debug=true
                export DEBUG=true
                export LOCAL_IP=$LOCAL_IP
                export SMTP_HOST=$SMTP_HOST
                export SMTP_ADDR=$SMTP_ADDR
                export MAIL_RECEIVER_1=$MAIL_RECEIVER_1
                export MAIL_RECEIVER_2=$MAIL_RECEIVER_2
                export MAIL_RECEIVER_COMM=$MAIL_RECEIVER_COMM
                export MAIL_SAV=$MAIL_SAV
                export MAIL_SAV_PASS=$MAIL_SAV_PASS
                export MAIL_RECEIVER_RC=$MAIL_RECEIVER_RC
                export MAIL_SENDER=$MAIL_SENDER
                export MAIL_PASS=$MAIL_PASS
                export GOOGLE_CLI_ID=$GOOGLE_CLI_ID
                export GOOGLE_SECRET=$GOOGLE_SECRET
                export MONDAY_API_KEY=$MONDAY_API_KEY
                export MONDAY_APP_KEY=$MONDAY_APP_KEY
                export DOMAIN_NAME=$DOMAIN_NAME
                export MYSQL_USER=$MYSQL_USER
                export MYSQL_PASSWORD=$MYSQL_PASSWORD
                export MYSQL_PORT=$MYSQL_PORT
                export MYSQL_DATABASE=$MYSQL_DATABASE
                export MAINTENER_EMAIL=$MAINTENER_EMAIL
                export SSL_CERT=$SSL_CERT
                export SSL_KEY=$SSL_KEY
            elif [ $ac == "-go" ]; then
                onlyGo=true
            elif [ $ac == "-sql" ]; then
                onlyMySQL=true
	        elif [ $ac == "-ftp" ]; then
                onlyFTP=true
            elif [ $ac == "-ngx" ]; then
                onlyNGX=true
            elif [ $ac == "-cad" ]; then
                onlyCAD=true
            elif [ $ac == "-exp" ]; then
                onlyEXP=true
            elif [ $ac == "-sqlexp" ]; then
                onlySQLexp=true
            elif [ $ac == "-lo" ]; then
                Local=true
                export URL_SCHEME="http://"
            elif [ $ac != $1 ]; then
                echo -e $USAGE
                exit 1
            fi

            if [ $Local ]; then
                DOCKER_CMD="docker-compose -f compose-base.yaml -f compose-local.yaml -p $DOCKER_NAMESPACE"
            else
                DOCKER_CMD="docker-compose -f compose-base.yaml -p $DOCKER_NAMESPACE"
            fi
        done

        if [ $1 == "install" ]; then
            #    git clone ssh://git@gitlab.e-frogg.com:65252/signatix/signatix-front.git # https://gitlab.e-frogg.com:65253/signatix/signatix-front.git

            if [ $debug ]; then
                echo "try apt install golang (1.13)"
                apt install golang
                go version
                echo "golang version must be >= 1.13 (enter to continue)"
                read
            fi
            echo 'Reboot? (y/n)' && read x && [[ "$x" == "y" ]] && /sbin/reboot;
            echo "Then use: systemctl enable signatix.service"

        elif [ $1 == "build" ]; then
            if [ $onlyMySQL ]; then
                echo -e "Building..."
                $DOCKER_CMD build mysql
            elif [ $onlyGo ]; then
                if [ $debug ]; then
                    cd $PROJECT; go build -o $EXEC_NAME ./src &
                    # loading animation
                    pid=$! ; i=0
                    while ps -a | awk '{print $1}' | grep -q "${pid}"
                    do
                        c=`expr ${i} % 14`
                        case ${c} in
                            0) echo -e "👩‍🔬 Building...🚧👷.............🧱🏗️\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\c" ;;
                            1) echo -e "👩‍🔬 Building...🚧.👷............🧱🏗️\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\c" ;;
                            2) echo -e "👩‍🔬 Building...🚧..👷...........🧱🏗️\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\c" ;;
                            3) echo -e "👩‍🔬 Building...🚧...👷..........🧱🏗️\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\c" ;;
                            4) echo -e "👩‍🔬 Building...🚧....👷.........🧱🏗️\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\c" ;;
                            5) echo -e "👩‍🔬 Building...🚧.....👷........🧱🏗️\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\c" ;;
                            6) echo -e "👩‍🔬 Building...🚧......👷.......🧱🏗️\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\c" ;;
                            7) echo -e "👩‍🔬 Building...🚧.......👷......🧱🏗️\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\c" ;;
                            8) echo -e "👩‍🔬 Building...🚧........👷.....🧱🏗️\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\c" ;;
                            9) echo -e "👩‍🔬 Building...🚧.........👷....🧱🏗️\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\c" ;;
                            10) echo -e "👩‍🔬 Building...🚧..........👷...🧱🏗️\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\c" ;;
                            11) echo -e "👩‍🔬 Building...🚧...........👷..🧱🏗️\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\c" ;;
                            12) echo -e "👩‍🔬 Building...🚧............👷.🧱🏗️\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\c" ;;
                            13) echo -e "👩‍🔬 Building...🚧.............👷🧱🏗️\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\c" ;;
                        esac
                        i=`expr ${i} + 1`
                        sleep 0.3
                    done
                    echo -e "--------------------------------------------------------------------------------"
                    wait ${pid}
                    # loading animation

                    if [ $? == "0" ]; then
                        mv $EXEC_NAME ../..
                    fi
                else
                    echo -e "Building..."
                    $DOCKER_CMD build golang
                    if [ $? == "0" ]; then
                        docker image prune --filter label=build=golang-temp
                    fi
                fi
            elif [ $onlyFTP ]; then
                echo -e "Building..."
                $DOCKER_CMD build pureftpd
            elif [ $onlyNGX ]; then
                echo -e "Building..."
                $DOCKER_CMD build nginxmonitor
            elif [ $onlyCAD ]; then
                echo -e "Building..."
                $DOCKER_CMD build cadvisor
            elif [ $onlyEXP ]; then
                echo -e "Building..."
                $DOCKER_CMD build sysexporter
            else
                echo -e "Building..."
                $DOCKER_CMD build
                if [ $? == "0" ]; then
                    docker image prune --filter label=build=golang-temp
                fi
            fi
            echo -e "Done"

        elif [ $1 == "start" ]; then
            if [ $onlyMySQL ]; then
                if [ $debug ]; then
                    $DOCKER_CMD up mysql
                else
                    $DOCKER_CMD up -d mysql
                fi
            elif [ $onlyGo ]; then
                if [ $debug ]; then
                    cd front; ../server
                else
                    $DOCKER_CMD up -d golang
                fi
            elif [ $onlyFTP ]; then
                if [ $debug ]; then
                    $DOCKER_CMD up pureftpd
                else
                    $DOCKER_CMD up -d pureftpd
                fi
            elif [ $onlyNGX ]; then
                if [ $debug ]; then
                    $DOCKER_CMD up nginxmonitor
                else
                    $DOCKER_CMD up -d nginxmonitor
                fi
            elif [ $onlyCAD ]; then
                if [ $debug ]; then
                    $DOCKER_CMD up cadvisor
                else
                    $DOCKER_CMD up -d cadvisor
                fi
            elif [ $onlyEXP ]; then
                if [ $debug ]; then
                    $DOCKER_CMD up sysexporter
                else
                    $DOCKER_CMD up -d sysexporter
                fi
            elif [ $onlySQLexp ]; then
                if [ $debug ]; then
                    $DOCKER_CMD up mysqlexporter
                else
                    $DOCKER_CMD up -d mysqlexporter
                fi
            else
                if [ $debug ]; then
                    $DOCKER_CMD up mysql
                    cd front; ../server
                else
                    $DOCKER_CMD up -d
                fi
            fi

        elif [ $1 == "stop" ]; then
            if [ $onlyMySQL ]; then
                $DOCKER_CMD stop mysql
            elif [ $onlyGo ]; then
                if [ $debug ]; then
                    pkill -f "../server"
                else
                    $DOCKER_CMD stop golang
                fi
            elif [ $onlyFTP ]; then
                $DOCKER_CMD stop pureftpd
            elif [ $onlyNGX ]; then
                $DOCKER_CMD stop nginxmonitor
            elif [ $onlyCAD ]; then
                $DOCKER_CMD stop cadvisor
            elif [ $onlyEXP ]; then
                $DOCKER_CMD stop sysexporter
            elif [ $onlySQLexp ]; then
                $DOCKER_CMD stop mysqlexporter
            else
                if [ $debug ]; then
                    $DOCKER_CMD stop mysql
                    pkill -f "../server"
                else
                    $DOCKER_CMD down
                fi
            fi

        elif [ $1 == "log" ]; then
            if [ $onlyMySQL ]; then
                $DOCKER_CMD logs mysql
            elif [ $onlyGo ]; then
                if [ $debug ]; then
                    echo "❌ Cannot read output outside Docker"
                else
                    $DOCKER_CMD logs golang
                fi
            elif [ $onlyFTP ]; then
                $DOCKER_CMD logs pureftpd
            elif [ $onlyNGX ]; then
                $DOCKER_CMD logs nginxmonitor
            elif [ $onlyCAD ]; then
                $DOCKER_CMD logs cadvisor
            elif [ $onlyEXP ]; then
                $DOCKER_CMD logs sysexporter
            else
                if [ $debug ]; then
                    echo "❌ Cannot read output outside Docker"
                else
                    $DOCKER_CMD logs golang
                fi
            fi

        elif [ $1 == "restart" ]; then
            if [ $onlyMySQL ]; then
                $DOCKER_CMD stop golang; $DOCKER_CMD start golang
            fi

        elif [ $1 == "reboot" ]; then
            sudo shutdown -r 04:00 &
        fi

    else
        echo -e $USAGE
    fi
else
    echo -e "\r\n❌ env.config is missing (use env.config-sample)\r\n"
fi
