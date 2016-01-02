def setConsulVariables server, configfile
  puts "Setting Consul key values for server: #{server}"


  yaml = YAML.load_file(configfile)
  config = []
  process_yaml("", yaml).each {|k| config << k}

  config.each do |c|
    if c[:v].is_a?(Array)
      value = c[:v].to_json
    else
      value = c[:v]
    end
    puts `curl -X PUT -d '#{value}' http://#{server}:9500/v1/kv#{c[:k]}`
  end
end

def process_yaml root, hash
  keys = []
  return [] unless hash

  hash.each do |key, value|
    if value.is_a?(Hash)
      process_yaml(root + "/" + key, value).each do |k|
        keys << k
      end
    else
      keys << {:k => root + "/" + key, :v => value}
    end
  end
  return keys
end
