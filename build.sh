#!	/usr/bin/ksh

#	set environment
export  GOROOT=${GOROOT:=/usr/local/go}
export	PATH=${PATH}:$( go env GOPATH )/bin

export  PROJECT_NAME='HTTP-server'
export  PROJECT_PATH=$( pwd )
export	PROJECT_DEPENDENCIES='
	'
export	PROJECT_SOURCE='app.go main.go'
export	PROJECT_PACKAGES='httpServer'
export	PROJECT_TARGET='bin/httpServer'
#export  GOPATH="${PROJECT_PATH}"

#	if necessary, set default target
export	TARGET_LIST="$*"

if [ -z "${TARGET_LIST}" ]
then
	export	TARGET_LIST='compile'
fi

if [ "${TARGET_LIST}" == 'all' ]
then
	export	TARGET_LIST='dependencies compile test package verify'
fi

#	help message
if [ "${TARGET_LIST}" == 'help' -o "${TARGET_LIST}" == '--help' ]
then
	cat <<HELP_MESSAGE
$( basename ${0} ): specify one or more targets:
	clean: remove all required packages and the target file
	dependencies: install all dependencies
	compile: compile and link source code
	test: execute unit tests
	package: package the target in a container
	verify: execute integration tests
	run: run the target
HELP_MESSAGE

	exit 0
fi

#	perform every specified target
for TARGET in ${TARGET_LIST}
do
	case ${TARGET} in 

	#	remove all required packages and the target file
		clean )
		echo '+++ removing required packages and the target file'

		rm -rf pkg/*
		rm -rf src/go.mod src/go.sum

		for PACKAGE in ${PROJECT_PACKAGES}
		do
			rm -rf src/${PACKAGE}/go.mod src/${PACKAGE}/go.sum
		done

		rm -f ${PROJECT_TARGET}
		;;

	#	install all dependencies
		dependencies )
		echo '+++ downloading and installing all dependencies'

		#	main packge dependencies
		cd src
		go mod init aldebap/${PROJECT_NAME}

		for MODULE in ${PROJECT_DEPENDENCIES}
		do
			go get -u ${MODULE}
		done

		for PACKAGE in ${PROJECT_PACKAGES}
		do
			echo "require ${PROJECT_NAME}/${PACKAGE} v0.0.0-unpublished" >> go.mod
			echo "replace ${PROJECT_NAME}/${PACKAGE} v0.0.0-unpublished => ../src/${PACKAGE}" >> go.mod
			#echo go mod edit -replace ${PROJECT_NAME}/${PACKAGE}=../src/${PACKAGE}
		done

		cd ..

		#	create project packages module files
		for PACKAGE in ${PROJECT_PACKAGES}
		do
			cd src/${PACKAGE}
			go mod init ${PROJECT_NAME}/${PACKAGE}
			cd ../..
		done
		;;

	#	compile and link source code
		compile )
		echo '+++ compiling and linking source code'

		rm -f ${PROJECT_TARGET}

		if [ ! -d $( dirname ${PROJECT_TARGET} ) ]
		then
			mkdir $( dirname ${PROJECT_TARGET} )
		fi

		cd src
		go build -o ../${PROJECT_TARGET} ${PROJECT_SOURCE}
		cd ..
		;;

	#	execute unit tests
		test )
		echo '+++ running unit tests'

		export	TEST_RESULT=0

		cd src
		go test -v
		if [ $? -ne 0 ]
		then
			export	TEST_RESULT=-1
		fi
		cd ..

		for PACKAGE in ${PROJECT_PACKAGES}
		do
			cd src/${PACKAGE}
			go test -v
			if [ $? -ne 0 ]
			then
				export	TEST_RESULT=-1
			fi
			cd ../..
		done

		if [ ${TEST_RESULT} -ne 0 ]
		then
			exit ${TEST_RESULT}
		fi
		;;

	#	package the target in a Docker image
		package )
		echo '+++ package the target in a Docker image'
		#docker build --tag shopping-cart .
		;;

	#	execute integration tests
		verify )
		echo '+++ execute integration tests'

		${PROJECT_TARGET} &
		PID=$!
		#docker-compose up -d
		newman run 'test/Error handling scenarios.postman_collection.json' --environment 'test/Localhost.postman_environment.json'
		newman run 'test/Getting resources scenarios.postman_collection.json' --environment 'test/Localhost.postman_environment.json'
		#docker-compose stop
		kill -9 ${PID}
		;;

	#	run the target
		run )
		echo '+++ run the target'

		#docker-compose up
		${PROJECT_TARGET}
		;;

	#	unknown target
		* )
		echo "[error] unknown target ${TARGET}"
		;;
	esac
done
