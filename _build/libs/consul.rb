def setConsulVariables host, port
  puts "Setting Consul key values for server: #{host}:#{port}"

  kvs = Consul::Client::KeyValue.new :api_host => host, :api_port => port, :logger => Logger.new("/dev/null")

  kvs.put('/api/eventsauce/stats_d_server_url','statsd:8125')

  kvs.put('/api/eventsauce/data_store/connection_string','mongodb://mongo/event-sauce')
  kvs.put('/api/eventsauce/data_store/database_name','event-sauce')

  kvs.put('/api/eventsauce/queue/connection_string','redis:6379')
  kvs.put('/api/eventsauce/queue/event_queue','event_queue')
  kvs.put('/api/eventsauce/queue/dead_letter_queue','dead_letter_queue')
end
