namespace :app do
  desc "updates build images for swagger and golang will overwrite existing images"
  task :update_images do
  	pull_image 'golang:latest'
  	#pull_image 'sandcastle/swagger-codegen-docker:latest'
  end

  desc "pull images for golang from Docker registry if not already downloaded"
  task :fetch_images do
  	pull_image 'golang' unless find_image 'golang:latest'
  	#pull_image 'sandcastle/swagger-codegen-docker' unless find_image 'sandcastle/swagger-codegen-docker:latest'
  end

  desc "run unit tests"
  task :test => [:fetch_images] do
  	p "Test Application"
  	container = get_container GO_CONTAINER_ARGS

  	begin
  		# Get go packages
  		ret = container.exec(['go','get','-t','-v','./...']) { |stream, chunk| puts "#{stream}: #{chunk}" }
  		raise Exception, 'Error running command' unless ret[2] == 0

  		# Test application
  		ret = container.exec(['go','test','./...']) { |stream, chunk| puts "#{stream}: #{chunk}" }
  		raise Exception, 'Error running command' unless ret[2] == 0
  	ensure
  		container.delete(:force => true)
  	end
  end

  desc "build and test application"
  task :build => [:fetch_images, :test] do
  	p "Build for Linux"
  	container = get_container GO_CONTAINER_ARGS

  	begin
  		# Build go server
  		ret = container.exec(['go','build','-a','-installsuffix','cgo','-ldflags','\'-s\'','-o','server']) { |stream, chunk| puts "#{stream}: #{chunk}" }
  		raise Exception, 'Error running command' unless ret[2] == 0
  	ensure
  		container.delete(:force => true)
  	end
  end

  desc "build Docker image for application"
  task :build_server => [:build] do
  	p "Building Docker Image:- #{DOCKER_IMAGE_NAME}"

  	FileUtils.cp "#{GOPATH}/src/#{GONAMESPACE}/#{DOCKER_IMAGE_NAME}/server", "./dockerfile/#{DOCKER_IMAGE_NAME}/server"
  	Dir.mkdir "./dockerfile/#{DOCKER_IMAGE_NAME}/swagger_spec/" unless Dir.exist? "./dockerfile/#{DOCKER_IMAGE_NAME}/swagger_spec/"
  	FileUtils.cp "./swagger_spec/swagger.yml", "./dockerfile/#{DOCKER_IMAGE_NAME}/swagger_spec/swagger.yml"

  	Docker.options = {:read_timeout => 6200}
  	image = Docker::Image.build_from_dir "./dockerfile/#{DOCKER_IMAGE_NAME}", {:t => DOCKER_IMAGE_NAME}
  end

  desc "run application with Docker Compose"
  task :run do
  	begin
      puts `docker-compose -f ./dockercompose/#{DOCKER_IMAGE_NAME}/docker-compose.yml up -d`
      sleep 2
  		setConsulVariables get_docker_ip_address, 9500

      sh "docker-compose -f ./dockercompose/#{DOCKER_IMAGE_NAME}/docker-compose.yml logs"
  	rescue SystemExit, Interrupt
  		sh "docker-compose -f ./dockercompose/#{DOCKER_IMAGE_NAME}/docker-compose.yml stop"
  		# remove stopped containers
  		sh "echo y | docker-compose -f ./dockercompose/#{DOCKER_IMAGE_NAME}/docker-compose.yml rm"
  	end
  end

  desc "build and run application with Docker Compose"
  task :build_and_run => [:build_server, :run]

  desc "run end to end Cucumber tests"
  task :e2e do
  	feature = ARGV.last
  	if feature != "app:e2e"
  		feature = "--tags #{feature}"
  	else
  		feature = ""
  	end

  	host = get_docker_ip_address

  	puts "Running Tests for #{host}"

  	ENV['WEB_SERVER_URI'] = "http://#{host}:8001"
  	ENV['MONGO_URI'] = "#{host}:27017"
  	ENV['EMAIL_SERVER_URI'] = "http://#{host}:1080"

    puts "Running Tests"
  	begin
  	  puts `docker-compose -f ./dockercompose/#{DOCKER_IMAGE_NAME}/docker-compose.yml up -d`
      sleep 2
  		setConsulVariables host, 9500
  		self.wait_until_server_running ENV['WEB_SERVER_URI']

  		p 'Running Tests'
  		puts `cucumber --color --strict -f pretty #{feature}`
  	ensure
      p 'Stopping Application'
  		# remove stop running application
  		puts `docker-compose -f ./dockercompose/#{DOCKER_IMAGE_NAME}/docker-compose.yml stop`
  		# remove stopped containers
  		puts `echo y | docker-compose -f ./dockercompose/#{DOCKER_IMAGE_NAME}/docker-compose.yml rm`
      abort "Cucumber steps failed" unless $FAILED == 0
  	end
  end

  desc "push built image to Docker registry"
  task :push do
  	p "Push image to registry"

  	image =  find_image "#{DOCKER_IMAGE_NAME}:latest"
  	image.tag('repo' => "#{DOCKER_NAMESPACE}#{DOCKER_IMAGE_NAME}", 'force' => true) unless image.info["RepoTags"].include? "#{DOCKER_NAMESPACE}#{DOCKER_IMAGE_NAME}:latest"

  	sh "docker login -u #{REGISTRY_USER} -p #{REGISTRY_PASS} -e #{REGISTRY_EMAIL} #{REGISTRY_URL}"
  	sh "docker push #{DOCKER_NAMESPACE}#{DOCKER_IMAGE_NAME}:latest"
  end
end
