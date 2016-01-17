

def self.wait_until_server_running server, count
  begin
    response = RestClient.send("get", "#{server}/v1/health")
  rescue

  end

  if response == nil || !response.code.to_i == 200
    puts "Waiting for server to start"
    sleep 1
    if count < 20
      self.wait_until_server_running server, count + 1
    else
      raise 'Server failed to start'
    end
  end
end
